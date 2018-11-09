package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	host = flag.String("h", "192.168.0.1", "hostname")
	port = flag.String("p", "80", "port number")
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

func con(t, a string) error {
	ln, err := net.Listen(t, a)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Printf("Listening on %v", ln.Addr())
	
	c, err := ln.Accept()
	if err != nil {
		log.Print(err)
		return err
	}
	log.Printf("Accepted %v", c)
	go func() {
		if _, err := io.Copy(c, os.Stdin); err != nil {
			log.Print(err)
		}
	}()
	if _, err = io.Copy(os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return nil
}
func main() {
	flag.Parse()
	go srvfiles()
	con("tcp", "192.168.0.1:8086")
}
