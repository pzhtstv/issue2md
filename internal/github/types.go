package github

import "time"

// Author represents a GitHub user
type Author struct {
	Login string
	URL   string
}

// Reaction represents a GitHub reaction
type Reaction struct {
	Content string // "+1", "-1", "heart", "hooray", "laugh", "confused"
	Count   int
}

// Comment represents a GitHub comment
type Comment struct {
	ID        int
	Body      string
	Author    Author
	CreatedAt time.Time
	Reactions []Reaction
}

// IssueData represents a GitHub issue
type IssueData struct {
	Title     string
	Body      string
	Author    Author
	CreatedAt time.Time
	State     string // "open" | "closed"
	Number    int
	Comments  []Comment
	Reactions []Reaction
}

// PullRequestData represents a GitHub pull request
type PullRequestData struct {
	Title     string
	Body      string
	Author    Author
	CreatedAt time.Time
	State     string // "open" | "closed" | "merged"
	Merged    bool
	Number    int
	Comments  []Comment
	Reactions []Reaction
}

// DiscussionData represents a GitHub discussion
type DiscussionData struct {
	Title     string
	Body      string
	Author    Author
	CreatedAt time.Time
	Category  string
	Number    int
	Answers   []Comment
	Comments  []Comment
	Reactions []Reaction
}
