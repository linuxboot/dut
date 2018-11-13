package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
	
	"golang.org/x/sys/unix"

)

type RPCRes struct {
	C   []byte
	Err error
}

type Command int

type RPCCmd struct {
	C []byte
}

func (*Command) Welcome(args *RPCCmd, r *RPCRes) error {
	r.C = []byte(welcome)
	r.Err = nil
	return nil
}

type RPCExit struct {
	When time.Duration
}

func (*Command) Die(args *RPCExit, r *RPCRes) error {
	go func() {
		time.Sleep(args.When)
		os.Exit(0)
	}()
	*r = RPCRes{}
	return nil
}

type RPCReboot struct {
	When time.Duration
}

func (*Command) Reboot(args *RPCReboot, r *RPCRes) error {
	go func() {
		time.Sleep(args.When)
		if err := unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART); err != nil {
			log.Printf("%v\n", err)
		}
	}()
	*r = RPCRes{}
	return nil
}

type RPCKexec struct {
	File string
	When time.Duration
}

func (*Command) Kexec(args *RPCReboot, r *RPCRes) error {
	go func() {
		time.Sleep(args.When)
	}()
	*r = RPCRes{Err: fmt.Errorf("Not yet"),}
	return nil
}

func dutIO(c net.Conn, b []byte, r []byte) error {
	if _, err := c.Write(b); err != nil {
		return err
	}
	if r != nil {
		_, err := c.Read(r)
		return err
	}
	return nil
}
