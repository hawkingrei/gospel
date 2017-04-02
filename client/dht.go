package dht

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"net"
	"time"
)

// DHT server
type DHT struct {
	conn     *net.UDPConn
	route    *Table
	secret   *secret
	searches *searches
	storages *storages
	tsecret  time.Time
}

func NewDHT(id *ID, conn *net.UDPConn, ksize int) *DHT {
	return &DHT{
		conn:     conn,
		route:    NewTable(id, ksize),
		secret:   newSecret(),
		searches: newSearches(),
		storages: newStorages(),
		tsecret:  time.Now(),
	}
}

// ID returns dht id
func (d *DHT) ID() *ID {
	return d.route.id
}

// Conn returns dht connection
func (d *DHT) Conn() *net.UDPConn {
	return d.conn
}

// Addr returns dht address
func (d *DHT) Addr() *net.UDPAddr {
	if conn := d.conn; conn != nil {
		return conn.LocalAddr().(*net.UDPAddr)
	}
	return nil
}

// Route returns route table
func (d *DHT) Route() *Table {
	return d.route
}

func (d *DHT) cleanNodes(tm time.Duration) {
	d.route.Map(func(b *Bucket) bool {
		if time.Since(b.time) > tm {
			if n := b.Random(); n != nil {
				d.FindNode(n.ID())
			}
		} else {
			b.clean(func(n *Node) bool {
				if n.pinged > 0 {
					return true
				}
				if time.Since(n.time) > tm {
					d.ping(n.addr)
					n.pinged++
				}
				return false
			})
		}
		return true
	})
}
