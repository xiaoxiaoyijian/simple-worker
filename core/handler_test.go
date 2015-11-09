package core

import (
	"fmt"
	"testing"
	"time"
)

func Test_Handler(t *testing.T) {
	p := func(input interface{}) (output interface{}, err error) {
		fmt.Printf("input: %v\n", input)
		return "xxxx", nil
	}

	h := NewHandler(p)
	fmt.Printf("handler id:%s\n", h.Id())

	ret, err := h.Run("dddddddddddddd")
	fmt.Printf("ret:%v, err:%v\n", ret, err)
}

func Test_NewHandlerPool(t *testing.T) {
	p := func(input interface{}) (output interface{}, err error) {
		fmt.Printf("input: %v\n", input)
		return "xxxx", nil
	}

	pool := NewHandlerPool(p, 10, 10)
	n1, n2 := pool.Status()
	fmt.Printf("active: %d, idle:%d\n", n1, n2)
}

func Test_Pool_Get(t *testing.T) {
	p := func(input interface{}) (output interface{}, err error) {
		fmt.Printf("input: %v\n", input)
		return "xxxx", nil
	}

	pool := NewHandlerPool(p, 10, 10)
	pool.Get()
	pool.Get()
	n1, n2 := pool.Status()
	fmt.Printf("active: %d, idle:%d\n", n1, n2)
}

func Test_Pool_ReleaseHandler(t *testing.T) {
	p := func(input interface{}) (output interface{}, err error) {
		fmt.Printf("input: %v\n", input)
		return "xxxx", nil
	}

	pool := NewHandlerPool(p, 10, 10)
	h, _ := pool.Get()
	pool.Get()
	n1, n2 := pool.Status()
	fmt.Printf("active: %d, idle:%d\n", n1, n2)

	pool.ReleaseHandler(h)
	n1, n2 = pool.Status()
	fmt.Printf("active: %d, idle:%d\n", n1, n2)

	time.Sleep(time.Second * 10)
	n1, n2 = pool.Status()
	fmt.Printf("active: %d, idle:%d\n", n1, n2)
}
