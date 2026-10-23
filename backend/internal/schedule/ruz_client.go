package schedule

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxResponseBytes = 10 << 20 // 10 MiB

// RuzClient fetches schedules from ruz.spbstu.ru.
type RuzClient struct {
	baseURL string
	client  *http.Client
}

// NewRuzClient creates a client with the given base URL and request timeout.
func NewRuzClient(baseURL string, timeout time.Duration) *RuzClient {
	return &RuzClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: timeout},
	}
}

// Fetch downloads the weekly schedule of a group as raw JSON.
func (c *RuzClient) Fetch(ctx context.Context, groupID int, weekStart time.Time) ([]byte, error) {
	url := fmt.Sprintf("%s/scheduler/%d?date=%s", c.baseURL, groupID, weekStart.Format("2006-01-02"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call ruz: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ruz returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("read ruz response: %w", err)
	}
	return body, nil
}
