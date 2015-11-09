package core

import (
	"sync"
	"time"
)

type Processer func(input interface{}) (output interface{}, err error)
type ErrorHandler func(input interface{}, err error) (err_output interface{})

type Worker struct {
	err_handler ErrorHandler

	done chan struct{}
	pool *HandlerPool
}

func NewWorker(h Processer, err_handler ErrorHandler, init, max int) *Worker {
	return &Worker{
		err_handler: err_handler,
		done:        make(chan struct{}),
		pool:        NewHandlerPool(h, init, max),
	}
}

func (this *Worker) Run(in_chan chan interface{}) (out_chan chan interface{}, err_chan chan interface{}) {
	out_chan = make(chan interface{}, 100)
	err_chan = make(chan interface{}, 100)

	go func() {
		defer func() {
			close(out_chan)
			close(err_chan)
		}()

		var wg sync.WaitGroup
		var v interface{}
		for {
			select {
			case v = <-in_chan:
				if v == nil {
					wg.Wait()
					return
				}
				handler := this.getHandler()
				if handler == nil {
					return
				}

				wg.Add(1)
				go func(h *Handler, v2 interface{}) {
					defer wg.Done()

					ret, err := h.Run(v2)
					this.pool.ReleaseHandler(h)

					if err != nil {
						err_chan <- this.err_handler(v2, err)
					} else {
						out_chan <- ret
					}
				}(handler, v)
			case <-this.done:
				wg.Wait()
				return
			}
		}
	}()

	return
}

func (this *Worker) Stop() {
	select {
	case <-this.done:
	default:
		close(this.done)
	}
}

func (this *Worker) getHandler() *Handler {
	handler, err := this.pool.Get()
	if err == nil {
		return handler
	}

	duration := time.Second * 30
	t := time.NewTimer(duration)
	for {
		select {
		case <-this.done:
			return nil
		case <-t.C:
			handler, err := this.pool.Get()
			if err == nil {
				return handler
			}
			t.Reset(duration)
		}
	}
}
