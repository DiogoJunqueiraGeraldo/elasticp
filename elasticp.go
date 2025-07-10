package elasticp

import (
	"fmt"
	"sync"
	"time"
)

type ElasticPool struct {
	workCh  chan func()
	monitor *time.Ticker
	workers []Worker
	config  ElasticPoolConfiguration

	m sync.RWMutex
}

func New(cfg ElasticPoolConfiguration) *ElasticPool {
	return &ElasticPool{
		workCh:  make(chan func(), cfg.workBuffer),
		monitor: time.NewTicker(cfg.idleTolerance),
		workers: make([]Worker, 0, cfg.maxWorkers),
		config:  cfg,
	}
}

func (p *ElasticPool) Start() {
	p.launchWorkers(p.config.minWorkers)
	go p.monitorWorkers()
}

func (p *ElasticPool) Go(f func()) {
	workAmount := len(p.workCh) + 1
	workCapacity := cap(p.workCh)
	workRatio := float32(workAmount) / float32(workCapacity)

	if workRatio > p.config.workThreshold {
		p.launchWorkers(p.config.growthRatio)
	}

	p.workCh <- f
}

func (p *ElasticPool) launchWorkers(n int) {
	available := p.config.maxWorkers - len(p.workers)
	if available <= 0 {
		p.log("Worker limit reached")
		return
	}

	newWorkers := min(available, n)
	for i := 0; i < newWorkers; i++ {
		worker := Worker{
			status:      WorkerStatusIdle,
			lastUpdated: time.Now(),
			hrCh:        make(chan struct{}),
		}

		go func() {
			for {
				p.log("Starting worker")
				select {
				case work := <-p.workCh:
					worker.Do(work)
					p.log("Worker doing nothing more than was paid to do")

				case <-worker.hrCh:
					p.log("Worker was fired")
					return
				}
			}
		}()

		p.safeHire(worker)
	}
}

func (p *ElasticPool) safeHire(worker Worker) {
	p.m.Lock()
	defer p.m.Unlock()
	p.workers = append(p.workers, worker)
}

func (p *ElasticPool) tryFire(i int) bool {
	p.m.Lock()
	defer p.m.Unlock()

	worker := p.workers[i]
	idleTime := time.Since(worker.lastUpdated)
	if worker.status == WorkerStatusIdle && idleTime > p.config.idleTolerance {
		worker.hrCh <- struct{}{}
		p.workers[i] = p.workers[len(p.workers)-1]
		p.workers = p.workers[:len(p.workers)-1]
		p.log("Lazy worker fired")

		return true
	}

	return false
}

func (p *ElasticPool) monitorWorkers() {
	for range p.monitor.C {
		i := 0
		for i < len(p.workers) {
			if len(p.workers) <= p.config.minWorkers {
				break
			}

			fired := p.tryFire(i)
			if fired {
				continue
			}

			i++
		}
	}
}

func (p *ElasticPool) log(msg string) {
	if p.config.isDebug {
		fmt.Println("[DEBUG] ", msg)
	}
}
