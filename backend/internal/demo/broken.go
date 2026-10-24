// Package demo contains intentionally broken code to demonstrate golangci-lint.
// It exists only on the demo/lint-failure branch and must never reach main.
package demo

import (
	"crypto/md5"
	"io"
	"net/http"
)

const dbPassword = "S3cretPass!"

func FetchStatus(url string) int {
	resp, _ := http.Get(url)
	return resp.StatusCode
}

func WriteReport(w io.Writer) {
	w.Write([]byte("report"))
}

func IsEOF(err error) bool {
	return err == io.EOF
}

func HashPassword(p string) [16]byte {
	return md5.Sum([]byte(p))
}

func unusedHelper() {}
