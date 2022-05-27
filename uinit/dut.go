// This is a very simple dut program. It builds into one binary to implement
// both client and server. It's just easier to see both sides of the code and test
// that way.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"

	"github.com/u-root/u-root/pkg/ulog"
)

var (
	debug = flag.Bool("debug", false, "Enable debug prints")
	host  = flag.String("h", "192.168.0.1", "hostname")
	klog  = flag.Bool("klog", false, "Direct all logging to klog -- depends on debug")
	port  = flag.String("p", "8080", "port number")
	dir   = flag.String("d", ".", "directory to serve")
	mode  = flag.String("m", "device", "what mode to run in -- device, tester, or ssh starter")

	// for debug
	v = func(string, ...interface{}) {}
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

func dutcpu(host, port, pubkey, hostkey, cpuport string) error {
	var req = &RPCCPU{Port: cpuport}
	var err error

	// we send the pubkey and hostkey as the value of the key, not the
	// name of the file.
	// TODO: maybe use ssh_config to find keys? the cpu client can do that.
	// Note: the public key is not optional. That said, we do not test
	// for len(*pubKey) > 0; if it is set to ""< ReadFile will return
	// an error.
	if req.PubKey, err = ioutil.ReadFile(pubkey); err != nil {
		return fmt.Errorf("Reading pubKey:%w", err)
	}
	if len(hostkey) > 0 {
		if req.HostKey, err = ioutil.ReadFile(hostkey); err != nil {
			return fmt.Errorf("Reading hostKey:%w", err)
		}
	}

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
		{"Command.Welcome", &RPCWelcome{}},
		{"Command.CPU", req},
	} {
		var r RPCRes
		if err := cl.Call(cmd.call, cmd.args, &r); err != nil {
			return err
		}
		fmt.Printf("%v(%v): %v\n", cmd.call, cmd.args, string(r.C))
	}
	return err
}

func main() {
	// for CPU
	flag.Parse()

	if *debug {
		v = log.Printf
		if *klog {
			ulog.KernelLog.Reinit()
			v = ulog.KernelLog.Printf
		}
	}
	var err error
	switch *mode {
	case "tester":
		err = dutRPC(*host, *port)
	case "cpu":
		var (
			pubKey  = flag.String("pubkey", "key.pub", "public key file")
			hostKey = flag.String("hostkey", "", "host key file -- usually empty")
			cpuPort = flag.String("cpuport", "17010", "cpu port -- IANA value is ncpu tcp/17010")
		)
		flag.Parse()
		dutcpu(*host, *port, *pubKey, *hostKey, *cpuPort)
	case "device":
		err = uinit(*host, *port)
	}
	log.Printf("We are now done ......................")
	if err != nil {
		log.Printf("%v", err)
		os.Exit(2)
	}
}
