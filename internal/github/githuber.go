package github

// GitHuber interface defines methods for fetching GitHub data
type GitHuber interface {
	FetchIssue(owner, repo string, number int) (*IssueData, error)
	FetchPullRequest(owner, repo string, number int) (*PullRequestData, error)
	FetchDiscussion(owner, repo string, number int) (*DiscussionData, error)
}
