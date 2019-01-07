// This is a very simple dut program. It builds into one binary to implement
// both client and server. It's just easier to see both sides of the code and test
// that way.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

var (
	host      = flag.String("h", "192.168.0.1", "hostname")
	me        = flag.String("me", "192.168.0.2", "dut hostname")
	port      = flag.String("p", "8080", "port number")
	dir       = flag.String("d", ".", "directory to serve")
	mode      = flag.String("m", "device", "what mode to run in -- device, tester, or ssh starter")
	configNet = flag.Bool("C", true, "configure the network")
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

func dutRPC(host, port string) error {
	l, err := dutStart("tcp", host, port)
	if err != nil {
		return err
	}
	c, err := dutAccept(l)
	if err != nil {
		return err
	}
	cl := rpc.NewClient(c)
	for _, cmd := range []struct {
		call string
		args interface{}
	}{
		{"Command.Welcome", &RPCWelcome{}},
		{"Command.Reboot", &RPCReboot{}},
	} {
		var r RPCRes
		if err := cl.Call(cmd.call, cmd.args, &r); err != nil {
			return err
		}
		fmt.Printf("%v(%v): %v\n", cmd.call, cmd.args, string(r.C))
	}

	if c, err = dutAccept(l); err != nil {
		return err
	}
	cl = rpc.NewClient(c)
	var r RPCRes
	if err := cl.Call("Command.Welcome", &RPCWelcome{}, &r); err != nil {
		return err
	}
	fmt.Printf("%v(%v): %v\n", "Command.Welcome", nil, string(r.C))

	return nil
}

func dutssh(host, port string, args ...string) error {
	l, err := dutStart("tcp", host, port)
	if err != nil {
		return err
	}
	c, err := dutAccept(l)
	if err != nil {
		return err
	}
	var r RPCRes
	cl := rpc.NewClient(c)
	err = cl.Call("Command.Ssh", &RPCSsh{args: args}, &r)
	return err
}

func main() {
	flag.Parse()
	var err error
	switch *mode {
	case "tester":
		err = dutRPC(*host, *port)
	case "ssh":
		dutssh(*host, *port, flag.Args()...)
	case "device":
		err = uinit(*host, *me, *port)
	}
	log.Printf("We are now done ......................")
	if err != nil {
		log.Printf("%v", err)
		os.Exit(2)
	}
}
