package core

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

// Specialized forms of the concurrent map
type BindingConcurrentMap struct {
	ConcurrentMap
}

type BindingTuple struct {
	Key string
	Val string
}

// Map of type [uint64: *boundEndpint]
func NewConcurrentBindingMap() BindingConcurrentMap {
	return BindingConcurrentMap{NewConcurrentMap(5)}
}

// func (m BindingConcurrentMap) Iter() <-chan BindingTuple {
// 	ch := make(chan BindingTuple, 10)
// 	go func() {
// 		for _, shard := range m.ConcurrentMap {
// 			shard.RLock()

// 			for key, val := range shard.items {
// 				ch <- BindingTuple{key, val.(string)}
// 			}

// 			shard.RUnlock()
// 		}

// 		close(ch)
// 	}()
// 	return ch
// }

// Sets the given value under the specified key.
func (m *BindingConcurrentMap) Set(key string, value string) {
	m.ConcurrentMap.Set(KeyFromString(key), value)
}

// Retrieves an element from map under given key.
func (m BindingConcurrentMap) Get(key string) (string, bool) {
	if val, ok := m.ConcurrentMap.Get(KeyFromString(key)); ok {
		return val.(string), true
	} else {
		return "", false
	}
}

func TestGet(t *testing.T) {
	b := NewConcurrentBindingMap()
	b.Set("greeting", "hello")
	c, ok := b.Get("greeting")

	True(t, ok)
	Equal(t, "hello", c)
}

func TestBadGet(t *testing.T) {
	b := NewConcurrentBindingMap()
	b.Set("greeting", "hello")
	_, ok := b.Get("notagreeting")

	False(t, ok)
}

func TestIterator(t *testing.T) {
	b := NewConcurrentBindingMap()
	b.Set("greeting", "hello")
	_, ok := b.Get("notagreeting")

	False(t, ok)
}
