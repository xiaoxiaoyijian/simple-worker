package core

import (
	"fmt"
	"testing"
	"time"
)

func Test_TimerTimeout(t *testing.T) {
	timer := NewTimer()
	timer.OnTimeOut = func() {
		fmt.Println("This is my timeout function.")
	}
	timer.OnStop = func() {
		fmt.Println("This is my stop function.")
	}

	timer.Start(time.Second * 3)
	time.Sleep(time.Second * 2)
	timer.Stop()

	time.Sleep(time.Second * 2)
}

func Test_TimerStopped(t *testing.T) {
	timer := NewTimer()
	timer.OnTimeOut = func() {
		fmt.Println("This is my timeout function.")
	}
	timer.OnStop = func() {
		fmt.Println("This is my stop function.")
	}

	timer.Start(time.Second * 3)
	time.Sleep(time.Second * 5)
}
