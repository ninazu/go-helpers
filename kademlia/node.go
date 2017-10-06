package kademlia

import (
	"crypto/rand"
	"net"
)

type Node struct {
	id   []byte
	ip   net.IP
	port int
}

type Options struct {
	self  *Node
	nodes []*Node
}

func NewNode(options *Options) (*Node, error) {
	node := &Node{}
	node.id, _ = generateId()

	return node, nil
}

func generateId() ([]byte, error) {
	result := make([]byte, 20)
	_, err := rand.Read(result)

	return result, err
}
