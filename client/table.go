package dht

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"sort"
)

// Table store all nodes
type Table struct {
	id      *ID
	ksize   int
	buckets *list.List
}


// NewTable returns a table
func NewTable(id *ID, ksize int) *Table {
	t := &Table{
		id:      id,
		ksize:   ksize,
		buckets: list.New(),
	}
	b := NewBucket(ZeroID, ksize)
	t.buckets.PushBack(b)
	return t
}

func (t *Table) KSize() int {
	return t.ksize
}

// NumNodes returns all node count
func (t *Table) NumNodes() (n int) {
	t.Map(func(b *Bucket) bool {
		n += b.Count()
		return true
	})
	return
}

// Insert a node
func (t *Table) Insert(id *ID, addr *net.UDPAddr) (*Node, error) {
	if id.Compare(t.id) == 0 {
		return nil, errors.New("id equal to table's id")
	}
	return t.insert(id, addr)
}

func (t *Table) insert(id *ID, addr *net.UDPAddr) (n *Node, err error) {
	if e := t.find(id); e != nil {
		if n = e.Value.(*Bucket).Insert(id, addr); n != nil {
			return
		}
		if inBucket(t.id, e) && t.split(e) {
			return t.insert(id, addr)
		}
	}
	err = errors.New("drop this node")
	return
}


func (t *Table) split(e *list.Element) bool {
	bit := e.Value.(*Bucket).first.LowBit()
	if next := e.Next(); next != nil {
		bit2 := next.Value.(*Bucket).first.LowBit()
		if bit < bit2 {
			bit = bit2
		}
	}
	if bit++; bit >= 160 {
		return false
	}

	b := e.Value.(*Bucket)
	first, _ := NewID(b.first.Bytes())
	first.SetBit(bit, true)
	b2 := NewBucket(first, b.cap)
	t.buckets.InsertAfter(b2, e)

	var eles []*list.Element
	b.handle(func(be *list.Element) bool {
		if inBucket(be.Value.(*Node).id, e) == false {
			eles = append(eles, be)
		}
		return true
	})
	for _, ele := range eles {
		b2.nodes.PushBack(b.nodes.Remove(ele))
	}

	return true
}

// Find returns bucket
func (t *Table) Find(id *ID) *Bucket {
	if e := t.find(id); e != nil {
		return e.Value.(*Bucket)
	}
	return nil
}

func (t *Table) find(id *ID) (ele *list.Element) {
	t.handle(func(e *list.Element) bool {
		if inBucket(id, e) {
			ele = e
			return false
		}
		return true
	})
	return
}