package arrmap

import (
	"sync"
	"sync/atomic"
)

// ArrayMap like a map, acctually an array, for append-only cases
type ArrayMap[K comparable, V any] struct {
	array     []*tuple[K, V]
	ls        []sync.Mutex
	size      atomic.Int32
	idxHasher func(K) uint32
}

type tuple[K comparable, V any] struct {
	key  K
	val  V
	next *tuple[K, V]
}

const (
	// expand practical capacity of array, to reduce the probability of hash collision
	// when factor = 1, more than 25% array elements become a linked list with at-least 2 nodes
	// when factor = 5, it's about 8%； when factor = 5, the number is less than 4%
	capacityFactor = 5
)

// NewArrayMap with an initialed capacity, decides the MAX elements this array-map can deal with
func NewArrayMap[T any](capacity int) ArrayMap[string, T] {
	if capacity <= 0 {
		panic("invalid capacity")
	}
	return ArrayMap[string, T]{
		array:     make([]*tuple[string, T], capacity*capacityFactor),
		ls:        make([]sync.Mutex, capacity*capacityFactor),
		idxHasher: fnv32,
	}
}

// NewArrayMapWithHasher custom key & hasher
func NewArrayMapWithHasher[K comparable, V any](capacity int, hasher func(K) uint32) *ArrayMap[K, V] {
	if capacity <= 0 {
		panic("invalid capacity")
	}
	return &ArrayMap[K, V]{
		array:     make([]*tuple[K, V], capacity*capacityFactor),
		ls:        make([]sync.Mutex, capacity*capacityFactor),
		idxHasher: hasher,
	}
}

func (m *ArrayMap[K, V]) getIdx(key K) int {
	if m.idxHasher != nil {
		return int(m.idxHasher(key)) % len(m.array)
	}
	return -1
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (m *ArrayMap[K, V]) Get(key K) (val V, ok bool) {
	idx := m.getIdx(key)
	if idx < 0 || idx >= len(m.array) {
		return
	}

	t := m.array[idx]
	if t == nil {
		return
	}
	for t != nil {
		if t.key == key {
			val = t.val
			ok = true
			return
		}
		t = t.next
	}
	return
}

func (m *ArrayMap[K, V]) Set(key K, val V) {
	idx := m.getIdx(key)
	if idx < 0 || idx >= len(m.array) {
		return
	}

	t := m.array[idx]
	// update 大多数无冲突的场景
	for t != nil {
		if t.key == key {
			t.val = val
			return
		}
		t = t.next
	}

	// insert 拉链法处理冲突
	m.ls[idx].Lock()
	defer m.ls[idx].Unlock()
	m.size.Add(1)
	t = m.array[idx]
	newT := &tuple[K, V]{
		key:  key,
		val:  val,
		next: t,
	}
	m.array[idx] = newT
	return
}

func (m *ArrayMap[K, V]) MSet(data map[K]V) {
	for k, v := range data {
		m.Set(k, v)
	}
}

func (m *ArrayMap[K, V]) Count() int {
	return int(m.size.Load())
}

func (m *ArrayMap[K, V]) Has(key K) bool {
	_, ok := m.Get(key)
	return ok
}

func (m *ArrayMap[K, V]) IsEmpty() bool {
	return m.Count() == 0
}

func (m *ArrayMap[K, V]) Items() map[K]V {
	tmp := make(map[K]V)
	for _, item := range m.array {
		for item != nil {
			tmp[item.key] = item.val
			item = item.next
		}
	}
	return tmp
}
