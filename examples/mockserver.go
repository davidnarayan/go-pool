// mockserver is a simple HTTP server used to test pool

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//-----------------------------------------------------------------------------
// Flags

var listen = flag.String("listen", "127.0.0.1:9000", "Listen on this address")

//-----------------------------------------------------------------------------

var counter int64

func handler(w http.ResponseWriter, r *http.Request) {
	counter++
	c := counter
	log.Printf("Received request %d from %s: %s", c, r.RemoteAddr, r.URL)

	var delay time.Duration

	if c%3 == 0 {
		delay = time.Duration(rand.Intn(15)) * time.Second
	} else {
		delay = time.Duration(1) * time.Nanosecond
	}

	select {
	case <-time.After(delay):
		fmt.Fprintf(w, "mockserver OK %s %s\n", r.URL, delay)
		log.Printf("Sent response %d to %s for %s after %s", c, r.RemoteAddr,
			r.URL, delay)
	}

}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
	http.HandleFunc("/", handler)
	log.Println("Listening on", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
