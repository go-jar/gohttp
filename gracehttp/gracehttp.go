package gracehttp

import (
	"net/http"
	"time"
)

const (
	DefaultReadTimeout  = 60 * time.Second
	DefaultWriteTimeout = DefaultReadTimeout
)

func ListenAndServe(addr string, handler http.Handler) error {
	return NewServer(addr, handler, DefaultReadTimeout, DefaultWriteTimeout).ListenAndServe()
}

func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	return NewServer(addr, handler, DefaultReadTimeout, DefaultWriteTimeout).ListenAndServeTLS(certFile, keyFile)
}
