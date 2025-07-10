package elasticp

import "time"

type WorkerStatus int

const (
	WorkerStatusIdle WorkerStatus = iota
	WorkerStatusBusy
	WorkerStatusFired
)

type Worker struct {
	status      WorkerStatus
	lastUpdated time.Time
	hrCh        chan struct{}
}

func (w *Worker) Status() WorkerStatus {
	return w.status
}

func (w *Worker) LastUpdated() time.Time {
	return w.lastUpdated
}

func (w *Worker) Do(work func()) {
	w.status = WorkerStatusBusy
	w.lastUpdated = time.Now()
	work()
	w.status = WorkerStatusIdle
	w.lastUpdated = time.Now()
}
