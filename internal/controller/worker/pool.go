package worker

import (
	"sync"
)

type Task func()

//go:generate mockery --name=PoolI
type PoolI interface {
	Submit(task Task)
	Shutdown()
}

type Pool struct {
	taskQueue chan Task
	wg        sync.WaitGroup
}

func NewWorkerPool(numWorkers, numTask int) *Pool {
	pool := &Pool{
		taskQueue: make(chan Task, numTask),
	}

	pool.wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go pool.worker()
	}

	return pool
}

func (p *Pool) worker() {
	defer p.wg.Done()

	for task := range p.taskQueue {
		task()
	}
}

func (p *Pool) Submit(task Task) {
	p.taskQueue <- task
}

func (p *Pool) Shutdown() {
	close(p.taskQueue)
	p.wg.Wait()
}
