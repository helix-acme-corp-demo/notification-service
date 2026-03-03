package store

import (
	"sync"

	"github.com/helix-acme-corp-demo/notification-service/internal/domain"
)

// Store is an in-memory notification store protected by a read-write mutex.
type Store struct {
	mu      sync.RWMutex
	entries map[string]*domain.Notification
}

// New creates an empty Store.
func New() *Store {
	return &Store{
		entries: make(map[string]*domain.Notification),
	}
}

// Save persists a notification in the store.
func (s *Store) Save(n *domain.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[n.ID] = n
}

// Find retrieves a notification by ID. The second return value indicates
// whether the notification was found.
func (s *Store) Find(id string) (*domain.Notification, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.entries[id]
	return n, ok
}

// All returns every notification in the store.
func (s *Store) All() []*domain.Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*domain.Notification, 0, len(s.entries))
	for _, n := range s.entries {
		result = append(result, n)
	}
	return result
}
