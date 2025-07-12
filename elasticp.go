package elasticp

import (
	"sync"
)

type Task struct {
	Input  []float64
	Output []float64
	Wg     *sync.WaitGroup
}

type Executor interface {
	Execute(Task)
}

type Worker struct {
	taskCh chan Task
}

type Pool struct {
	workers []Worker
	c       uint
	m       *sync.RWMutex
}

func New() *Pool {
	var rw sync.RWMutex

	return &Pool{
		workers: make([]Worker, 0, 32),
		c:       0,
		m:       &rw,
	}
}

func (p *Pool) Grow(amount int) {
	p.m.Lock()
	defer p.m.Unlock()
	for range amount {
		taskCh := make(chan Task, 1024)

		go func() {
			for {
				select {
				case task, ok := <-taskCh:
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

		p.workers = append(p.workers, Worker{
			taskCh: taskCh,
		})
	}
}

func (p *Pool) Shrink(amount int) {
	p.m.Lock()
	defer p.m.Unlock()

	nSize := max(len(p.workers)-amount, 0)
	p.workers = append(make([]Worker, nSize), p.workers[:nSize]...)
}

func (p *Pool) Size() int {
	p.m.RLock()
	defer p.m.RUnlock()
	return len(p.workers)
}

func (p *Pool) Submit(task Task) {
	wLen := len(p.workers)
	for i := 0; i < wLen; i++ {
		p.c++
		i := int(p.c) % wLen
		select {
		case p.workers[i].taskCh <- task:
			// scheduled
			return
		default:
			// let it go through
		}
	}
	p.workers[int(p.c)%wLen].taskCh <- task
}
