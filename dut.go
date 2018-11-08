package main

import (
	"fmt"
)

import (
	"flag"
	"log"
	"net/http"
)

var (
	host = flag.String("h", "127.0.0.1", "hostname")
	port = flag.String("p", "8080", "port number")
	dir  = flag.String("d", ".", "directory to serve")
)

var cacheHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func maxAgeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for _, v := range cacheHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		h.ServeHTTP(w, r)
	})
}

func srvfiles() {
	http.Handle("/", maxAgeHandler(http.FileServer(http.Dir(*dir))))
	log.Fatal(http.ListenAndServe(*host+":"+*port, nil))
}

func main() {
	flag.Parse()
}

