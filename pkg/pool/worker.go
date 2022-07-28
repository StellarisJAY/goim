package pool

import (
	"context"
	"sync"
)

type Worker struct {
	tasks chan func()
}

type WorkerPool struct {
	workers    []*Worker
	maxWorker  int
	close      context.Context
	cancel     context.CancelFunc
	nextWorker int
	sync.Mutex
}

func NewWorkerPool(workers int) *WorkerPool {
	pool := &WorkerPool{
		workers:   make([]*Worker, workers),
		maxWorker: workers,
		Mutex:     sync.Mutex{},
	}
	pool.close, pool.cancel = context.WithCancel(context.Background())
	for i := 0; i < workers; i++ {
		tasks := make(chan func(), 1204)
		pool.workers[i] = &Worker{tasks: tasks}
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case task := <-tasks:
					task()
				}
			}
		}(pool.close)
	}
	return pool
}

func (p *WorkerPool) Submit(task func()) {
	p.Lock()
	p.nextWorker = p.nextWorker + 1
	if p.nextWorker == p.maxWorker {
		p.nextWorker = 0
	}
	worker := p.workers[p.nextWorker]
	p.Unlock()
	worker.tasks <- task
}
