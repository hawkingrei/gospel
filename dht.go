package main

import (
	"encoding/hex"
	"errors"
	"math"
	"net"
	"time"
)

// Config represents the configure of dht.
type Config struct {
	// in mainline dht, k = 8
	K int
	// for crawling mode, we put all nodes in one bucket, so KBucketSize may
	// not be K
	KBucketSize int
	// candidates are udp, udp4, udp6
	Network string
	// format is `ip:port`
	Address string
	// the prime nodes through which we can join in dht network
	PrimeNodes []string
	// the kbucket expired duration
	KBucketExpiredAfter time.Duration
	// the node expired duration
	NodeExpriedAfter time.Duration
	// how long it checks whether the bucket is expired
	CheckKBucketPeriod time.Duration
	// peer token expired duration
	TokenExpiredAfter time.Duration
	// the max transaction id
	MaxTransactionCursor uint64
	// how many nodes routing table can hold
	MaxNodes int
	// callback when got get_peers request
	OnGetPeers func(string, string, int)
	// callback when got announce_peer request
	OnAnnouncePeer func(string, string, int)
	// blcoked ips
	BlockedIPs []string
	// blacklist size
	BlackListMaxSize int
	// the times it tries when send fails
	Try int
	// the size of packet need to be dealt with
	PacketJobLimit int
	// the size of packet handler
	PacketWorkerLimit int
	// the nodes num to be fresh in a kbucket
	RefreshNodeNum int
}

// NewStandardConfig returns a Config pointer with default values.
func NewStandardConfig() *Config {
	return &Config{
		K:           8,
		KBucketSize: 8,
		Network:     "udp4",
		Address:     ":6881",
		PrimeNodes: []string{
			"router.bittorrent.com:6881",
			"router.utorrent.com:6881",
			"dht.transmissionbt.com:6881",
		},
		NodeExpriedAfter:     time.Duration(time.Minute * 15),
		KBucketExpiredAfter:  time.Duration(time.Minute * 15),
		CheckKBucketPeriod:   time.Duration(time.Second * 30),
		TokenExpiredAfter:    time.Duration(time.Minute * 10),
		MaxTransactionCursor: math.MaxUint32,
		MaxNodes:             5000,
		BlockedIPs:           make([]string, 0),
		BlackListMaxSize:     65536,
		Try:                  2,
		PacketJobLimit:       1024,
		PacketWorkerLimit:    256,
		RefreshNodeNum:       8,
	}
}

// DHT represents a DHT node.
type DHT struct {
	*Config
	node               *node
	conn               *net.UDPConn
	routingTable       *routingTable
	transactionManager *transactionManager
	peersManager       *peersManager
	tokenManager       *tokenManager
	blackList          *blackList
	Ready              bool
	packets            chan packet
	workerTokens       chan struct{}
}

// New returns a DHT pointer. If config is nil, then config will be set to
// the default config.
func New(config *Config) *DHT {
	if config == nil {
		config = NewStandardConfig()
	}

	node, err := newNode(randomString(20), config.Network, config.Address)
	if err != nil {
		panic(err)
	}

	d := &DHT{
		Config:       config,
		node:         node,
		blackList:    newBlackList(config.BlackListMaxSize),
		packets:      make(chan packet, config.PacketJobLimit),
		workerTokens: make(chan struct{}, config.PacketWorkerLimit),
	}

	for _, ip := range config.BlockedIPs {
		d.blackList.insert(ip, -1)
	}

	go func() {
		for _, ip := range getLocalIPs() {
			d.blackList.insert(ip, -1)
		}

		ip, err := getRemoteIP()
		if err != nil {
			d.blackList.insert(ip, -1)
		}
	}()

	return d

}

// init initializes global varables.
func (dht *DHT) init() {
	listener, err := net.ListenPacket(dht.Network, dht.Address)
	if err != nil {
		panic(err)
	}

	dht.conn = listener.(*net.UDPConn)
	dht.routingTable = newRoutingTable(dht.KBucketSize, dht)
	dht.peersManager = newPeersManager(dht)
	dht.tokenManager = newTokenManager(dht.TokenExpiredAfter, dht)
	dht.transactionManager = newTransactionManager(
		dht.MaxTransactionCursor, dht)

	go dht.transactionManager.run()
	go dht.tokenManager.clear()
	go dht.blackList.clear()
}

// join makes current node join the dht network.
func (dht *DHT) join() {
	for _, addr := range dht.PrimeNodes {
		raddr, err := net.ResolveUDPAddr(dht.Network, addr)
		if err != nil {
			continue
		}

		// NOTE: Temporary node has NOT node id.
		dht.transactionManager.findNode(
			&node{addr: raddr},
			dht.node.id.RawString(),
		)
	}
}

// listen receives message from udp.
func (dht *DHT) listen() {
	go func() {
		buff := make([]byte, 8192)
		for {
			n, raddr, err := dht.conn.ReadFromUDP(buff)
			if err != nil {
				continue
			}

			dht.packets <- packet{buff[:n], raddr}
		}
	}()
}

// id returns a id near to target if target is not null, otherwise it returns
// the dht's node id.
func (dht *DHT) id(target string) string {
	if target == "" {
		return dht.node.id.RawString()
	}
	return target[:15] + dht.node.id.RawString()[15:]
}

// GetPeers returns peers who have announced having infoHash.
func (dht *DHT) GetPeers(infoHash string) ([]*Peer, error) {
	if !dht.Ready {
		return nil, errors.New("dht not ready")
	}

	if len(infoHash) == 40 {
		data, err := hex.DecodeString(infoHash)
		if err != nil {
			return nil, err
		}
		infoHash = string(data)
	}

	peers := dht.peersManager.GetPeers(infoHash, dht.K)
	if len(peers) != 0 {
		return peers, nil
	}

	ch := make(chan struct{})

	go func() {
		neighbors := dht.routingTable.GetNeighbors(
			newBitmapFromString(infoHash), dht.K)

		for _, no := range neighbors {
			dht.transactionManager.getPeers(no, infoHash)
		}

		i := 0
		for _ = range time.Tick(time.Second * 1) {
			i++
			peers = dht.peersManager.GetPeers(infoHash, dht.K)
			if len(peers) != 0 || i == 30 {
				break
			}
		}

		ch <- struct{}{}
	}()

	<-ch
	return peers, nil
}

// Run starts the dht.
func (dht *DHT) Run() {
	dht.init()
	dht.listen()
	dht.join()

	dht.Ready = true

	var pkt packet
	tick := time.Tick(dht.CheckKBucketPeriod)

	for {
		select {
		case pkt = <-dht.packets:
			handle(dht, pkt)
		case <-tick:
			if dht.routingTable.Len() == 0 {
				dht.join()
			} else if dht.transactionManager.len() == 0 {
				go dht.routingTable.Fresh()
			}
		}
	}
}