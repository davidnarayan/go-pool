go-pool
=======

go-pool implements a basic worker pool for Go


Installation
------------

```
go get github.com/davidnarayan/go-pool
````


Quick Start
-----------

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/davidnarayan/go-pool"
)

// 5 URLs to fetch
var urls = []string{
    "https://example.com/1",
    "https://example.com/2",
    "https://example.com/3",
    "https://example.com/4",
    "https://example.com/5",
}

// Worker function to GET an URL
func worker(args ...interface{}) (interface{}, error) {
    url := args[0].(string)
    resp, err := http.Get(url)

    return resp, err
}

// Handle results from workers
func handleResult(job *pool.Job) {
    if job.Error != nil {
        fmt.Println("Error running job: ", job.Error)
    } else {
        fmt.Println(job.Result)
    }
}

func main() {
    // Create a new pool with 4 workers
    mypool := pool.NewPool(4)
    defer close(wp.In)

    // Keep track of the number of jobs
    var numJobs int

    // Add tasks to the pool using the worker function and the list of URLs
    for _, url := range urls {
        go mypool.Add(worker, url)
        numJobs++
    }

    // Print out the results as they become available
    for i := 0; i < numJobs; i++ {
        go handleResult(wp.Result())
    }
}
```
