package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
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

func uinit(r, p string) error {
	log.Printf("here we are in uinit")
	log.Printf("UINIT uid is %d", os.Getuid())

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
