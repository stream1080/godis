package consistenthash

import "hash/crc32"

type HashFunc func(data []byte) uint32

type NodeMap struct {
	hashFunc    HashFunc
	nodeHashs   []int
	nodeHashMap map[int]string
}

func NewNodeMap(hashFunc HashFunc) *NodeMap {
	if hashFunc == nil {
		hashFunc = crc32.ChecksumIEEE
	}

	return &NodeMap{
		hashFunc:    hashFunc,
		nodeHashMap: make(map[int]string),
	}
}

func (n *NodeMap) IsEmpty() bool {
	return len(n.nodeHashs) == 0
}
