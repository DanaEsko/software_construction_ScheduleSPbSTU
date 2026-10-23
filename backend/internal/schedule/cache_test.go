package schedule

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeFetcher struct {
	payload []byte
	err     error
	calls   int
}

func (f *fakeFetcher) Fetch(_ context.Context, _ int, _ time.Time) ([]byte, error) {
	f.calls++
	return f.payload, f.err
}

func TestClassify(t *testing.T) {
	tests := []struct {
		name string
		age  time.Duration
		want State
	}{
		{"just fetched", 0, StateFresh},
		{"29 minutes", 29 * time.Minute, StateFresh},
		{"exactly 30 minutes", 30 * time.Minute, StateStale},
		{"2 days", 48 * time.Hour, StateStale},
		{"exactly 7 days", 7 * 24 * time.Hour, StateStale},
		{"8 days", 8 * 24 * time.Hour, StateTooOld},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Classify(tc.age); got != tc.want {
				t.Fatalf("Classify(%v) = %v, want %v", tc.age, got, tc.want)
			}
		})
	}
}

func TestServiceGet(t *testing.T) {
	now := time.Date(2026, 10, 23, 12, 0, 0, 0, time.UTC)
	week := WeekStart(now)
	errRuz := errors.New("ruz is down")

	tests := []struct {
		name      string
		cacheAge  time.Duration // negative means no cache entry
		fetcher   *fakeFetcher
		wantBody  string
		wantState State
		wantCalls int
		wantErr   error
		anyErr    bool
	}{
		{
			name:      "fresh cache is returned without calling ruz",
			cacheAge:  10 * time.Minute,
			fetcher:   &fakeFetcher{payload: []byte("new")},
			wantBody:  "cached",
			wantState: StateFresh,
			wantCalls: 0,
		},
		{
			name:      "stale cache is refreshed from ruz",
			cacheAge:  2 * time.Hour,
			fetcher:   &fakeFetcher{payload: []byte("new")},
			wantBody:  "new",
			wantState: StateFresh,
			wantCalls: 1,
		},
		{
			name:      "stale cache is served when ruz is down",
			cacheAge:  2 * time.Hour,
			fetcher:   &fakeFetcher{err: errRuz},
			wantBody:  "cached",
			wantState: StateStale,
			wantCalls: 1,
		},
		{
			name:      "too old cache is rejected when ruz is down",
			cacheAge:  8 * 24 * time.Hour,
			fetcher:   &fakeFetcher{err: errRuz},
			wantCalls: 1,
			wantErr:   ErrTooOld,
		},
		{
			name:      "no cache and ruz is down gives an error",
			cacheAge:  -1,
			fetcher:   &fakeFetcher{err: errRuz},
			wantCalls: 1,
			anyErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := NewMemoryStore()
			if tc.cacheAge >= 0 {
				err := store.Save(context.Background(), Entry{
					GroupID:   45608,
					WeekStart: week,
					Payload:   []byte("cached"),
					FetchedAt: now.Add(-tc.cacheAge),
				})
				if err != nil {
					t.Fatalf("prepare cache: %v", err)
				}
			}

			svc := NewService(store, tc.fetcher)
			svc.now = func() time.Time { return now }

			res, err := svc.Get(context.Background(), 45608, week)

			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, want %v", err, tc.wantErr)
			}
			if tc.anyErr && err == nil {
				t.Fatal("expected an error, got nil")
			}
			if tc.wantErr == nil && !tc.anyErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if string(res.Payload) != tc.wantBody {
					t.Fatalf("payload = %q, want %q", res.Payload, tc.wantBody)
				}
				if res.State != tc.wantState {
					t.Fatalf("state = %v, want %v", res.State, tc.wantState)
				}
			}
			if tc.fetcher.calls != tc.wantCalls {
				t.Fatalf("ruz calls = %d, want %d", tc.fetcher.calls, tc.wantCalls)
			}
		})
	}
}
