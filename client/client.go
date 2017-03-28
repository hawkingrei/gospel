package main

import (
	"fmt"
	"math/rand"
	"net"
	/*
	"crypto/sha1"
	"os"
	"net"
	"github.com/zeebo/bencode"
	*/
)

func main(){
	id := new([20]byte);
	rand.Read(id[:])
	fmt.println(id)
	var routers = []string{
		"router.magnets.im:6881",
		"router.bittorrent.com:6881",
		"dht.transmissionbt.com:6881",
		"router.utorrent.com:6881",
	}
	for _, addr := range routers {
		addr, err := net.ResolveUDPAddr("udp", addr)
	}
	/*
	var id [20]byte
 	h := sha1.New()
	ss, err := os.Hostname()
    h.Write([]byte(ss))
	h.Sum(id[:0:20])
	targetid := h.Sum(nil)

	fmt.Println(targetid)
	addr, err := net.ResolveUDPAddr("udp", "router.bittorrent.com:6881")
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	d := map[string]interface{}{
		"t": "aa",
		"y": "q",
		"q": "get_peers",
		"a": map[string]interface{}{"info_hash":"abcdefghij0123456789",
				"id":"mnopqrstuvwxyz123456"},
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}
	defer conn.Close()
	data, err := bencode.EncodeBytes(d)
	fmt.Println(string(data))
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("failed:", err)
		
	}
	fmt.Println("ok")
	buffer := make([]byte, 10240)
	conn.Read(buffer)
	fmt.Printf("%q\n", string(buffer[:]))
	fmt.Println("===================================")
	var torrent map[interface{}]
	err = bencode.DecodeBytes(buffer,&torrent)
	if err != nil {
                fmt.Println("failed:", err)

        }
	fmt.Println(torrent)
	*/
}
