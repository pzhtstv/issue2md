package parser

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	prefixHTTPS   = "https://"
	prefixHTTP    = "http://"
	prefixWWW     = "www."
	minPartsCount = 5
)

// urlParser is the default implementation of Parser
type urlParser struct{}

// New creates a new Parser
func New() Parser {
	return &urlParser{}
}

// Parse parses a GitHub URL and returns a ParsedURL
func (p *urlParser) Parse(rawURL string) (*ParsedURL, error) {
	// Remove protocol prefix with single operation
	url := rawURL
	if len(url) > len(prefixHTTPS) {
		if url[:len(prefixHTTPS)] == prefixHTTPS {
			url = url[len(prefixHTTPS):]
		} else if len(url) > len(prefixHTTP) && url[:len(prefixHTTP)] == prefixHTTP {
			url = url[len(prefixHTTP):]
		}
	}

	// Handle www prefix
	if len(url) > len(prefixWWW) && url[:len(prefixWWW)] == prefixWWW {
		url = url[len(prefixWWW):]
	}

	// Split by /
	parts := strings.SplitN(url, "/", minPartsCount)
	// Expected format: github.com/owner/repo/type/number (5 parts minimum)
	if len(parts) < minPartsCount {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, rawURL)
	}

	// parts[0] = "github.com" (hostname, ignore)
	// parts[1] = owner
	// parts[2] = repo
	// parts[3] = type (issues/pull/discussions)
	// parts[4] = number
	owner := parts[1]
	repo := parts[2]
	pathType := parts[3]
	numberStr := parts[4]

	// Validate owner and repo (should not be empty)
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, rawURL)
	}

	// Parse number
	number, err := strconv.Atoi(numberStr)
	if err != nil || number <= 0 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, rawURL)
	}

	// Determine URL type
	var urlType URLType
	switch pathType {
	case "issues":
		urlType = TypeIssue
	case "pull":
		urlType = TypePullRequest
	case "discussions":
		urlType = TypeDiscussion
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedURLType, rawURL)
	}

	return &ParsedURL{
		Type:   urlType,
		Owner:  owner,
		Repo:   repo,
		Number: number,
		RawURL: rawURL,
	}, nil
}
