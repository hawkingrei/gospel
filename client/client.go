package client

import (
	"crypto"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/anacrolix/missinggo"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/iplist"
	"github.com/anacrolix/torrent/logonce"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/tylertreat/BoomFilters"

	"github.com/hawkingrei/gospel/dht/krpc"
)


type Server struct {
	id               string
	socket           net.PacketConn
	transactions     map[transactionKey]*Transaction
	transactionIDInt uint64
	nodes            map[string]*node // Keyed by dHTAddr.String().
	mu               sync.Mutex
	closed           missinggo.Event
	ipBlockList      iplist.Ranger
	badNodes         *boom.BloomFilter
	tokenServer      tokenServer

	numConfirmedAnnounces int
	bootstrapNodes        []string
	config                ServerConfig
}