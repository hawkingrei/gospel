package main

import (
	"fmt"
	"crypto/sha1"
	"os"
	"net"
)

func main(){
	var id [20]byte
        h := sha1.New()
	ss := "wasdfadsf"
        h.Write([]byte(ss))
	h.Sum(id[:0:20])
	fmt.Println(id)
	addr, err := net.ResolveUDPAddr("udp","router.bittorrent.com:6881")
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}
	defer conn.Close()
}
