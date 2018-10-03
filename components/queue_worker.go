package components

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// LoggerWithTask adds task information to a logger.
func LoggerWithTask(logger log.Logger, task *Task) log.Logger {
	return log.With(
		logger,
		"id", task.ID,
		"objectId", task.ObjectID,
		"retries", task.Retries,
	)
}

// NewQueueWorker returns a worker that synchronizes a SQL table formed as a list of pair documentID/timestamp with a channel.
func NewQueueWorker(l log.Logger, repo QueueRepository, table string, syncTick time.Duration, outCh chan<- Task, inCh <-chan Task) Worker {
	return NewWorker(
		l,
		&queueWorker{
			table:  table,
			repo:   repo,
			ticker: time.NewTicker(syncTick),
			outCh:  outCh,
			inCh:   inCh,
		},
	)
}

type queueWorker struct {
	table  string
	repo   QueueRepository
	ticker *time.Ticker
	lastID uint64
	outCh  chan<- Task
	inCh   <-chan Task
}

func (w *queueWorker) Name() string {
	return w.table
}

func (w *queueWorker) Start(ctx context.Context, l log.Logger) error { return nil }

func (w *queueWorker) Work(ctx context.Context, l log.Logger) {
	select {
	case <-w.ticker.C:
		if err := w.repo.RetryTasks(ctx, w.table, w.lastID); err != nil {
			l.Log("err", err)
			return
		}
		tasks, err := w.repo.GetQueue(ctx, w.table, w.lastID, cap(w.outCh)-len(w.outCh))
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
			if err := w.repo.UpdateTask(ctx, w.table, task.ID, task.Retries); err != nil {
				l.Log("err", err)
				return
			}
			return
		}
		if err := w.repo.DeleteFromQueue(ctx, w.table, task.ID); err != nil {
			l.Log("err", err)
			return
		}
	case <-ctx.Done():
		return
	}
}

func (w *queueWorker) Shutdown(ctx context.Context, l log.Logger) error { return nil }
