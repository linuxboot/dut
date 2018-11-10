package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

var welcome = `  ______________
< welcome to DUT >
  --------------
         \   ^__^ 
          \  (oo)\_______
             (__)\       )\/\
                 ||----w |
                 ||     ||
`

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
func main() {
	fmt.Print(welcome)
	flag.Parse()
	up("127.0.0.1/8", "lo")
	up("192.168.0.2/24", "eth0")
	cmd := exec.Command("wget", "http://192.168.0.1:80/bzImage")
	if o, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ip link up failed(%v, %v); continuing", string(o), err)
	}
	c, err := net.Dial("tcp", "192.168.0.1:8086")
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
				if err := unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART); err != nil {
					fmt.Fprintf(c, "%v\n", err)
					fmt.Print(err)
				}

			default:
				fmt.Fprintf(c, welcome)
			}
		}
	}()
	if _, err = io.Copy(c, os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
