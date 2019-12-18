package session

import (
	"sync"
)

//TODO F-SEID (CP) <-> F-SEID(DP)

// SessionEntity is safe to use concurrently.
type CPSEIDDPSEIDEntity struct {
	M   map[uint64]uint64
	mux sync.Mutex
}

// Inc add value for the given key.
func (c *CPSEIDDPSEIDEntity) Inc(key uint64, value uint64) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map s.m.
	c.M[key] = value
	c.mux.Unlock()
}

// Value returns the current value  for the given key.
func (c *CPSEIDDPSEIDEntity) Value(key uint64) uint64 {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map s.M.
	defer c.mux.Unlock()
	return c.M[key]
}
