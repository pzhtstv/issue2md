package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v45/github"
)

func TestFetchIssue(t *testing.T) {
	tests := []struct {
		name       string
		owner      string
		repo       string
		number     int
		statusCode int
		issueResp  interface{}
		commentsResp interface{}
		wantErr    bool
		errType    error
	}{
		{
			name:   "successful fetch",
			owner:  "golang",
			repo:   "go",
			number: 123,
			statusCode: http.StatusOK,
			issueResp: &github.Issue{
				Title:  github.String("Test Issue"),
				Body:   github.String("Issue body content"),
				State:  github.String("open"),
				Number: github.Int(123),
				User:   &github.User{Login: github.String("testuser")},
			},
			commentsResp: []*github.IssueComment{},
			wantErr: false,
		},
		{
			name:   "not found",
			owner:  "golang",
			repo:   "go",
			number: 999999999,
			statusCode: http.StatusNotFound,
			issueResp: map[string]string{
				"message": "Not Found",
			},
			commentsResp: nil,
			wantErr: true,
			errType: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server with multiple endpoints
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Handle comments endpoint
				if r.URL.Path == "/repos/golang/go/issues/123/comments" {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(tt.commentsResp)
					return
				}
				// Handle issue endpoint
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.issueResp)
			}))
			defer server.Close()

			// Create client pointing to test server
			client := &client{
				httpClient: github.NewClient(server.Client()),
			}
			client.httpClient.BaseURL, _ = client.httpClient.BaseURL.Parse(server.URL + "/")

			got, err := client.FetchIssue(tt.owner, tt.repo, tt.number)

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchIssue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errType != nil {
				if err != tt.errType {
					t.Errorf("FetchIssue() error type = %v, want %v", err, tt.errType)
				}
				return
			}

			if !tt.wantErr && got == nil {
				t.Error("FetchIssue() got = nil, want non-nil")
			}
		})
	}
}

func TestFetchPullRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/golang/go/issues/456/comments" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]*github.IssueComment{})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&github.PullRequest{
			Title:  github.String("Test PR"),
			Body:   github.String("PR body content"),
			State:  github.String("open"),
			Number: github.Int(456),
			User:   &github.User{Login: github.String("testuser")},
			Merged: github.Bool(false),
		})
	}))
	defer server.Close()

	client := &client{
		httpClient: github.NewClient(server.Client()),
	}
	client.httpClient.BaseURL, _ = client.httpClient.BaseURL.Parse(server.URL + "/")

	got, err := client.FetchPullRequest("golang", "go", 456)
	if err != nil {
		t.Errorf("FetchPullRequest() error = %v", err)
	}
	if got == nil {
		t.Error("FetchPullRequest() got = nil, want non-nil")
	}
}

func TestFetchDiscussion(t *testing.T) {
	// Discussions require GraphQL API
	client := &client{}

	_, err := client.FetchDiscussion("owner", "repo", 1)
	if err == nil {
		t.Error("FetchDiscussion() want error, got nil")
	}
}

// Test token authentication
func TestNewWithToken(t *testing.T) {
	client := New(WithToken("test-token"))
	if client == nil {
		t.Error("New() returned nil")
	}
}

// Test error handling
func TestHandleError(t *testing.T) {
	client := &client{}

	tests := []struct {
		name      string
		statusCode int
		wantErr   error
	}{
		{"not found", http.StatusNotFound, ErrNotFound},
		{"forbidden", http.StatusForbidden, ErrRateLimited},
		{"unauthorized", http.StatusUnauthorized, ErrUnauthorized},
		{"nil response", 0, ErrNetwork},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *github.Response
			if tt.statusCode > 0 {
				resp = &github.Response{Response: &http.Response{StatusCode: tt.statusCode}}
			}
			err := client.handleError(resp)
			if err != tt.wantErr {
				t.Errorf("handleError() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
