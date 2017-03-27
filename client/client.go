package main

import (
	"fmt"
	"crypto/sha1"
	"os"
	"net"
	"github.com/zeebo/bencode"

)

func main(){
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
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	fmt.Println(string(buffer))
	fmt.Println("===================================")
	var torrent interface{}
	err = bencode.DecodeBytes(buffer,&torrent)
	if err != nil {
                fmt.Println("failed:", err)

        }
	fmt.Println(torrent)
}
