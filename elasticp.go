package elasticp

import (
	"fmt"
	"sync"
)

type ElasticPool[I, O any] struct {
	poolSize int
	workCh   chan WorkUnit[I, O]
	killCh   chan int
	debug    bool

	m sync.RWMutex
}

func New[I, O any]() *ElasticPool[I, O] {
	return &ElasticPool[I, O]{
		poolSize: 0,
		workCh:   make(chan WorkUnit[I, O], 1024),
		killCh:   make(chan int, 1024),
		debug:    true,
	}
}

func (p *ElasticPool[I, O]) Launch(poolSize int) {
	p.Grow(poolSize)
}

func (p *ElasticPool[I, O]) Submit(wu WorkUnit[I, O]) {
	p.workCh <- wu
}

func (p *ElasticPool[I, O]) Size() int {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.poolSize
}

func (p *ElasticPool[I, O]) Grow(incr int) {
	p.m.Lock()
	defer p.m.Unlock()
	p.spawn(p.poolSize, incr)
	p.poolSize += incr
}

func (p *ElasticPool[I, O]) Shrink(decr int) {
	p.m.Lock()
	defer p.m.Unlock()
	newPoolSize := p.poolSize - decr

	for i := p.poolSize - 1; i >= newPoolSize; i-- {
		p.killCh <- i
	}

	p.poolSize -= decr
}

func (p *ElasticPool[I, O]) Close() {
	p.Shrink(p.poolSize)
	close(p.workCh)
	close(p.killCh)
}

func (p *ElasticPool[I, O]) spawn(offset int, n int) {
	for i := range n {
		go func() {
			id := offset + i
			p.log(fmt.Sprintf("gopher %d spawned", id))

			for {
				select {
				case wu, ok := <-p.workCh:
					if !ok {
						p.log("work channel closed")
						return
					}

					wu.Execute()

				case killSignal, ok := <-p.killCh:
					if !ok {
						p.log("kill channel closed")
						return
					}

					if killSignal == id {
						p.log(fmt.Sprintf("gopher %d killed", id))
						return
					}
				}
			}
		}()
	}
}

func (p *ElasticPool[I, O]) log(msg string) {
	if p.debug {
		fmt.Printf("[elasticp] %s\n", msg)
	}
}
