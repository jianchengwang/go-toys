package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	mu       sync.Mutex
	hash     Hash
	replicas int
	values   atomic.Value
}

type values struct {
	keys    []int
	hashMap map[int]string
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
	}
	m.values.Store(&values{
		hashMap: make(map[int]string),
	})
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
func (m *Map) Add(keys ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	newValues := m.loadValues()
	for _, key := range keys {
		// 对每个 key(节点) 创建 m.replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			newValues.keys = append(newValues.keys, hash)
			newValues.hashMap[hash] = key
		}
	}
	sort.Ints(newValues.keys)
	m.values.Store(newValues)
}

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	values := m.loadValues()
	if len(values.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(values.keys), func(i int) bool {
		return values.keys[i] >= hash
	})
	// 如果 idx == len(m.keys)，说明应选择 m.keys[0]，
	// 因为 m.keys 是一个环状结构，用取余数的方式来处理这种情况
	return values.hashMap[values.keys[idx%len(values.keys)]]
}

func (m *Map) Remove(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	newValues := m.loadValues()

	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(newValues.keys, hash)
		if newValues.keys[idx] != hash {
			return
		}
		newValues.keys = append(newValues.keys[:idx], newValues.keys[idx+1:]...)
		delete(newValues.hashMap, hash)
	}

	m.values.Store(newValues)
}

func (m *Map) loadValues() *values {
	return m.values.Load().(*values)
}

func (m *Map) copyValues() *values {
	oldValues := m.loadValues()
	newValues := &values{
		keys:    make([]int, len(oldValues.keys)),
		hashMap: make(map[int]string),
	}
	copy(newValues.keys, oldValues.keys)
	for k, v := range oldValues.hashMap {
		newValues.hashMap[k] = v
	}
	return newValues
}
