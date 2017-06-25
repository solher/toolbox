package components

import (
	"context"
	"errors"

	"sync"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Worker allows to convert an endpoint to a worker listening on an event channel.
type Worker struct {
	// We also use cancelProcessCtx as a semaphore to know when the worker is processing.
	cancelProcessCtx context.CancelFunc
	mutex            sync.Mutex

	logger    log.Logger
	endpoint  endpoint.Endpoint
	submitter func(ctx context.Context, response interface{}) error

	inputCh    <-chan interface{}
	shutdownCh chan chan error
}

// NewWorker returns a new instance of Worker.
func NewWorker(logger log.Logger, name string, inputCh <-chan interface{}, endpoint endpoint.Endpoint, submitter func(ctx context.Context, response interface{}) error) *Worker {
	return &Worker{
		logger:   log.With(logger, "component", "worker", "name", name),
		endpoint: endpoint,
		inputCh:  inputCh,
	}
}

// Shutdown shuts down synchronously and gracefully the worker processing.
func (w *Worker) Shutdown(ctx context.Context) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.cancelProcessCtx == nil {
		return nil
	}

	w.cancelProcessCtx()
	w.cancelProcessCtx = nil

	callback := make(chan error)
	w.shutdownCh <- callback
	select {
	case err := <-callback:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Process launches the worker processing on the input channel.
func (w *Worker) Process(ctx context.Context) error {
	w.mutex.Lock()
	if w.cancelProcessCtx != nil {
		w.mutex.Unlock()
		return errors.New("worker is already processing")
	}
	w.shutdownCh = make(chan chan error)
	ctx, w.cancelProcessCtx = context.WithCancel(ctx)
	w.mutex.Unlock()

	for {
		select {
		case <-ctx.Done():
			w.Shutdown(ctx)
		case callback := <-w.shutdownCh:
			close(w.shutdownCh)
			callback <- nil
			return nil
		default:
			select {
			case input := <-w.inputCh:
				output, err := w.endpoint(ctx, input)
				if err != nil {
					w.logger.Log("err", err)
					continue
				}
				if err := w.submitter(ctx, output); err != nil {
					w.logger.Log("err", err)
					continue
				}
			default:
				continue
			}
		}
	}
}
