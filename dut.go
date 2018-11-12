package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	host = flag.String("h", "192.168.0.1", "hostname")
	port = flag.String("p", "8080", "port number")
	conv = flag.String("c", "8086", "conv port number")
	dir  = flag.String("d", ".", "directory to serve")
	real = flag.Bool("r", true, "Run a real test")
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

// have a conversation
func con(t, a string) error {
	ln, err := net.Listen(t, a)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Printf("Listening on %v", ln.Addr())

	c, err := ln.Accept()
	if err != nil {
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
func test(t, a string) error {
	ln, err := net.Listen(t, a)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Printf("Listening on %v at %v", ln.Addr(), time.Now())

	if err := ln.(*net.TCPListener).SetDeadline(time.Now().Add(1 * time.Minute)); err != nil {
		return err
	}
	c, err := ln.Accept()
	if err != nil {
		log.Printf("Listen failed: %v at %v", err, time.Now())
		log.Print(err)
		return err
	}
	log.Printf("Accepted %v", c)
	go func() {
		if _, err := io.Copy(os.Stdout, c); err != nil {
			log.Print(err)
			return
		}
	}()

	if _, err := c.Write([]byte("w")); err != nil {
		log.Fatal(err)
	}
	if _, err := c.Write([]byte("r")); err != nil {
		log.Fatal(err)
	}
	// other end reboots; do an accept
	if err := ln.(*net.TCPListener).SetDeadline(time.Now().Add(3 * time.Minute)); err != nil {
		return err
	}
	c, err = ln.Accept()
	if err != nil {
		log.Printf("Listen failed: %v at %v", err, time.Now())
		log.Print(err)
		return err
	}
	log.Printf("Accepted %v", c)
	go func() {
		if _, err := io.Copy(os.Stdout, c); err != nil {
			log.Print(err)
		}
	}()
	if _, err := c.Write([]byte("w")); err != nil {
		log.Fatal(err)
	}
	return nil
}
func main() {
	flag.Parse()
	go srvfiles()
	a := *host + ":" + *conv
	if !*real {
		con("tcp", a)
		os.Exit(0)
	}
	if err := test("tcp", a); err != nil {
		log.Print(err)
		os.Exit(2)
	}
	os.Exit(0)
}
