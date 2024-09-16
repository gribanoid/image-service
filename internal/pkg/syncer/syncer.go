package syncer

import "sync"

type Syncer struct {
	dates map[string]*sync.Mutex
}

func New() *Syncer {
	return &Syncer{dates: make(map[string]*sync.Mutex, 1<<10)}
}

func (s *Syncer) AddDate(date string) {
	mu := sync.Mutex{}
	s.dates[date] = &mu
}

func (s *Syncer) Lock(date string) bool {
	mu, ok := s.dates[date]
	if !ok {
		return false
	}
	mu.Lock()

	return true
}

func (s *Syncer) Unlock(date string) bool {
	mu, ok := s.dates[date]
	if !ok {
		return false
	}

	mu.Unlock()

	return false
}
