//-----------------------------------------------------------------------------
// go-pool is a basic worker pool used to manage a large number of jobs across
// a fixed number of "workers" (i.e. goroutines)
//-----------------------------------------------------------------------------

package pool

import (
	"sync"
)

//-----------------------------------------------------------------------------

// JobFunc represents the user-defined function to do some type of work.
type JobFunc func(...interface{}) (interface{}, error)

// Job holds everything required for a single unit of work
type Job struct {
	F      JobFunc
	Args   []interface{}
	Result interface{}
	Error  error
}

// Pool manages and distributes work to the workers
type Pool struct {
	in  chan *Job
	out chan *Job

	wg *sync.WaitGroup
}

// NewPool creates a new pool with a fixed number of workers
func NewPool(size int) *Pool {
	var wg sync.WaitGroup
	p := &Pool{
		wg:  &wg,
		in:  make(chan *Job),
		out: make(chan *Job),
	}

	for i := 0; i < size; i++ {
		p.wg.Add(1)
		go p.worker()
	}

	return p
}

// Add submits creates a new job for the pool
func (p *Pool) Add(fn JobFunc, args ...interface{}) {
	go func() {
		job := &Job{
			F:    fn,
			Args: args,
		}

		p.in <- job
	}()
}

// Return one result
func (p *Pool) Result() *Job {
	return <-p.out
}

// worker runs the next job available
func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		job := <-p.in
		job.Result, job.Error = job.F(job.Args...)
		p.out <- job
	}
}
