package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	rebooting = "Rebooting!"
	welcome = `  ______________
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
func uinit(t, a string) error{
	h := strings.Split(a, ":")
	if os.Getuid() == 0 {
		up("127.0.0.1/8", "lo")
		up(h[0]+"/24", "eth0")
	}
	c, err := net.Dial(t, a)
	if err != nil {
		log.Fatal(err)
	}
	c.Write([]byte(welcome))
	go func() {
		var nerr int
		var b = make([]byte, 1)
		for {
			if _, err := c.Read(b); err != nil {
				fmt.Print(err)
				if nerr > 128 {
					return
				}
				nerr++
			}
			os.Stdout.Write(b)
			switch b[0] {
			case 'e':
				os.Exit(0)
			case 'k':
				fmt.Println("kexec: not yet")
			case 'a':
				c.Write([]byte(welcome))
			case 'r':
				c.Write([]byte(rebooting))
				// well, this better never run in a test ...
				if err := unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART); err != nil {
					fmt.Fprintf(c, "%v\n", err)
					fmt.Print(err)
				}
				c.Close()
				return
			default:
				fmt.Fprintf(c, welcome)
			}
		}
	}()
	if _, err = io.Copy(c, os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return err

}
