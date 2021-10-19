package consistenthash

import (
	"hash/crc32"
	"sort"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed nodesHash
type Consistenthash struct {
	hash     Hash
	replicas int
	nodesHash     []int // hash ring,Sorted
	hashToRNodeName  map[int]string //virtual hash node value -> real node name
}

// New creates a Map instance
func New(fn Hash) *Consistenthash {
	m := &Consistenthash{
		replicas: 1,
		hash:     fn,
		hashToRNodeName:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some nodesHash to the hash.
// allow to input one or more real node name
// real node name
func (m *Consistenthash) Add(nodes ...string) {
	for _, name := range nodes {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(name)))
			m.nodesHash = append(m.nodesHash, hash)
			m.hashToRNodeName[hash] = name
		}
	}
	sort.Ints(m.nodesHash)
}

// Get gets the closest item in the hash to the provided key.
func (m *Consistenthash) Get(key string) string {
	if len(m.nodesHash) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.nodesHash), func(i int) bool {
		return m.nodesHash[i] >= hash
	})

	// return the real node name, which stores the key's value
	return m.hashToRNodeName[m.nodesHash[idx%len(m.nodesHash)]]
}