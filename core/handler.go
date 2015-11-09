package core

import (
	"errors"
	"github.com/nu7hatch/gouuid"
	"sync"
	"time"
)

type Handler struct {
	p  Processer
	id string
}

func NewHandler(p Processer) *Handler {
	uuid_val, _ := uuid.NewV4()

	return &Handler{
		p:  p,
		id: uuid_val.String(),
	}
}

func (this *Handler) Id() string {
	return this.id
}

func (this *Handler) Run(input interface{}) (output interface{}, err error) {
	return this.p(input)
}

const (
	IDLE_TIMEOUT = time.Second * 5
)

type HandlerPool struct {
	p          Processer
	activePool *Pool
	idlePool   *Pool

	m sync.RWMutex
}

func NewHandlerPool(p Processer, maxIdle, maxActive int) *HandlerPool {
	return &HandlerPool{
		p:          p,
		activePool: NewPool(maxActive),
		idlePool:   NewPool(maxIdle),
	}
}

func (this *HandlerPool) Get() (*Handler, error) {
	this.m.Lock()
	defer this.m.Unlock()

	activeLen, _ := this.activePool.Len()
	if activeLen >= this.activePool.Max() {
		return nil, errors.New("Reach maxinum active handlers.")
	}

	idleLen, _ := this.idlePool.Len()
	if idleLen > 0 {
		id, v := this.idlePool.Pop()
		if h, ok := v.(*Handler); ok {
			this.activePool.Add(id, h, 0)
			return h, nil
		}
	}

	h := NewHandler(this.p)
	this.activePool.Add(h.Id(), h, 0)

	return h, nil
}

func (this *HandlerPool) ReleaseHandler(h *Handler) error {
	this.m.Lock()
	defer this.m.Unlock()

	v, _, _ := this.activePool.Get(h.Id())
	if v == nil {
		return errors.New("Active handler not found.")
	}

	this.activePool.Remove(h.Id())
	idleLen, _ := this.idlePool.Len()
	if idleLen < this.idlePool.Max() {
		return this.idlePool.Add(h.Id(), h, IDLE_TIMEOUT)
	}

	return nil
}

func (this *HandlerPool) Status() (int, int) {
	this.m.RLock()
	defer this.m.RUnlock()

	activeLen, _ := this.activePool.Len()
	idleLen, _ := this.idlePool.Len()

	return activeLen, idleLen
}
