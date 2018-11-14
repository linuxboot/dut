package main

import (
	"log"
	"os"
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
	*r = RPCRes{Err: "Not yet"}
	return nil
}
