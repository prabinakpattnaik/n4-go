package session

import (
        "sync"
)

// SessionInfo is session information
type SessionInfo struct {
	Sid string
	Rule_id string
	Ue_ipv4 string
	Teid uint32
	Session_state uint8
}

// SEIDSessionInfo store session info with key
type SEIDSessionInfo struct {
	M   map[uint64]SessionInfo
        mux sync.Mutex
}

// Inc add value for the given key.
func (s *SEIDSessionInfo) Inc(key uint64, si SessionInfo) {
        s.mux.Lock()
        // Lock so only one goroutine at a time can access the map s.m.
        s.M[key] = si
        s.mux.Unlock()
}

// Value returns the current value  for the given key.
func (s *SEIDSessionInfo) Value(key uint64) SessionInfo {
        s.mux.Lock()
        // Lock so only one goroutine at a time can access the map s.M.
        defer s.mux.Unlock()
        return s.M[key]
}

// Delete the entry for the given key.
func (s *SEIDSessionInfo) Delete(key uint64) {
        s.mux.Lock()
        // Lock so only one goroutine at a time can access the map s.M.
        defer s.mux.Unlock()
        delete(s.M, key)
}

