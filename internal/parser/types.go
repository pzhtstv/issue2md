package parser

// URLType represents the type of GitHub resource
type URLType string

const (
	TypeIssue        URLType = "issue"
	TypePullRequest  URLType = "pull"
	TypeDiscussion   URLType = "discussion"
	TypeUnknown      URLType = "unknown"
)

// ParsedURL represents a parsed GitHub URL
type ParsedURL struct {
	Type   URLType
	Owner  string
	Repo   string
	Number int
	RawURL string
}
