package session

import (
	"sync"
)

// SessionEntity is safe to use concurrently.
type SEIDEntity struct {
	M   map[uint64]uint64
	mux sync.Mutex
}

// Inc add value for the given key.
func (s *SEIDEntity) Inc(key uint64, value uint64) {
	s.mux.Lock()
	s.M[key] = value
	s.mux.Unlock()
}

// Value returns the current value  for the given key.
func (s *SEIDEntity) Value(key uint64) uint64 {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.M[key]
}
