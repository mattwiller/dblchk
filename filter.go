package dblchk

import (
	"hash"
	"hash/fnv"
	"sync"
	"sync/atomic"
)

const DEFAULT_CAPACITY uint = 1 << 16 /* 256 KiB / sizeof(uint32) */

var hashPool *sync.Pool = &sync.Pool{
	New: func() any { return fnv.New32() },
}

func hashElement(element []byte) hashKey {
	h := hashPool.Get().(hash.Hash32)
	defer hashPool.Put(h)

	h.Reset()
	h.Write(element)
	return hashKey(h.Sum32())
}

// Key into the filter, with a composite internal structure.
//
// The 32-bit (4-byte) layout of the key is as follows:
// Upper 16 bits (two bytes): idx
// Low 5 bits of middle byte: posA
// Low 5 bits of lower byte: posB
type hashKey uint32

func (key hashKey) index() uint16 {
	return uint16(key >> 16)
}

func (key hashKey) mask() uint32 {
	posA := byte((key & 0x0000_1f00) >> 8)
	posB := byte(key & 0x0000_001f)
	return uint32((1 << posA) | (1 << posB))
}

type Filter []uint32

// Create a new filter with the given number of 32-bit blocks.
//
// If capacity is set to 0, the default capacity of 64k (=256 KiB) is used.
func NewFilter(capacity uint) Filter {
	if capacity == 0 {
		capacity = DEFAULT_CAPACITY
	}
	return make(Filter, capacity)
}

// Add the given element to the filter.
func (filter Filter) Add(element []byte) {
	hash := hashElement(element)
	atomic.OrUint32(&filter[hash.index()], hash.mask())
}

// Check the filter for an element.
//
// If the method returns `true`, the element may be present in the filter.
// Otherwise, it is guaranteed not to be in the filter.
func (filter Filter) MayContain(element []byte) bool {
	hash := hashElement(element)
	mask := hash.mask()

	block := filter[hash.index()]
	return (block & mask) == mask
}

// Clear the filter, resetting it back to an empty state.
func (filter Filter) Reset() {
	clear(filter)
}
