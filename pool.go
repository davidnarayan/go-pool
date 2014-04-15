//-----------------------------------------------------------------------------
// go-pool is a basic worker pool used to manage a large number of jobs across
// a fixed number of "workers" (i.e. goroutines)
//-----------------------------------------------------------------------------

package pool

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

//-----------------------------------------------------------------------------

// JobFunc represents the user-defined function to do some type of work.
type JobFunc func(...interface{}) (interface{}, error)

// Job holds everything required for a single unit of work
type Job struct {
	Id     string
	F      JobFunc
	Args   []interface{}
	Result interface{}
	Error  error
}

// IdFunc represents a user-defined function to assign Job IDs
type IdFunc func() string

// Pool manages and distributes work to the workers
type Pool struct {
	In  chan *Job
	Out chan *Job

	wg *sync.WaitGroup

	nextId IdFunc

	Stats struct {
		Submitted int64
		Pending   int64
		Running   int64
		Completed int64
	}
}

// NewPool creates a new pool with a fixed number of workers
func NewPool(size int) *Pool {
	var wg sync.WaitGroup
	p := &Pool{
		wg:  &wg,
		In:  make(chan *Job),
		Out: make(chan *Job),
	}

	for i := 0; i < size; i++ {
		p.wg.Add(1)
		go p.worker()
	}

	go func() {
		p.wg.Wait()
		close(p.Out)
	}()

	return p
}

// SetIdFunc adds an id generator to the pool
func (p *Pool) SetIdFunc(fn IdFunc) {
	p.nextId = fn
}

// Add submits creates a new job for the pool
func (p *Pool) Add(fn JobFunc, args ...interface{}) {
	go func() {
		job := &Job{
			F:    fn,
			Args: args,
		}

		if p.nextId != nil {
			job.Id = p.nextId()
		}

		atomic.AddInt64(&p.Stats.Submitted, 1)
		atomic.AddInt64(&p.Stats.Pending, 1)
		p.In <- job
	}()
}

// Return one result
func (p *Pool) Result() *Job {
	return <-p.Out
}

// worker runs the next job available
func (p *Pool) worker() {
	defer p.wg.Done()
	for job := range p.In {
		atomic.AddInt64(&p.Stats.Pending, -1)
		atomic.AddInt64(&p.Stats.Running, 1)
		job.Result, job.Error = job.F(job.Args...)
		atomic.AddInt64(&p.Stats.Running, -1)
		atomic.AddInt64(&p.Stats.Completed, 1)
		p.Out <- job
	}
}

// String provides statistics about the pool and its jobs
func (p *Pool) String() string {
	s := fmt.Sprintf("%+v", p.Stats)
	s = strings.Replace(s, ":", "=", -1)
	s = strings.Trim(s, "{}")
	return s
}
