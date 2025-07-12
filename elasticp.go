package elasticp

import (
	"sync"
)

type Task struct {
	Input  []float64
	Output []float64
	Wg     *sync.WaitGroup
}

type Pool struct {
	workers []chan Task
	c       uint
	m       *sync.RWMutex
}

func New() *Pool {
	var rw sync.RWMutex

	return &Pool{
		workers: make([]chan Task, 0, 32),
		c:       0,
		m:       &rw,
	}
}

func (p *Pool) Grow(amount int) {
	p.m.Lock()
	defer p.m.Unlock()
	for range amount {
		worker := make(chan Task, 1024)

		go func() {
			for {
				select {
				case task, ok := <-worker:
					if !ok {
						return
					}

					for i, v := range task.Input {
						task.Output[i] += v
					}
					task.Wg.Done()
				}
			}
		}()

		p.workers = append(p.workers, worker)
	}
}

func (p *Pool) Shrink(amount int) {
	p.m.Lock()
	defer p.m.Unlock()

	oSize := len(p.workers)
	nSize := max(oSize-amount, 0)
	for i := nSize; i < oSize; i++ {
		close(p.workers[i])
	}

	p.workers = append(make([]chan Task, nSize), p.workers[:nSize]...)
}

func (p *Pool) Size() int {
	p.m.RLock()
	defer p.m.RUnlock()
	return len(p.workers)
}

func (p *Pool) Submit(task Task) {
	for i := 0; i < len(p.workers); i++ {
		select {
		case p.nextWorker() <- task:
			// scheduled
			return
		default:
			// let it go through
		}
	}
	p.nextWorker() <- task
}

func (p *Pool) nextWorker() chan Task {
	p.c++
	return p.workers[int(p.c)%len(p.workers)]
}
