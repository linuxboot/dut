package main

import (
	"net"
)

type RPCCmd struct {
	C []byte
}

type RPCRes struct {
	C   []byte
	Err error
}

type Command int

func (*Command) Welcome(args *RPCCmd, r *RPCRes) error {
	r.C = []byte(welcome)
	r.Err = nil
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
