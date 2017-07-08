package components

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/solher/toolbox/sql/types"
)

var getPostgresQueueTmpl = template.Must(template.New("get_postgresQueue").Parse(`
		WITH queue AS (
			SELECT id, created_at
			FROM {{.Table}}
			WHERE TRUE
			{{if not .FromTime.IsZero -}}
			AND created_at >= :from_time
			{{end -}}
			{{if .FromID -}}
			AND id > :from_id
			{{end -}}
			ORDER BY created_at, id ASC
			LIMIT :id_limit
		)
		SELECT ARRAY_AGG(id) AS ids, MAX(id) AS last_id, MAX(created_at) AS last_timestamp
		FROM queue
	`))

// NewPostgresQueue returns a worker that synchronizes a postgresQueue table formed as a list of pair documentID/timestamp with a channel.
func NewPostgresQueue(l log.Logger, db *sqlx.DB, table string, syncTick time.Duration, idCh chan<- uint64, doneCh <-chan uint64) Worker {
	return NewWorker(
		l,
		&postgresQueue{
			table:  table,
			db:     db,
			ticker: time.NewTicker(syncTick),
			idCh:   idCh,
			doneCh: doneCh,
		},
	)
}

type postgresQueue struct {
	table         string
	db            *sqlx.DB
	ticker        *time.Ticker
	lastID        uint64
	lastTimestamp time.Time
	idCh          chan<- uint64
	doneCh        <-chan uint64
}

func (w *postgresQueue) Name() string {
	return w.table
}

func (w *postgresQueue) Work(ctx context.Context, l log.Logger) {
	select {
	case <-w.ticker.C:
		ids, lastID, lastTimestamp, err := w.GetQueue(ctx, w.lastID, w.lastTimestamp, cap(w.idCh)-len(w.idCh))
		if err != nil {
			l.Log("err", err)
			return
		}
		for _, id := range ids {
			w.idCh <- id
		}
		if lastID != 0 && !lastTimestamp.IsZero() {
			w.lastID, w.lastTimestamp = lastID, lastTimestamp
		}
	case id := <-w.doneCh:
		if err := w.DeleteFromQueue(ctx, id); err != nil {
			l.Log("err", err)
			return
		}
	case <-ctx.Done():
		l.Log("err", ctx.Err())
		return
	}
}

func (w *postgresQueue) GetQueue(ctx context.Context, fromID uint64, fromTime time.Time, limit int) (ids []uint64, lastID uint64, lastTimestamp time.Time, err error) {
	if limit == 0 {
		return []uint64{}, 0, lastTimestamp, nil
	}

	arg := &struct {
		FromID   types.NullUint64 `db:"from_id"`
		FromTime types.NullTime   `db:"from_time"`
		IDLimit  int              `db:"id_limit"`
		Table    string
	}{
		FromID:   types.NullUint64(fromID),
		FromTime: types.NullTime(fromTime.UTC()),
		IDLimit:  limit,
		Table:    w.table,
	}
	query := bytes.NewBuffer(nil)
	if err := getPostgresQueueTmpl.Execute(query, arg); err != nil {
		return nil, 0, lastTimestamp, errors.WithStack(err)
	}

	stmt, err := w.db.PrepareNamedContext(ctx, query.String())
	if err != nil {
		return nil, 0, lastTimestamp, errors.WithStack(err)
	}

	dest := &struct {
		IDs           types.Uint64Array `db:"ids"`
		LastID        types.NullUint64  `db:"last_id"`
		LastTimestamp types.NullTime    `db:"last_timestamp"`
	}{}
	if err := stmt.GetContext(ctx, dest, arg); err != nil {
		return nil, 0, lastTimestamp, errors.WithStack(err)
	}

	return dest.IDs, uint64(dest.LastID), time.Time(dest.LastTimestamp), nil
}

func (w *postgresQueue) DeleteFromQueue(ctx context.Context, id uint64) error {
	if _, err := w.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", w.table), id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
