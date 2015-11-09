package core

import (
	"fmt"
	"time"
)

type Timer struct {
	done chan struct{}

	OnTimeOut func()
	OnStop    func()
}

func NewTimer() *Timer {
	return &Timer{
		done:      make(chan struct{}),
		OnTimeOut: func() {},
		OnStop:    func() {},
	}
}

func (this *Timer) Start(d time.Duration) *Timer {
	timer := &Timer{}

	go func() {
		t := time.NewTimer(d)
		for {
			select {
			case <-this.done:
				fmt.Println("Timer is stopped.")
				this.OnStop()
				return
			case <-t.C:
				fmt.Println("Timer is time out.")
				this.OnTimeOut()
				return
			}
		}
	}()

	fmt.Println("Timer is running.")
	return timer
}

func (this *Timer) Stop() {
	select {
	case <-this.done:
	default:
		close(this.done)
	}
}
