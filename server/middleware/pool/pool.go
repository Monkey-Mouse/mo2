package pool

import (
	"container/list"
	"sync"
)

type Pool struct {
	pool list.List
	fac  ObjFactory
}

type ObjFactory func() interface{}

var lock = sync.Mutex{}

func (p *Pool) Rent() interface{} {
	lock.Lock()
	defer lock.Unlock()
	if p.pool.Len() == 0 {
		return p.fac()
	}
	return p.pool.Remove(p.pool.Back())
}

func (p *Pool) Return(e interface{}) {
	lock.Lock()
	defer lock.Unlock()
	p.pool.PushBack(e)
}
func (p *Pool) Init(fac ObjFactory) {
	p.pool.Init()
	p.fac = fac
}
