package components

import (
	"context"
	"errors"

	"sync"

	"github.com/go-kit/kit/log"
)

// Workable defines objects usable by the worker.
type Workable interface {
	Name() string
	Work(ctx context.Context, l log.Logger)
	Shutdown(ctx context.Context, l log.Logger) error
}

// Worker provides a synchronization api to a workable.
type Worker interface {
	Name() string
	Shutdown(ctx context.Context) error
	Start(ctx context.Context) error
}

type shutdownCallback struct {
	ctx      context.Context
	callback chan error
}

// worker allows to convert an endpoint to a worker listening on an event channel.
type worker struct {
	// We also use cancelWorkCtx as a semaphore to know when the worker is processing.
	cancelWorkCtx context.CancelFunc
	mutex         sync.Mutex

	name     string
	l        log.Logger
	workable Workable

	shutdownCh chan shutdownCallback
}

// NewWorker returns a new Worker.
func NewWorker(l log.Logger, workable Workable) Worker {
	return &worker{
		cancelWorkCtx: nil,
		name:          workable.Name(),
		l:             log.With(l, "component", workable.Name()),
		workable:      workable,
	}
}

// Shutdown shuts down synchronously and gracefully the worker processing.
func (w *worker) Shutdown(ctx context.Context) error {
	w.l.Log("msg", "shutting down")

	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.cancelWorkCtx == nil {
		return nil
	}

	w.cancelWorkCtx()
	w.cancelWorkCtx = nil

	callback := make(chan error)
	w.shutdownCh <- shutdownCallback{
		ctx:      ctx,
		callback: callback,
	}
	select {
	case err := <-callback:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Start launches the worker on the input channel.
func (w *worker) Start(ctx context.Context) error {
	w.mutex.Lock()
	if w.cancelWorkCtx != nil {
		w.mutex.Unlock()
		return errors.New("worker is already processing")
	}
	w.shutdownCh = make(chan shutdownCallback)
	workCtx, cancelWork := context.WithCancel(ctx)
	w.cancelWorkCtx = cancelWork
	w.mutex.Unlock()

	w.l.Log("msg", "successfully started")

	for {
		select {
		case <-ctx.Done():
			w.Shutdown(ctx)
		case shutdown := <-w.shutdownCh:
			close(w.shutdownCh)
			shutdown.callback <- w.workable.Shutdown(shutdown.ctx, w.l)
			return nil
		default:
			w.workable.Work(workCtx, w.l)
		}
	}
}

// Name returns the name of the worker.
func (w *worker) Name() string {
	return w.name
}
