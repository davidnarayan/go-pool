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

func main() {
    // Create a new pool with 4 workers
    mypool := pool.New(4)

    // Launch the workers
    mypool.Run()

    // Add tasks to the pool using the worker function and the list of URLs
    for _, url := range urls {
        mypool.Add(worker, url)
    }

    // Print out the results as they become available
    for {
        if job, ok := mypool.GetResult(); ok {
            if job.Error != nil {
                fmt.Println("Error running job: ", job.Error)
            } else {
                fmt.Println(job.Result)
            }
        } else {
            // No more results - exit loop!
            break
        }
    }
}
```
