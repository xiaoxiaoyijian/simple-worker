package core

import (
	"errors"
	"sync"
	"time"
)

var (
	ERR_POOL_DUPLICATE = errors.New("Already exist in pool.")
	ERR_POOL_MAXINUM   = errors.New("Reach max num of pool.")
	ERR_POOL_NOTFOUND  = errors.New("Not exist in pool.")
)

type Pool struct {
	values map[string]interface{}
	timers map[string]*Timer

	maxNum int

	m sync.RWMutex
}

func NewPool(max int) *Pool {
	return &Pool{
		values: make(map[string]interface{}, max),
		timers: make(map[string]*Timer, max),
		maxNum: max,
	}
}

func (this *Pool) Add(id string, v interface{}, timeout time.Duration) error {
	this.m.Lock()
	defer this.m.Unlock()

	if _, ok := this.values[id]; ok {
		return ERR_POOL_DUPLICATE
	}

	if this.maxNum > 0 && len(this.values) >= this.maxNum {
		return ERR_POOL_MAXINUM
	}

	this.values[id] = v

	if timeout > 0 {
		timer := NewTimer()
		timer.OnTimeOut = func() {
			this.Remove(id)
		}

		timer.Start(timeout)
		this.timers[id] = timer
	}

	return nil
}

func (this *Pool) Remove(id string) error {
	this.m.Lock()
	defer this.m.Unlock()

	if _, ok := this.values[id]; !ok {
		return ERR_POOL_NOTFOUND
	}

	delete(this.values, id)

	if t, ok := this.timers[id]; ok {
		t.Stop()
		delete(this.timers, id)
	}

	return nil
}

func (this *Pool) Get(id string) (v interface{}, hasTimer bool, err error) {
	this.m.RLock()
	defer this.m.RUnlock()

	var ok bool
	v, ok = this.values[id]
	if !ok {
		err = ERR_POOL_NOTFOUND
		return
	}

	if _, ok = this.timers[id]; ok {
		hasTimer = true
	}

	return
}

func (this *Pool) Len() (int, int) {
	this.m.RLock()
	defer this.m.RUnlock()

	return len(this.values), len(this.timers)
}

func (this *Pool) Max() int {
	return this.maxNum
}

func (this *Pool) Pop() (string, interface{}) {
	this.m.RLock()
	defer this.m.RUnlock()

	for k, v := range this.values {
		if t, ok := this.timers[k]; ok {
			t.Stop()
		}

		return k, v
	}

	return "", nil
}
