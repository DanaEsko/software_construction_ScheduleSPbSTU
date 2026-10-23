// Package schedule implements fetching and caching of group schedules
// from the external ruz.spbstu.ru API with graceful degradation.
package schedule

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// State describes how fresh the returned schedule is.
type State string

const (
	// StateFresh means the data is younger than FreshTTL.
	StateFresh State = "fresh"
	// StateStale means the data is older than FreshTTL but still usable.
	StateStale State = "stale"
	// StateTooOld means the data is older than MaxStaleAge and must not be served.
	StateTooOld State = "too_old"
)

const (
	// FreshTTL is the age until which a cache entry is returned without calling ruz.
	FreshTTL = 30 * time.Minute
	// MaxStaleAge is the maximum age of a cache entry that may be served when ruz is down.
	MaxStaleAge = 7 * 24 * time.Hour
)

var (
	// ErrNotFound is returned by a Store when there is no entry for the key.
	ErrNotFound = errors.New("cache entry not found")
	// ErrTooOld is returned when ruz is unavailable and the cache is too old.
	ErrTooOld = errors.New("cached schedule is too old")
)

// Entry is a cached weekly schedule of one group.
type Entry struct {
	GroupID   int
	WeekStart time.Time
	Payload   []byte
	FetchedAt time.Time
}

// Store persists cache entries.
type Store interface {
	Get(ctx context.Context, groupID int, weekStart time.Time) (*Entry, error)
	Save(ctx context.Context, entry Entry) error
}

// Fetcher loads a fresh schedule from the external API.
type Fetcher interface {
	Fetch(ctx context.Context, groupID int, weekStart time.Time) ([]byte, error)
}

// Result is the outcome of Service.Get.
type Result struct {
	Payload   []byte
	State     State
	FromCache bool
}

// Service implements the fault-tolerant schedule lookup algorithm.
type Service struct {
	store   Store
	fetcher Fetcher
	now     func() time.Time
}

// NewService creates a Service.
func NewService(store Store, fetcher Fetcher) *Service {
	return &Service{store: store, fetcher: fetcher, now: time.Now}
}

// Classify returns the cache state for an entry of the given age.
func Classify(age time.Duration) State {
	switch {
	case age < FreshTTL:
		return StateFresh
	case age <= MaxStaleAge:
		return StateStale
	default:
		return StateTooOld
	}
}

// WeekStart returns Monday 00:00 UTC of the week that contains t.
func WeekStart(t time.Time) time.Time {
	offset := (int(t.Weekday()) + 6) % 7
	d := t.AddDate(0, 0, -offset)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

// Get returns the schedule of a group for the week that starts at weekStart.
//
// Fresh cache is returned without calling ruz. Otherwise ruz is asked for new
// data; if it fails, the cached copy is returned (marked stale) unless it is
// older than MaxStaleAge.
func (s *Service) Get(ctx context.Context, groupID int, weekStart time.Time) (Result, error) {
	entry, err := s.store.Get(ctx, groupID, weekStart)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return Result{}, fmt.Errorf("read cache: %w", err)
	}
	found := err == nil

	if found && Classify(s.now().Sub(entry.FetchedAt)) == StateFresh {
		return Result{Payload: entry.Payload, State: StateFresh, FromCache: true}, nil
	}

	payload, fetchErr := s.fetcher.Fetch(ctx, groupID, weekStart)
	if fetchErr == nil {
		saveErr := s.store.Save(ctx, Entry{
			GroupID:   groupID,
			WeekStart: weekStart,
			Payload:   payload,
			FetchedAt: s.now(),
		})
		if saveErr != nil {
			return Result{}, fmt.Errorf("save cache: %w", saveErr)
		}
		return Result{Payload: payload, State: StateFresh}, nil
	}

	if !found {
		return Result{}, fmt.Errorf("fetch schedule: %w", fetchErr)
	}
	if Classify(s.now().Sub(entry.FetchedAt)) == StateTooOld {
		return Result{}, fmt.Errorf("%w: %w", ErrTooOld, fetchErr)
	}
	return Result{Payload: entry.Payload, State: StateStale, FromCache: true}, nil
}
