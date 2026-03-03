package convert

import (
	"strconv"
	"strings"

	"github.com/pzhtstv/issue2md/internal/github"
)

// ConvertIssue converts an IssueData to Markdown
func (c *Converter) ConvertIssue(data *github.IssueData) (string, error) {
	// Pre-calculate approximate size for better performance
	// Title + body + comments estimate
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
	sb.WriteString(formatState(data.State))
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

// formatAuthor formats the author based on options
func (c *Converter) formatAuthor(author github.Author) string {
	if c.userLinks {
		// Pre-allocate for link format
		var sb strings.Builder
		sb.Grow(len(author.Login) + 25)
		sb.WriteString("[")
		sb.WriteString(author.Login)
		sb.WriteString("](https://github.com/")
		sb.WriteString(author.Login)
		sb.WriteString(")")
		return sb.String()
	}
	return "@" + author.Login
}

// formatState capitalizes the first letter of state
func formatState(state string) string {
	if state == "" {
		return ""
	}
	if len(state) == 1 {
		return strings.ToUpper(state)
	}
	return strings.ToUpper(state[:1]) + state[1:]
}
