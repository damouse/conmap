package core

import (
	"hash/fnv"
	"sync"
)

type ShardSet []*Shard

// A sharded concurrent map with internal RW locks
type ConcurrentMap struct {
	shards []*Shard
	count  int // The number of shards in the map
}

// A goroutine safe string to anything map.
type Shard struct {
	items map[uint64]interface{}
	sync.RWMutex
}

type Tuple struct {
	Key uint64
	Val interface{}
}

func NewConcurrentMap(numShards int) ConcurrentMap {
	m := ConcurrentMap{
		shards: make([]*Shard, numShards),
		count:  numShards,
	}

	for i := 0; i < numShards; i++ {
		m.shards[i] = &Shard{items: make(map[uint64]interface{})}
	}

	return m
}

func (m ConcurrentMap) GetShard(key uint64) *Shard {
	k := m.shards[key%uint64(m.count)]

	return k
}

func (m *ConcurrentMap) Set(key uint64, value interface{}) {
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.items[key] = value
}

func (m ConcurrentMap) Get(key uint64) (interface{}, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

func (m ConcurrentMap) Count() int {
	count := 0

	for i := 0; i < m.count; i++ {
		shard := m.shards[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}

	return count
}

func (m *ConcurrentMap) Has(key uint64) bool {
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	_, ok := shard.items[key]
	return ok
}

func (m *ConcurrentMap) RemoveKey(key uint64) {
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

func (m *ConcurrentMap) IsEmpty() bool {
	return m.Count() == 0
}

// Returns a buffered iterator which could be used in a for range loop.
func (m ConcurrentMap) Iter() <-chan Tuple {
	ch := make(chan Tuple, 10)

	go func() {
		for _, shard := range m.shards {
			shard.RLock()

			for key, val := range shard.items {
				ch <- Tuple{key, val}
			}

			shard.RUnlock()
		}
		close(ch)
	}()

	return ch
}

func (m ConcurrentMap) Items() map[uint64]interface{} {
	tmp := make(map[uint64]interface{})

	for item := range m.Iter() {
		tmp[item.Key] = item.Val
	}

	return tmp
}

// Remove an element by value and return it
// func (m *ConcurrentMap) RemoveValue(v interface{}) {
//  // Try to get shard.
//  shard := m.GetShard(key)
//  shard.Lock()
//  defer shard.Unlock()
//  delete(shard.items, key)
// }

// Maps that want to use strings as keys should call this function when resolving keys
func KeyFromString(k string) uint64 {
	h := fnv.New64()
	h.Write([]byte(k))
	return h.Sum64()
}
