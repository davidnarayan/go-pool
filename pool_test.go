package pool

import (
	"net/http"
	"testing"
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
	pool := New(2)
	pool.Run()

	for _, url := range urls {
		pool.Add(fetchURL, url)
		m[url] = true
	}

	for {
		if job, ok := pool.GetResult(); ok {
			// log.Printf("%+v\n", job)
			// Should have a result or an error for each url
			url := job.Args[0]

			if _, ok := m[url.(string)]; ok {
				if job.Result == nil && job.Error == nil {
					t.Errorf("No result or error found for job for url %s",
						url)
				}
			} else {
				t.Errorf("No job found for url: %s", url)
			}

		} else {
			break
		}
	}
}
