package github

import "errors"

// ErrNotFound indicates the requested resource was not found
var ErrNotFound = errors.New("resource not found")

// ErrRateLimited indicates the API rate limit has been exceeded
var ErrRateLimited = errors.New("API rate limit exceeded")

// ErrUnauthorized indicates authentication is required
var ErrUnauthorized = errors.New("unauthorized")

// ErrNetwork indicates a network error occurred
var ErrNetwork = errors.New("network error")
