package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/davidnarayan/go-logging"
	"github.com/davidnarayan/go-pool"

	//"github.com/davecgh/go-spew/spew"
)

//-----------------------------------------------------------------------------
// Flags

var host = flag.String("host", "localhost:9000", "HTTP server host:port")
var requests = flag.Int("requests", 100, "Number of requests to send")
var workers = flag.Int("workers", 5, "Number of workers to use")

//-----------------------------------------------------------------------------

// Worker function to GET an URL
func worker(args ...interface{}) (interface{}, error) {
	url := args[0].(string)
	log.Printf("Sending request %s", url)
	resp, err := http.Get(url)

	return resp, err
}

func main() {
	flag.Parse()
	logging.SetLevel(logging.TRACE)
	var urls []string

	for i := 1; i <= *requests; i++ {
		urls = append(urls, fmt.Sprintf("http://%s/%d", *host, i))
	}

	mypool := pool.NewPool(*workers)
	logging.Trace("Sending %d jobs to pool", len(urls))

	for _, url := range urls {
		mypool.Add(worker, url)
	}

	logging.Trace("Reading results from pool")

	// Print out the results as they become available
	for i := 0; i < len(urls); i++ {
		job := mypool.Result()

		if job.Error != nil {
			log.Println("Error running job: ", job.Error)
		} else {
			log.Println(job.Result)
		}

		logging.Info("Pool stats: %+v", mypool.Stats)
	}
}
