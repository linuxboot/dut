package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"time"
)

var (
	rebooting = "Rebooting!"
	welcome   = `  ______________
< welcome to DUT >
  --------------
         \   ^__^ 
          \  (oo)\_______
             (__)\       )\/\
                 ||----w |
                 ||     ||
`
)

func up(ip, dev string) {
	cmd := exec.Command("ip", "link", "set", "dev", dev, "up")
	if o, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ip link up failed(%q, %v); continuing", o, err)
	}
	cmd = exec.Command("ip", "addr", "add", ip, dev)
	if o, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ip addr add failed(%q, %v); continuing", o, err)
	}
	cmd = exec.Command("ip", "addr")
	if o, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ip addr failed%q, (%v); continuing", o, err)
	}
	log.Printf("Sleeping 16 seconds for stupid network to come up")
	time.Sleep(16 * time.Second)
	log.Printf("up all done")
}

func uinit(r, l, p string) error {
	log.Printf("here we are in uinit")
	log.Printf("UINIT uid is %d", os.Getuid())
	if os.Getuid() == 0 && *configNet {
		//go up("127.0.0.1/8", "lo")
		up(l+"/24", "eth0")
	}

	na := r + ":" + p
	log.Printf("Now dial %v", na)
	c, err := net.Dial("tcp", na)
	if err != nil {
		log.Printf("Dial went poorly")
		return err
	}
	log.Printf("Start the RPC server")
	var Cmd Command
	s := rpc.NewServer()
	log.Printf("rpc server is %v", s)
	if err := s.Register(&Cmd); err != nil {
		log.Printf("register failed: %v", err)
		return err
	}
	log.Printf("Serve and protect")
	s.ServeConn(c)
	log.Printf("And uinit is all done.")
	return err

}
