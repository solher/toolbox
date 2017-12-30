package components

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

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

// QueueRepository allows interaction with a Task queue.
type QueueRepository interface {
	GetQueue(ctx context.Context, table string, fromID uint64, limit int) (tasks []Task, err error)
	RetryTasks(ctx context.Context, table string, fromID uint64) error
	DeleteFromQueue(ctx context.Context, table string, id uint64) error
	DeleteFromQueueBulk(ctx context.Context, table string, ids []uint64) error
	UpdateTask(ctx context.Context, table string, id, retries uint64) error
	InsertTask(ctx context.Context, table string, task *Task) error
}

// NewPostgresQueueRepo returns a new Postgres backed QueueRepository.
func NewPostgresQueueRepo(db *sqlx.DB) QueueRepository {
	return &postgresQueueRepo{db: db}
}

type postgresQueueRepo struct {
	db *sqlx.DB
}

func (r *postgresQueueRepo) GetQueue(ctx context.Context, table string, fromID uint64, limit int) (tasks []Task, err error) {
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
		Table:   table,
	}
	query := bytes.NewBuffer(nil)
	if err := getPostgresQueueTmpl.Execute(query, arg); err != nil {
		return nil, errors.WithStack(err)
	}

	stmt, err := r.db.PrepareNamed(query.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer stmt.Close()

	if err := stmt.SelectContext(ctx, &tasks, arg); err != nil {
		return nil, errors.WithStack(err)
	}

	return tasks, nil
}

func (r *postgresQueueRepo) RetryTasks(ctx context.Context, table string, fromID uint64) error {
	_, err := r.db.ExecContext(
		ctx,
		fmt.Sprintf(`
		WITH t AS (
			SELECT id, object_id, retries
			FROM %[1]s
			WHERE id < $1
			AND created_at < NOW() - INTERVAL '10 minutes'
		), d AS (
      DELETE FROM %[1]s WHERE id IN (SELECT id FROM t)
		)
		INSERT INTO %[1]s (object_id, retries) SELECT t.object_id, t.retries FROM t ON CONFLICT (object_id) DO NOTHING 
		`, table),
		fromID,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *postgresQueueRepo) DeleteFromQueue(ctx context.Context, table string, id uint64) error {
	if _, err := r.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", table), id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *postgresQueueRepo) DeleteFromQueueBulk(ctx context.Context, table string, ids []uint64) error {
	if _, err := r.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ANY($1)", table), types.Uint64Array(ids)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *postgresQueueRepo) UpdateTask(ctx context.Context, table string, id, retries uint64) error {
	_, err := r.db.ExecContext(
		ctx,
		fmt.Sprintf("UPDATE %s SET retries = $1, created_at = CURRENT_TIMESTAMP WHERE id = $2", table),
		retries,
		id,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *postgresQueueRepo) InsertTask(ctx context.Context, table string, task *Task) error {
	_, err := r.db.ExecContext(
		ctx,
		fmt.Sprintf("INSERT INTO %s (object_id, retries) VALUES ($1, $2) ON CONFLICT (object_id) DO NOTHING", table),
		task.ObjectID,
		task.Retries,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
