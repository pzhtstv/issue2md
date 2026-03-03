package convert

import (
	"strconv"
	"strings"

	"github.com/pzhtstv/issue2md/internal/github"
)

// ConvertDiscussion converts a DiscussionData to Markdown
func (c *Converter) ConvertDiscussion(data *github.DiscussionData) (string, error) {
	// Pre-calculate approximate size for better performance
	estimatedSize := len(data.Title) + len(data.Body) + (len(data.Answers)+len(data.Comments))*200
	if estimatedSize < 1024 {
		estimatedSize = 1024
	}

	var sb strings.Builder
	sb.Grow(estimatedSize)

	// Title
	sb.WriteString("# ")
	sb.WriteString(data.Title)
	sb.WriteString("\n\n")

	// Author, Created, Category
	sb.WriteString("**Author:** ")
	sb.WriteString(c.formatAuthor(data.Author))
	sb.WriteString(" | **Created:** ")
	sb.WriteString(data.CreatedAt.Format("2006-01-02 15:04:05"))
	sb.WriteString(" | **Category:** ")
	sb.WriteString(data.Category)
	sb.WriteString("\n\n")

	// Divider
	sb.WriteString("---\n\n")

	// Body
	sb.WriteString(data.Body)
	sb.WriteString("\n\n")

	// Divider
	sb.WriteString("---\n\n")

	// Answers header
	sb.WriteString("## Answers (")
	sb.WriteString(strconv.Itoa(len(data.Answers)))
	sb.WriteString(")\n\n")

	// Answers
	for _, answer := range data.Answers {
		sb.WriteString("### Answer by ")
		sb.WriteString(c.formatAuthor(answer.Author))
		sb.WriteString(" | ")
		sb.WriteString(answer.CreatedAt.Format("2006-01-02"))
		sb.WriteString("\n\n")
		sb.WriteString(answer.Body)
		sb.WriteString("\n\n")
	}

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
