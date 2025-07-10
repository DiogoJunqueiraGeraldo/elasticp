package elasticp

import (
	"sync/atomic"
)

type ElasticPool struct {
	workers    []chan func()
	nextWorker *atomic.Uint32
	buff       int
}

func New(size int, buff int) *ElasticPool {
	return &ElasticPool{
		workers:    make([]chan func(), size),
		nextWorker: &atomic.Uint32{},
		buff:       buff,
	}
}

func (p *ElasticPool) Start() {
	for i := 0; i < len(p.workers); i++ {
		workCh := make(chan func(), p.buff)
		p.workers[i] = workCh

		go func(ch <-chan func()) {
			for f := range ch {
				f()
			}
		}(workCh)
	}
}

func (p *ElasticPool) Go(f func()) {
	workerCount := len(p.workers)
	nextWorker := p.nextWorker.Add(1)
	i := int(nextWorker) % workerCount
	p.workers[i] <- f
}
