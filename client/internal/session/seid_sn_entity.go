package session

import (
	"sync"
)

type SNCollection struct {
	SESN uint32
	SMSN uint32
	SDSN uint32
}

// SessionEntity is safe to use concurrently.
type SEIDSNEntity struct {
	M   map[uint64]SNCollection
	mux sync.Mutex
}

// Inc add value for the given key.
func (s *SEIDSNEntity) Inc(key uint64, whichType uint8, sn uint32) {
	s.mux.Lock()
	// Lock so only one goroutine at a time can access the map s.m.
	sncValueExist, exist := s.M[key]
	if exist {
		if whichType == 1 {
			sncValueExist.SESN = sn
		} else if whichType == 2 {
			sncValueExist.SMSN = sn
		} else if whichType == 3 {
			sncValueExist.SMSN = sn
		}

		s.M[key] = sncValueExist
	} else {
		var snc SNCollection
		if whichType == 1 {
			snc.SESN = sn
		} else if whichType == 2 {
			snc.SMSN = sn
		} else if whichType == 3 {
			snc.SMSN = sn
		}
		s.M[key] = snc
	}
	s.mux.Unlock()
}

// Value returns the current value of the SessionEntity for the given key.
func (s *SEIDSNEntity) Value(key uint64) SNCollection {
	s.mux.Lock()
	// Lock so only one goroutine at a time can access the map s.M.
	defer s.mux.Unlock()
	return s.M[key]
}
