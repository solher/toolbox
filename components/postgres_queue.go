package components

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/solher/toolbox/sql/types"
)

var getPostgresQueueTmpl = template.Must(template.New("get_postgres_queue").Parse(`
		SELECT id, object_id, retries, payload, created_at
		FROM {{.Table}}
		WHERE TRUE
		{{if .FromID -}}
		AND id > :from_id
		{{end -}}
		ORDER BY id ASC
		LIMIT :id_limit
	`))

// Task is a standardised queued task.
type Task struct {
	ID        uint64         `json:"id"        db:"id"`
	ObjectID  uint64         `json:"objectId"  db:"object_id"`
	Retries   uint64         `json:"retries"   db:"retries"`
	Payload   types.JSONText `json:"payload"   db:"payload"`
	CreatedAt time.Time      `json:"createdAt" db:"created_at"`
	Err       error
}

// NewPostgresQueue returns a worker that synchronizes a postgresQueue table formed as a list of pair documentID/timestamp with a channel.
func NewPostgresQueue(l log.Logger, db *sqlx.DB, table string, syncTick time.Duration, outCh chan<- Task, inCh <-chan Task) Worker {
	return NewWorker(
		l,
		&postgresQueue{
			table:  table,
			db:     db,
			ticker: time.NewTicker(syncTick),
			outCh:  outCh,
			inCh:   inCh,
		},
	)
}

type postgresQueue struct {
	table  string
	db     *sqlx.DB
	ticker *time.Ticker
	lastID uint64
	outCh  chan<- Task
	inCh   <-chan Task
}

func (w *postgresQueue) Name() string {
	return w.table
}

func (w *postgresQueue) Work(ctx context.Context, l log.Logger) {
	select {
	case <-w.ticker.C:
		tasks, err := w.GetQueue(ctx, w.lastID, cap(w.outCh)-len(w.outCh))
		if err != nil {
			l.Log("err", err)
			return
		}
		for _, task := range tasks {
			w.outCh <- task
		}
		if len(tasks) > 0 {
			w.lastID = tasks[len(tasks)-1].ID
		}
	case task := <-w.inCh:
		if task.Err != nil && task.Retries < 5 {
			if err := w.DeleteFromQueue(ctx, task.ID); err != nil {
				l.Log("err", err)
				return
			}
			task.Retries++
			if err := w.InsertTask(ctx, &task); err != nil {
				l.Log("err", err)
				return
			}
			return
		}
		if err := w.DeleteFromQueue(ctx, task.ID); err != nil {
			l.Log("err", err)
			return
		}
	case <-ctx.Done():
		l.Log("err", ctx.Err())
		return
	}
}

func (w *postgresQueue) GetQueue(ctx context.Context, fromID uint64, limit int) (tasks []Task, err error) {
	if limit == 0 {
		return tasks, nil
	}

	arg := &struct {
		FromID  types.NullUint64 `db:"from_id"`
		IDLimit int              `db:"id_limit"`
		Table   string
	}{
		FromID:  types.NullUint64(fromID),
		IDLimit: limit,
		Table:   w.table,
	}
	query := bytes.NewBuffer(nil)
	if err := getPostgresQueueTmpl.Execute(query, arg); err != nil {
		return nil, errors.WithStack(err)
	}

	stmt, err := w.db.PrepareNamedContext(ctx, query.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := stmt.SelectContext(ctx, &tasks, arg); err != nil {
		return nil, errors.WithStack(err)
	}

	return tasks, nil
}

func (w *postgresQueue) DeleteFromQueue(ctx context.Context, id uint64) error {
	if _, err := w.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", w.table), id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (w *postgresQueue) InsertTask(ctx context.Context, task *Task) error {
	_, err := w.db.ExecContext(
		ctx,
		fmt.Sprintf("INSERT INTO %s (object_id, retries, payload) VALUES ($1, $2, $3) ON CONFLICT (object_id) DO NOTHING", w.table),
		task.ObjectID,
		task.Retries,
		task.Payload,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
