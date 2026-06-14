// Package regex101 provides access to the regex101.com pattern library.
package regex101

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Config holds client configuration.
type Config struct {
	BaseURL   string
	Rate      time.Duration
	Timeout   time.Duration
	Retries   int
	UserAgent string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://regex101.com",
		Rate:      300 * time.Millisecond,
		Timeout:   30 * time.Second,
		Retries:   3,
		UserAgent: "regex101-cli/0.1 (https://github.com/tamnd/regex101-cli)",
	}
}

// Client fetches data from regex101.com.
type Client struct {
	cfg  Config
	http *http.Client
	last time.Time
}

// NewClient creates a new Client.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

type apiItem struct {
	Title         string `json:"title"`
	Author        string `json:"author"`
	Flavor        string `json:"flavor"`
	PermalinkFrag string `json:"permalinkFragment"`
	Upvotes       int    `json:"upvotes"`
	Downvotes     int    `json:"downvotes"`
	DateCreated   string `json:"dateCreated"`
}

type apiResponse struct {
	Data       []apiItem `json:"data"`
	NextCursor string    `json:"nextCursor"`
	HasMore    bool      `json:"hasMore"`
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	if c.cfg.Rate > 0 {
		if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.cfg.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	var body []byte
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		resp, err := c.http.Do(req)
		c.last = time.Now()
		if err != nil {
			if attempt < c.cfg.Retries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, err
		}
		body, err = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			if attempt < c.cfg.Retries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		break
	}
	return body, nil
}

func itemToPattern(rank int, d apiItem, baseURL string) Pattern {
	author := strings.TrimSpace(d.Author)
	if len(author) > 50 {
		author = author[:50]
	}
	created := d.DateCreated
	if i := strings.Index(created, "T"); i > 0 {
		created = created[:i]
	}
	return Pattern{
		Rank:        rank,
		Title:       d.Title,
		Flavor:      d.Flavor,
		Author:      author,
		Upvotes:     d.Upvotes,
		Downvotes:   d.Downvotes,
		DateCreated: created,
		URL:         baseURL + "/r/" + d.PermalinkFrag,
	}
}

// List returns top regex patterns sorted by upvotes, optionally filtered by flavor.
func (c *Client) List(ctx context.Context, flavor string, limit int) ([]Pattern, error) {
	if limit <= 0 {
		limit = 50
	}
	var out []Pattern
	cursor := ""
	for len(out) < limit {
		var path string
		if cursor != "" {
			p := url.Values{}
			p.Set("orderBy", "MOST_UPVOTES")
			p.Set("cursor", cursor)
			if flavor != "" {
				p.Set("flavor", flavor)
			}
			path = "/api/library?" + p.Encode()
		} else {
			p := url.Values{}
			p.Set("orderBy", "MOST_UPVOTES")
			if flavor != "" {
				p.Set("flavor", flavor)
			}
			path = "/api/library?" + p.Encode()
		}
		data, err := c.get(ctx, path)
		if err != nil {
			return nil, err
		}
		var resp apiResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("parse response: %w", err)
		}
		for _, d := range resp.Data {
			out = append(out, itemToPattern(len(out)+1, d, c.cfg.BaseURL))
			if len(out) >= limit {
				break
			}
		}
		if !resp.HasMore || resp.NextCursor == "" {
			break
		}
		cursor = resp.NextCursor
	}
	return out, nil
}

// Search returns regex patterns matching query, optionally filtered by flavor.
func (c *Client) Search(ctx context.Context, query, flavor string, limit int) ([]Pattern, error) {
	if limit <= 0 {
		limit = 20
	}
	p := url.Values{}
	p.Set("orderBy", "RELEVANCE")
	p.Set("search", query)
	if flavor != "" {
		p.Set("flavor", flavor)
	}

	data, err := c.get(ctx, "/api/library?"+p.Encode())
	if err != nil {
		return nil, err
	}
	var resp apiResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	var out []Pattern
	for _, d := range resp.Data {
		out = append(out, itemToPattern(len(out)+1, d, c.cfg.BaseURL))
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}
