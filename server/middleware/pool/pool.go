package pool

import "container/list"

type Pool struct {
	pool list.List
	fac  ObjFactory
}

type ObjFactory func() interface{}

func (p *Pool) Rent() interface{} {
	if p.pool.Len() == 0 {
		return p.fac()
	}
	return p.pool.Remove(p.pool.Back())
}

func (p *Pool) Return(e interface{}) {
	p.pool.PushBack(e)
}
func (p *Pool) Init(fac ObjFactory) {
	p.pool.Init()
	p.fac = fac
}
