package schedule

import (
	"context"
	"sync"
	"time"
)

type cacheKey struct {
	groupID   int
	weekStart string
}

// MemoryStore is an in-memory Store. In the full project it is replaced by PostgreSQL.
type MemoryStore struct {
	mu   sync.RWMutex
	data map[cacheKey]Entry
}

// NewMemoryStore creates an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[cacheKey]Entry)}
}

func makeKey(groupID int, weekStart time.Time) cacheKey {
	return cacheKey{groupID: groupID, weekStart: weekStart.Format("2006-01-02")}
}

// Get returns the cached entry or ErrNotFound.
func (m *MemoryStore) Get(_ context.Context, groupID int, weekStart time.Time) (*Entry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.data[makeKey(groupID, weekStart)]
	if !ok {
		return nil, ErrNotFound
	}
	return &entry, nil
}

// Save stores or replaces an entry.
func (m *MemoryStore) Save(_ context.Context, entry Entry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[makeKey(entry.GroupID, entry.WeekStart)] = entry
	return nil
}
