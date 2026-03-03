package convert

import (
	"strconv"
	"strings"

	"github.com/pzhtstv/issue2md/internal/github"
)

// ConvertPullRequest converts a PullRequestData to Markdown
func (c *Converter) ConvertPullRequest(data *github.PullRequestData) (string, error) {
	// Pre-calculate approximate size for better performance
	estimatedSize := len(data.Title) + len(data.Body) + len(data.Comments)*200
	if estimatedSize < 1024 {
		estimatedSize = 1024
	}

	var sb strings.Builder
	sb.Grow(estimatedSize)

	// Title
	sb.WriteString("# ")
	sb.WriteString(data.Title)
	sb.WriteString("\n\n")

	// Author, Created, Status
	sb.WriteString("**Author:** ")
	sb.WriteString(c.formatAuthor(data.Author))
	sb.WriteString(" | **Created:** ")
	sb.WriteString(data.CreatedAt.Format("2006-01-02 15:04:05"))
	sb.WriteString(" | **Status:** ")
	sb.WriteString(c.formatPRState(data.State, data.Merged))
	sb.WriteString("\n\n")

	// Divider
	sb.WriteString("---\n\n")

	// Body
	sb.WriteString(data.Body)
	sb.WriteString("\n\n")

	// Divider
	sb.WriteString("---\n\n")

	// Comments header
	sb.WriteString("## Comments (")
	sb.WriteString(strconv.Itoa(len(data.Comments)))
	sb.WriteString(")\n\n")

	// Comments
	for i, comment := range data.Comments {
		sb.WriteString("### Comment #")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString(" | ")
		sb.WriteString(comment.CreatedAt.Format("2006-01-02"))
		sb.WriteString("\n\n")
		sb.WriteString(comment.Body)
		sb.WriteString("\n\n")
	}

	return sb.String(), nil
}

// formatPRState formats the PR state, showing "Merged" if merged
func (c *Converter) formatPRState(state string, merged bool) string {
	if merged {
		return "Merged"
	}
	if len(state) == 1 {
		return strings.ToUpper(state)
	}
	return strings.ToUpper(state[:1]) + state[1:]
}
