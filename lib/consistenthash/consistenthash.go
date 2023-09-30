package consistenthash

import (
	"hash/crc32"
	"sort"
)

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

func (n *NodeMap) AddNode(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := int(n.hashFunc([]byte(key)))
		n.nodeHashs = append(n.nodeHashs, hash)
		n.nodeHashMap[hash] = key
	}

	sort.Ints(n.nodeHashs)
}

func (n *NodeMap) PickNode(key string) string {
	if n.IsEmpty() {
		return ""
	}

	hash := int(n.hashFunc([]byte(key)))
	index := sort.Search(len(n.nodeHashs), func(i int) bool {
		return n.nodeHashs[i] >= hash
	})

	return n.nodeHashMap[n.nodeHashs[index%len(n.nodeHashs)]]
}
