package pool

import (
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

var urls = []string{
	"https://google.com",
	"https://yahoo.com",
	"https://httpbin.org",
	"https://golang.org",
}

func fetchURL(args ...interface{}) (interface{}, error) {
	url := args[0].(string)
	resp, err := http.Get(url)

	return resp, err
}

func TestNewPool(t *testing.T) {
	m := make(map[string]bool)
	pool := NewPool(2)

	for _, url := range urls {
		pool.Add(fetchURL, url)
		m[url] = true
	}

	for i := 0; i < len(urls); i++ {
		job := pool.Result()
		url := job.Args[0]

		if _, ok := m[url.(string)]; ok {
			if job.Result == nil && job.Error == nil {
				t.Errorf("No result or error found for job for url %s",
					url)
			}
		} else {
			t.Errorf("No job found for url: %s", url)
		}

	}
}

func TestJobId(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	nextId := func() string {
		var n int

		for {
			n = rand.Intn(10)

			if n > 0 {
				break
			}
		}

		return strconv.Itoa(n)
	}

	pool := NewPool(2)
	pool.SetIdFunc(nextId)

	for _, url := range urls {
		pool.Add(fetchURL, url)
	}

	for i := 0; i < len(urls); i++ {
		job := pool.Result()
		id, err := strconv.Atoi(job.Id)
		if err != nil {
			t.Errorf("Unable to convert job id to integer: %s", err)
		}
		if id < 1 {
			t.Errorf("Invalid Id found in Job: %s", job.Id)
		}
	}
}
