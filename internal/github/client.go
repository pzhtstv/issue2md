package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

// defaultContext is used for API calls
var defaultCtx = context.Background()

// client implements GitHuber interface using go-github library
type client struct {
	httpClient *github.Client
	token     string
}

// Option configures the GitHub client
type Option func(*client)

// WithToken sets the GitHub token
func WithToken(token string) Option {
	return func(c *client) {
		c.token = token
	}
}

// New creates a new GitHub client
func New(opts ...Option) GitHuber {
	c := &client{}
	for _, opt := range opts {
		opt(c)
	}

	var hc *http.Client
	if c.token != "" {
		// Use OAuth2 for token authentication
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.token},
		)
		tc := oauth2.NewClient(defaultCtx, ts)
		tc.Timeout = 30 * time.Second
		hc = tc
	} else {
		hc = &http.Client{Timeout: 30 * time.Second}
	}

	c.httpClient = github.NewClient(hc)
	return c
}

// FetchIssue fetches an issue from GitHub
func (c *client) FetchIssue(owner, repo string, number int) (*IssueData, error) {
	issue, resp, err := c.httpClient.Issues.Get(defaultCtx, owner, repo, number)
	if err != nil {
		return nil, c.handleError(resp)
	}

	// Fetch comments
	comments, _, err := c.httpClient.Issues.ListComments(defaultCtx, owner, repo, number, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	// Convert to IssueData
	data := &IssueData{
		Title:     issue.GetTitle(),
		Body:      issue.GetBody(),
		Author:    Author{Login: issue.GetUser().GetLogin()},
		CreatedAt: issue.GetCreatedAt(),
		State:     issue.GetState(),
		Number:    issue.GetNumber(),
	}

	// Convert comments
	for _, comment := range comments {
		data.Comments = append(data.Comments, Comment{
			ID:        int(comment.GetID()),
			Body:      comment.GetBody(),
			Author:    Author{Login: comment.GetUser().GetLogin()},
			CreatedAt: comment.GetCreatedAt(),
		})
	}

	return data, nil
}

// FetchPullRequest fetches a pull request from GitHub
func (c *client) FetchPullRequest(owner, repo string, number int) (*PullRequestData, error) {
	pr, resp, err := c.httpClient.PullRequests.Get(defaultCtx, owner, repo, number)
	if err != nil {
		return nil, c.handleError(resp)
	}

	// Fetch comments (issue comments)
	comments, _, err := c.httpClient.Issues.ListComments(defaultCtx, owner, repo, number, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	data := &PullRequestData{
		Title:     pr.GetTitle(),
		Body:      pr.GetBody(),
		Author:    Author{Login: pr.GetUser().GetLogin()},
		CreatedAt: pr.GetCreatedAt(),
		State:     pr.GetState(),
		Merged:    pr.GetMerged(),
		Number:    pr.GetNumber(),
	}

	for _, comment := range comments {
		data.Comments = append(data.Comments, Comment{
			ID:        int(comment.GetID()),
			Body:      comment.GetBody(),
			Author:    Author{Login: comment.GetUser().GetLogin()},
			CreatedAt: comment.GetCreatedAt(),
		})
	}

	return data, nil
}

// FetchDiscussion fetches a discussion from GitHub
// Note: go-github v45 does not support GraphQL Discussions API
// This is a placeholder that returns an error
func (c *client) FetchDiscussion(owner, repo string, number int) (*DiscussionData, error) {
	return nil, errors.New("discussions require GraphQL API (not implemented in REST client)")
}

// handleError converts GitHub API error responses to custom errors
func (c *client) handleError(resp *github.Response) error {
	if resp == nil {
		return ErrNetwork
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusForbidden:
		return ErrRateLimited
	case http.StatusUnauthorized:
		return ErrUnauthorized
	default:
		return fmt.Errorf("API error: status %d", resp.StatusCode)
	}
}
