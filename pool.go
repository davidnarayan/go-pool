//-----------------------------------------------------------------------------
// go-pool is a basic worker pool used to manage a large number of jobs across
// a fixed number of "workers" (i.e. goroutines)
//-----------------------------------------------------------------------------

package pool

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"

	//	"github.com/davidnarayan/go-logging"
	//"github.com/davecgh/go-spew/spew"
)

//-----------------------------------------------------------------------------

type Pool struct {
	MaxWorkers int        // Maximum number of workers (i.e. goroutines) to run
	jobs       *list.List // List of jobs that need to be run

	in      chan *Job // Channel for incoming work
	Results chan *Job // Channel for outgoing results
	wg      *sync.WaitGroup
	counter int32
}

type JobFunc func(...interface{}) (interface{}, error)

type Job struct {
	F        JobFunc
	Args     []interface{}
	Result   interface{}
	Error    error
	Interval time.Duration
}

//-----------------------------------------------------------------------------

// Create a new worker pool
func New(workers int) (pool *Pool) {
	var wg sync.WaitGroup

	// Channels for managing the work
	in := make(chan *Job)
	out := make(chan *Job, workers)

	return &Pool{
		MaxWorkers: workers,
		jobs:       list.New(),
		in:         in,
		Results:    out,
		wg:         &wg,
	}
}

// Add jobs to the pool
func (pool *Pool) Add(f JobFunc, args ...interface{}) {
	job := &Job{
		F:    f,
		Args: args,
	}

	atomic.AddInt32(&pool.counter, 1)
	pool.in <- job
}

func (pool *Pool) Run() {
	// Launch the workers
	for i := 0; i < pool.MaxWorkers; i++ {
		pool.wg.Add(1)
		go worker(pool.in, pool.Results, pool.wg)
	}
}

func (pool *Pool) GetResult() (*Job, bool) {
	if atomic.LoadInt32(&pool.counter) == 0 {
		return nil, false
	}

	result := <-pool.Results
	atomic.AddInt32(&pool.counter, -1)

	return result, true
}

func (pool *Pool) Wait() {
	close(pool.in)
	pool.wg.Wait()
}

func worker(in, out chan *Job, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range in {
		result, err := job.F(job.Args...)
		job.Result = result
		job.Error = err
		out <- job
	}
}

/*

// Run the pool by launching all the workers
func (pool *Pool) Start() {
	if pool.started {
		panic("Pool is already started!")
	}

	// Launch the workers
	for i := 0; i < pool.MaxWorkers; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	pool.started = true
}

// Stop the pool (waiting on any workers to finish)
func (pool *Pool) Stop() {
	if !pool.started {
		panic("Pool is already stopped!")
	}

    logging.Trace("1. ===> %+v", pool.wg)
	pool.wg.Wait()
    pool.started = false
    logging.Trace("%+v", pool.wg)
    close(pool.Out)
    close(pool.In)
}

func (pool *Pool) WaitForJob() (*Job) {
    logging.Trace("POOL COUNTER IS %d", pool.counter)

    if pool.counter == 0 {
        return nil
    }

    pool.counter--

    return <-pool.Out
}

// Run jobs once. This function is useful when the "supervisor" of the process
// restarts frequently and the job interval may not be met.
/*
func (pool *Pool) RunJobsOnce(delay time.Duration) {
	var seen = make(map[string]bool)

	for _, job := range pool.Jobs {
		if _, ok := seen[job.Name]; ok {
			continue
		}

		go func(d time.Duration) {
			logging.Info("Scheduling job: name=%s delay=%s", job.Name, d)

			for {
				select {
				case <-time.After(d):
					pool.In <- job
					return
				}
			}
		}(delay)

		seen[job.Name] = true
	}
}
*/
/*

// Run puts all jobs into the incoming queue
func (pool *Pool) Run() {
	for _, job := range pool.Jobs {
        logging.Trace("adding job: %+v", job)
		pool.In <- job
	}
}

// Run puts jobs into the incoming queue for the specified interval
func (pool *Pool) RunAt(interval time.Duration) {
	for _, job := range pool.Jobs {
		if job.Interval == interval {
			logging.Trace("Matched interval job.Interval=%s interval=%s",
				job.Interval, interval)
			pool.In <- job
		} else {
			logging.Trace("Did not match interval job.Interval=%s interval=%s",
				job.Interval, interval)
		}
	}
}

// worker runs jobs in the incoming queue
func (pool *Pool) worker() {
    //defer pool.wg.Done()

	for job := range pool.In {
		result, err := job.F(job.Args...)
		job.Result = result
		job.Error = err
		pool.Out <- job
        pool.wg.Done()
	}
}
*/
