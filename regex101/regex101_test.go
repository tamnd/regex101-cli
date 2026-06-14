package regex101_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/regex101-cli/regex101"
)

func mockResp(items []map[string]any, hasMore bool, cursor string) string {
	resp := map[string]any{
		"data":       items,
		"hasMore":    hasMore,
		"nextCursor": cursor,
	}
	b, _ := json.Marshal(resp)
	return string(b)
}

func makeItems(n int) []map[string]any {
	var items []map[string]any
	for i := 1; i <= n; i++ {
		items = append(items, map[string]any{
			"title":             fmt.Sprintf("Pattern %d", i),
			"author":            "tester",
			"flavor":            "javascript",
			"permalinkFragment": fmt.Sprintf("abc%d", i),
			"upvotes":           100 - i,
			"downvotes":         i,
			"dateCreated":       "2023-01-01T00:00:00.000Z",
		})
	}
	return items
}

func newTestClient(ts *httptest.Server) *regex101.Client {
	cfg := regex101.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return regex101.NewClient(cfg)
}

func TestList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, mockResp(makeItems(5), false, ""))
	}))
	defer ts.Close()

	patterns, err := newTestClient(ts).List(context.Background(), "", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(patterns) != 5 {
		t.Fatalf("got %d patterns, want 5", len(patterns))
	}
	if patterns[0].Title != "Pattern 1" {
		t.Errorf("patterns[0].Title = %q", patterns[0].Title)
	}
	if patterns[0].URL == "" {
		t.Error("URL should not be empty")
	}
}

func TestListPagination(t *testing.T) {
	page := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page++
		if page == 1 {
			_, _ = fmt.Fprint(w, mockResp(makeItems(3), true, "cursor123"))
		} else {
			_, _ = fmt.Fprint(w, mockResp(makeItems(3), false, ""))
		}
	}))
	defer ts.Close()

	patterns, err := newTestClient(ts).List(context.Background(), "", 6)
	if err != nil {
		t.Fatal(err)
	}
	if len(patterns) != 6 {
		t.Fatalf("got %d patterns, want 6", len(patterns))
	}
}

func TestSearch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("search")
		if q != "email" {
			t.Errorf("unexpected search query %q", q)
		}
		_, _ = fmt.Fprint(w, mockResp(makeItems(3), false, ""))
	}))
	defer ts.Close()

	patterns, err := newTestClient(ts).Search(context.Background(), "email", "", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(patterns) != 3 {
		t.Fatalf("got %d patterns, want 3", len(patterns))
	}
}
