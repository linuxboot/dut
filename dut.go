// This is a very simple dut program. It builds into one binary to implement
// both client and server. It's just easier to see both sides of the code and test
// that way.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

var (
	host   = flag.String("h", "192.168.0.1", "hostname")
	port   = flag.String("p", "8080", "port number")
	dir    = flag.String("d", ".", "directory to serve")
	runDUT = flag.Bool("r", true, "run as the DUT controller")
)

func dutStart(t, host, port string) (net.Listener, error) {
	ln, err := net.Listen(t, host+":"+port)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	log.Printf("Listening on %v at %v", ln.Addr(), time.Now())
	return ln, nil
}

func dutAccept(l net.Listener) (net.Conn, error) {
	if err := l.(*net.TCPListener).SetDeadline(time.Now().Add(3 * time.Minute)); err != nil {
		return nil, err
	}
	c, err := l.Accept()
	if err != nil {
		log.Printf("Listen failed: %v at %v", err, time.Now())
		log.Print(err)
		return nil, err
	}
	log.Printf("Accepted %v", c)
	return c, nil
}

func dutRun(host, port string) error {
	l, err := dutStart("tcp", host, port)
	if err != nil {
		return err
	}
	c, err := dutAccept(l)
	if err != nil {
		return err
	}
	go func() {
		if _, err := io.Copy(os.Stdout, c); err != nil {
			log.Print(err)
			return
		}
	}()
	if err := dutIO(c, []byte("w"), nil); err != nil {
		log.Fatal(err)
	}
	if err := dutIO(c, []byte("r"), nil); err != nil {
		log.Fatal(err)
	}
	// other end reboots; do an accept
	if c, err = dutAccept(l); err != nil {
		return err
	}

	log.Printf("Accepted %v", c)
	go func() {
		if _, err := io.Copy(os.Stdout, c); err != nil {
			log.Print(err)
		}
	}()
	if err := dutIO(c, []byte("w"), nil); err != nil {
		log.Fatal(err)
	}
	return nil
}

func dutRPC(host, port string) error {
	l, err := dutStart("tcp", host, port)
	if err != nil {
		return err
	}
	c, err := dutAccept(l)
	if err != nil {
		return err
	}
	// other end reboots; do an accept
	if c, err = dutAccept(l); err != nil {
		return err
	}
	cl := rpc.NewClient(c)
	log.Printf("Accepted %v", cl)
	for _, cmd := range []*RPCCmd{
		{[]byte("w")},
	} {
		var r RPCRes
		if err := cl.Call("Command", cmd, &r); err != nil {
			return err
		}
		fmt.Printf("%v: %v\n", string(r.C), r.Err)
	}

	return nil
}

func main() {
	flag.Parse()
	var err error
	if *runDUT {
		if false {
			err = dutRun(*host, *port)
		}
		err = dutRPC(*host, *port)
	} else {
		err = uinit(*host, *port)
	}
	if err != nil {
		log.Fatal(err)
	}
}
