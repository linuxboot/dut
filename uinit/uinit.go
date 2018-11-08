package main

import (
	"flag"
	"log"
	"os/exec"
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
func main() {
	flag.Parse()
	up("127.0.0.1/8", "lo")
	up("192.168.0.2", "eth0")
	cmd := exec.Command("wget", "http://192.168.0.1/etc/hosts")
	if o, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ip link up failed(%v, %v); continuing", string(o), err)
	}
}
