package session

import (
	"sync"
)

// SESEIDSNEntity is  Session Establishment SEID against SN.
type SESEIDSNEntity struct {
	M   map[uint64]uint32
	mux sync.Mutex
}

// Inc add value for the given key.
func (s *SESEIDSNEntity) Inc(key uint64, value uint32) {
	s.mux.Lock()
	// Lock so only one goroutine at a time can access the map s.m.
	s.M[key] = value
	s.mux.Unlock()
}

// Value returns the current value  for the given key.
func (s *SESEIDSNEntity) Value(key uint64) uint32 {
	s.mux.Lock()
	// Lock so only one goroutine at a time can access the map s.M.
	defer s.mux.Unlock()
	return s.M[key]
}
