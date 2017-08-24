package pool

import "sync"

// Task interface that requires Execute() method
type Task interface {
	Execute()
}

// Pool pool
type Pool struct {
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

// DefaultPoolTaskChannelSize default size of pool
const DEFAULT_POOL_TASK_CHANEL_SIZE = 128

// NewPool create new pool of given size
func NewPool(size int) *Pool {
	pool := &Pool{
		tasks: make(chan Task, DEFAULT_POOL_TASK_CHANEL_SIZE),
		kill:  make(chan struct{}),
	}
	pool.Resize(size)
	return pool
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok { // if chanel is closed, die
				return
			}
			task.Execute()
		case <-p.kill: // it's time to die!
			return
		}
	}
}

// Resize resize pool to specified size
func (p *Pool) Resize(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}
	for p.size > n {
		p.size--
		p.kill <- struct{}{}
	}
}

// Close close pool
func (p *Pool) Close() {
	close(p.tasks)
}

// Wait wait for all tasks in pool
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Exec schedule Task for execution
func (p *Pool) Exec(e Task) {
	p.tasks <- e
}
