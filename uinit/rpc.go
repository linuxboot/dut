package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"golang.org/x/sys/unix"
)

type RPCRes struct {
	C   []byte
	Err string
}

type Command int

type RPCWelcome struct {
}

func (*Command) Welcome(args *RPCWelcome, r *RPCRes) error {
	r.C = []byte(welcome)
	r.Err = ""
	log.Printf("welcome")
	return nil
}

type RPCExit struct {
	When time.Duration
}

func (*Command) Die(args *RPCExit, r *RPCRes) error {
	go func() {
		time.Sleep(args.When)
		log.Printf("die exits")
		os.Exit(0)
	}()
	*r = RPCRes{}
	log.Printf("die returns")
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
	log.Printf("reboot returns")
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
	*r = RPCRes{Err: "Not yet"}
	log.Printf("kexec returns")
	return nil
}

type RPCSsh struct {
}

func (*Command) Ssh(args *RPCSsh, r *RPCRes) error {
	res := make(chan error)
	go func() {
		c := exec.Command("/bbin/sshd")
		err := c.Start()
		res <- err
	}()
	err := <-res
	*r = RPCRes{Err: fmt.Sprintf("%v", err)}
	log.Printf("sshd returns")
	return nil
}
