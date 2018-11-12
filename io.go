package main

import (
	"net"
)

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

