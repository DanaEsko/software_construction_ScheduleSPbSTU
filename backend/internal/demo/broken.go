// Package demo shows the corrected versions of the code flagged by golangci-lint.
package demo

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// FetchStatus returns the HTTP status code of url.
func FetchStatus(ctx context.Context, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode, nil
}

// WriteReport writes a report to w.
func WriteReport(w io.Writer) error {
	if _, err := io.WriteString(w, "report"); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	return nil
}

// IsEOF reports whether err is io.EOF, including wrapped errors.
func IsEOF(err error) bool {
	return errors.Is(err, io.EOF)
}

// HashPassword returns the SHA-256 digest of p.
func HashPassword(p string) [32]byte {
	return sha256.Sum256([]byte(p))
}
