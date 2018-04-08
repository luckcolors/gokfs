package test

import (
	"math/rand"
	"sync"
	"testing"
)

type thing struct{}

var mutex = &sync.RWMutex{}
var m = make(map[uint64]*thing)
var rng = rand.New(rand.NewSource(0))

func GetIndex(index uint64) *thing {
	mutex.RLock()
	v, ok := m[index]
	mutex.RUnlock()
	if !ok {
		mutex.Lock()
		v = new(thing)
		m[index] = v
		mutex.Unlock()
	}
	return v
}

func BenchmarkGetIndex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		go GetIndex(uint64(rng.Intn(4)))
	}
}
