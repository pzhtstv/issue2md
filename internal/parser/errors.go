package parser

import "errors"

// ErrInvalidURL indicates the URL is not a valid GitHub URL
var ErrInvalidURL = errors.New("invalid GitHub URL")

// ErrUnsupportedURLType indicates the URL type is not supported
var ErrUnsupportedURLType = errors.New("unsupported URL type")
