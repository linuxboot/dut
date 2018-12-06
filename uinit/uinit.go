package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
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
	cmd := exec.Command("ip", "addr", "add", ip, dev)
	if o, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ip link failed(%v, %v); continuing", string(o), err)
	}
	cmd = exec.Command("ip", "link", "set", "dev", dev, "up")
	if o, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ip link up failed(%v, %v); continuing", string(o), err)
	}

}

func uinit(t, a string) error {
	h := strings.Split(a, ":")
	log.Printf("here we are in uinit")
	if os.Getuid() == 0 {
		up("127.0.0.1/8", "lo")
		up(h[0]+"/24", "eth0")
	}
	c, err := net.Dial(t, a)
	if err != nil {
		log.Fatal(err)
	}

	var Cmd Command
	s := rpc.NewServer()
	if err := s.Register(&Cmd); err != nil {
		return err
	}
	s.ServeConn(c)
	log.Printf("And uinit is all done.")
	return err

}
