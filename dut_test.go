package main

import (
	"net/rpc"
	"testing"
	"time"
)

func TestUinit(t *testing.T) {
	l, err := dutStart("tcp", "localhost", "")
	if err != nil {
		t.Fatal(err)
	}

	a := l.Addr()
	t.Logf("listening on %v", a)
	// Kick off our node.
	go func() {
		time.Sleep(1)
		if err := uinit(a.Network(), a.String()); err != nil {
			t.Fatalf("starting uinit: got %v, want nil", err)
		}
	}()

	c, err := dutAccept(l)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Connected on %v", c)

	cl := rpc.NewClient(c)
	var b = make([]byte, len(welcome))
	// issue a command
	b = make([]byte, len(welcome))
	call := &RPCCmd{C: b}
	var r RPCRes
	if err = cl.Call("Command.Welcome", call, &r); err != nil {
		t.Fatal(err)
	}
	if r.Err != nil {
		t.Fatal(err)
	}
	t.Logf("welcome comand? %v", string(r.C))
	if string(r.C) != welcome {
		t.Errorf("welcome: got %s, want %s", string(b), welcome)
	}
	if false {
		b = make([]byte, len(welcome))
		if err := dutIO(c, []byte("r"), b); err != nil {
			t.Error(err)
		}
		r := string(b[:len(rebooting)])
		t.Logf("welcome? %v", r)
		if r != rebooting {
			t.Errorf("rebooting: got %q, want %q", r, rebooting)
		}
	}

}
