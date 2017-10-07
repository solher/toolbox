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
		SELECT id, object_id, retries, created_at
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
	ID        uint64    `json:"id"        db:"id"`
	ObjectID  uint64    `json:"objectId"  db:"object_id"`
	Retries   uint64    `json:"retries"   db:"retries"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	Err       error
}

// LoggerWithTask adds task information to a logger.
func LoggerWithTask(logger log.Logger, task *Task) log.Logger {
	return log.With(
		logger,
		"id", task.ID,
		"objectId", task.ObjectID,
		"retries", task.Retries,
	)
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
		if err := w.RetryTasks(ctx, w.lastID); err != nil {
			l.Log("err", err)
			return
		}
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
			task.Retries++
			if err := w.UpdateTask(ctx, task.ID, task.Retries); err != nil {
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
		return
	}
}

func (w *postgresQueue) Shutdown(ctx context.Context, l log.Logger) error { return nil }

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

	stmt, err := w.db.PrepareNamed(query.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer stmt.Close()

	if err := stmt.SelectContext(ctx, &tasks, arg); err != nil {
		return nil, errors.WithStack(err)
	}

	return tasks, nil
}

func (w *postgresQueue) RetryTasks(ctx context.Context, fromID uint64) error {
	_, err := w.db.ExecContext(
		ctx,
		fmt.Sprintf(`
		WITH t AS (
			SELECT id, object_id, retries
			FROM %s
			WHERE id < $1
			AND created_at < NOW() - INTERVAL '10 minutes'
		), d AS (
      DELETE FROM %s WHERE id IN (SELECT id FROM t)
		)
		INSERT INTO %s (object_id, retries) SELECT t.object_id, t.retries FROM t ON CONFLICT (object_id) DO NOTHING 
		`, w.table, w.table, w.table),
		fromID,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (w *postgresQueue) DeleteFromQueue(ctx context.Context, id uint64) error {
	if _, err := w.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", w.table), id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (w *postgresQueue) UpdateTask(ctx context.Context, id, retries uint64) error {
	_, err := w.db.ExecContext(
		ctx,
		fmt.Sprintf("UPDATE %s SET retries = $1, created_at = CURRENT_TIMESTAMP WHERE id = $2", w.table),
		retries,
		id,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (w *postgresQueue) InsertTask(ctx context.Context, task *Task) error {
	_, err := w.db.ExecContext(
		ctx,
		fmt.Sprintf("INSERT INTO %s (object_id, retries) VALUES ($1, $2) ON CONFLICT (object_id) DO NOTHING", w.table),
		task.ObjectID,
		task.Retries,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
