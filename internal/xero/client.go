package xero

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	apiBaseURL = "https://api.xero.com/api.xro/2.0/"

	// Xero allows 60 calls per minute per organisation. Stay just under it.
	minRequestInterval = 1100 * time.Millisecond
	maxRetries         = 4
)

type Client struct {
	tokens   *TokenSource
	tenantID string

	mu          sync.Mutex
	lastRequest time.Time
}

func NewClientFromEnv(ctx context.Context) (*Client, error) {
	cfg, err := ConfigFromEnv()
	if err != nil {
		return nil, err
	}
	ts, err := NewTokenSource(cfg)
	if err != nil {
		return nil, err
	}

	c := &Client{tokens: ts, tenantID: ts.TenantID()}

	// Custom connections are bound to one organisation and need no tenant header.
	// Web apps need one, so pick it automatically when there is no ambiguity.
	if c.tenantID == "" && ts.token.RefreshToken != "" {
		conns, err := ListConnections(ctx, ts)
		if err != nil {
			return nil, err
		}
		if len(conns) != 1 {
			return nil, fmt.Errorf("%d xero organisations connected - set XERO_TENANT_ID (see `car xero connections`)", len(conns))
		}
		c.tenantID = conns[0].TenantID
	}

	return c, nil
}

func (c *Client) Tokens() *TokenSource {
	return c.tokens
}

func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodGet, path, query, nil, out)
}

func (c *Client) put(ctx context.Context, path string, query url.Values, body, out any) error {
	return c.do(ctx, http.MethodPut, path, query, body, out)
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	var payload []byte
	if body != nil {
		var err error
		if payload, err = json.Marshal(body); err != nil {
			return err
		}
	}

	u := apiBaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	for attempt := 0; ; attempt++ {
		if err := c.throttle(ctx); err != nil {
			return err
		}

		token, err := c.tokens.AccessToken(ctx)
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, method, u, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if c.tenantID != "" {
			req.Header.Set("xero-tenant-id", c.tenantID)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}

		retryable := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
		if retryable && attempt < maxRetries {
			wait := time.Duration(1<<attempt) * 2 * time.Second
			if s, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil {
				wait = time.Duration(s) * time.Second
			}
			select {
			case <-time.After(wait):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// Xero returns 400 with per-item validation errors when summarizeErrors=false;
		// the body is still the normal envelope, so let the caller inspect it.
		perItemErrors := resp.StatusCode == http.StatusBadRequest && out != nil && bytes.Contains(respBody, []byte(`"StatusAttributeString"`))
		if resp.StatusCode >= 300 && !perItemErrors {
			return fmt.Errorf("xero %s %s failed: %d %s", method, path, resp.StatusCode, respBody)
		}

		if out == nil {
			return nil
		}
		return json.Unmarshal(respBody, out)
	}
}

func (c *Client) throttle(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	wait := time.Until(c.lastRequest.Add(minRequestInterval))
	if wait > 0 {
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	c.lastRequest = time.Now()
	return nil
}
