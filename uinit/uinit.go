package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
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
	flag.Parse()
	up("127.0.0.1/8", "lo")
	up("192.168.0.2/24", "eth0")
	cmd := exec.Command("wget", "http://192.168.0.1/etc/hosts")
	if o, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ip link up failed(%v, %v); continuing", string(o), err)
	}
	c, err := net.Dial("tcp", "192.168.0.1:8086")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if _, err := io.Copy(c, os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	if _, err = io.Copy(os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
