package main

import (
	"fmt"
	"math/rand"
	"math"
	"net"
	"github.com/zeebo/bencode"
	/*
	"crypto/sha1"
	"os"
	"net"
	
	*/
)

func encodeMessage(msg interface{}) ([]byte, error) {
	return bencode.EncodeBytes(msg)
}

var (
	tidVals = map[string]string{
		"ping": "pn", "find_node": "fn", "get_peers": "gp", "announce_peer": "ap",
	}
	tidFuns = map[string]string{
		"pn": "ping", "fn": "find_node", "gp": "get_peers", "ap": "announce_peer",
	}
)


func encodeTID(q string, id int16) (tid []byte) {
	if val, ok := tidVals[q]; ok {
		uid := uint16(id)
		if id <= 0 {
			uid = math.MaxUint16
		}
		tid = make([]byte, 4)
		copy(tid[:2], val[:2])
		tid[2] = byte(uid & 0xFF00 >> 8)
		tid[3] = byte(uid & 0x00FF)
	}
	return
}

type kadQueryMessage struct {
	T []byte                 `bencode:"t"`
	Y string                 `bencode:"y"`
	Q string                 `bencode:"q"`
	A map[string]interface{} `bencode:"a"`
}

type kadReplyMessage struct {
	T []byte                 `bencode:"t"`
	Y string                 `bencode:"y"`
	R map[string]interface{} `bencode:"r"`
}

func newQueryMessage(tid []byte, q string, data map[string]interface{}) *kadQueryMessage {
	return &kadQueryMessage{tid, "q", q, data}
}

func FindNodeFromAddr(id [20]byte, addr *net.UDPAddr,conn *net.UDPConn) error {
	data := map[string]interface{}{
		"id":     id[:],
		"target": id[:],
	}
	return queryMessage("find_node", 0, addr, data,conn)
}

func sendMessage(addr *net.UDPAddr, data []byte,conn     *net.UDPConn) (err error) {
	for n, nn := 0, 0; nn < len(data); nn += n {
		n, err = conn.WriteToUDP(data[nn:], addr)
		if err != nil {
			break
		}
	}
	return
}

func queryMessage(q string, no int16, addr *net.UDPAddr, data map[string]interface{},conn     *net.UDPConn) (err error) {
	msg := newQueryMessage(encodeTID(q, no), q, data)
	if b, err := encodeMessage(msg); err == nil {
		err = sendMessage(addr, b,conn)
	}
	return
}



func main(){

	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return
	}
	id := new([20]byte);
	rand.Read(id[:])
	fmt.Println(id)
	var routers = []string{
		"router.magnets.im:6881",
		"router.bittorrent.com:6881",
		"dht.transmissionbt.com:6881",
		"router.utorrent.com:6881",
	}
	for _, addr := range routers {
		addr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			break
		}
		err = FindNodeFromAddr(*id, addr,conn.(*net.UDPConn))
		if err != nil {
			break
		}
	}
}
