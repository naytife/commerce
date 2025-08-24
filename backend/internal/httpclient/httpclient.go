package httpclient

import (
	"context"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// DefaultClient is a shared, tuned HTTP client for internal service communication.
var DefaultClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	},
	// Rely on request contexts for per-request timeouts; leave Timeout=0
}

// DoWithRetry performs the request with simple exponential backoff retries for idempotent methods (GET/HEAD).
// ctx is used for cancellation. maxAttempts should be >=1.
func DoWithRetry(ctx context.Context, req *http.Request, maxAttempts int) (*http.Response, error) {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	method := req.Method
	attempt := 0
	backoff := 100 * time.Millisecond

	for {
		attempt++
		// Clone request with the provided context so body and headers are preserved
		reqWithCtx := req.Clone(ctx)
		resp, err := DefaultClient.Do(reqWithCtx)
		if err == nil {
			if resp.StatusCode < 500 {
				return resp, nil
			}
			// treat 5xx as retryable
			// close body before retrying
			_ = resp.Body.Close()
			log.Printf("httpclient: attempt %d got status %d for %s %s", attempt, resp.StatusCode, method, req.URL)
		} else {
			log.Printf("httpclient: attempt %d error for %s %s: %v", attempt, method, req.URL, err)
		}

		if attempt >= maxAttempts || !(method == http.MethodGet || method == http.MethodHead) {
			// If we have a resp (non-nil) return last resp error case; else return err
			if err != nil {
				return nil, err
			}
			return resp, nil
		}

		// Sleep with jitter
		jitter := time.Duration(rand.Int63n(int64(backoff)))
		select {
		case <-time.After(backoff + jitter):
			backoff *= 2
			if backoff > 5*time.Second {
				backoff = 5 * time.Second
			}
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
