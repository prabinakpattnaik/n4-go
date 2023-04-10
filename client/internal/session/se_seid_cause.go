package session

import (
        "sync"
)

// SEIDCausentity is  Session Establishment SEID against cause value return from response.
type SEIDCauseEntity struct {
        M   map[uint64]uint8
        mux sync.Mutex
}
// Inc add value for the given key.
func (s *SEIDCauseEntity) Inc(key uint64, value uint8) {
        s.mux.Lock()
        // Lock so only one goroutine at a time can access the map s.m.
        s.M[key] = value
        s.mux.Unlock()
}

// Value returns the current value  for the given key.
func (s *SEIDCauseEntity) Value(key uint64) uint8 {
        s.mux.Lock()
        // Lock so only one goroutine at a time can access the map s.M.
        defer s.mux.Unlock()
        return s.M[key]
}

// Delete the entry for the given key.
func (s *SEIDCauseEntity) Delete(key uint64) {
        s.mux.Lock()
        // Lock so only one goroutine at a time can access the map s.M.
        defer s.mux.Unlock()
        delete(s.M, key)
}
