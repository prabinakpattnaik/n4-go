package session

import (
	"sync"

	"github.com/prabinakpattnaik/n4-go/msg"
)

type SessionRequestResponse struct {
	SRequest  msg.PFCP
	SResponse msg.PFCP
}

// SessionEntity is safe to use concurrently.
type SessionEntity struct {
	M   map[uint32]SessionRequestResponse
	mux sync.Mutex
}

// Inc add value for the given key.
func (s *SessionEntity) Inc(key uint32, srr SessionRequestResponse) {
	s.mux.Lock()
	// Lock so only one goroutine at a time can access the map s.m.
	srrValueExist, exist := s.M[key]
	if exist {
		if srrValueExist.SRequest != nil && srrValueExist.SResponse == nil {
			srrValueExist.SResponse = srr.SResponse
			s.M[key] = srrValueExist
		}

	} else {
		s.M[key] = srr
	}
	s.mux.Unlock()
}

// Value returns the current value of the SessionEntity for the given key.
func (s *SessionEntity) Value(key uint32) SessionRequestResponse {
	s.mux.Lock()
	// Lock so only one goroutine at a time can access the map s.M.
	defer s.mux.Unlock()
	return s.M[key]
}
