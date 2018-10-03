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
	Start(ctx context.Context, l log.Logger) error
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
	mutex   sync.Mutex
	running bool

	name     string
	l        log.Logger
	workable Workable

	shutdownCh, loopShutdownCh chan shutdownCallback
}

// NewWorker returns a new Worker.
func NewWorker(l log.Logger, workable Workable) Worker {
	return &worker{
		name:           workable.Name(),
		l:              log.With(l, "component", workable.Name()),
		workable:       workable,
		shutdownCh:     make(chan shutdownCallback, 1),
		loopShutdownCh: make(chan shutdownCallback, 1),
	}
}

// Shutdown shuts down synchronously and gracefully the worker processing.
func (w *worker) Shutdown(ctx context.Context) error {
	w.l.Log("msg", "shutting down")

	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.running = false

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
	if w.running {
		w.mutex.Unlock()
		return errors.New("worker is already processing")
	}
	w.running = true
	w.mutex.Unlock()

	if err := w.workable.Start(ctx, w.l); err != nil {
		return err
	}

	w.l.Log("msg", "successfully started")

	workCtx, cancelWork := context.WithCancel(ctx)
	defer cancelWork()

	go func() {
		for {
			select {
			case <-ctx.Done():
				w.Shutdown(ctx)
			case shutdown := <-w.loopShutdownCh:
				shutdown.callback <- w.workable.Shutdown(shutdown.ctx, w.l)
				return
			default:
				w.workable.Work(workCtx, w.l)
			}
		}
	}()

	select {
	case shutdown := <-w.shutdownCh:
		w.loopShutdownCh <- shutdown
		return nil
	}
}

// Name returns the name of the worker.
func (w *worker) Name() string {
	return w.name
}
