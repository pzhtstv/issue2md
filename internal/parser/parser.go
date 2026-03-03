package parser

// Parser parses GitHub URLs
type Parser interface {
	Parse(rawURL string) (*ParsedURL, error)
}
