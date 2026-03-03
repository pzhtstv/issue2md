package convert

import (
	"testing"
	"time"

	"github.com/pzhtstv/issue2md/internal/github"
)

func TestConvertIssue(t *testing.T) {
	tests := []struct {
		name    string
		data    *github.IssueData
		options []ConverterOption
		want    string
	}{
		{
			name: "basic issue",
			data: &github.IssueData{
				Title:     "Test Issue Title",
				Body:      "This is the issue body content.",
				Author:    github.Author{Login: "testuser"},
				CreatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				State:     "open",
				Number:    123,
				Comments:  []github.Comment{},
			},
			options: []ConverterOption{},
			want: `# Test Issue Title

**Author:** @testuser | **Created:** 2024-01-15 10:30:00 | **Status:** Open

---

This is the issue body content.

---

## Comments (0)

`,
		},
		{
			name: "issue with comments",
			data: &github.IssueData{
				Title:     "Issue With Comments",
				Body:      "Main issue body.",
				Author:    github.Author{Login: "owner"},
				CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				State:     "closed",
				Number:    456,
				Comments: []github.Comment{
					{
						ID:        1,
						Body:      "First comment",
						Author:    github.Author{Login: "user1"},
						CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:        2,
						Body:      "Second comment",
						Author:    github.Author{Login: "user2"},
						CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			options: []ConverterOption{},
			want: `# Issue With Comments

**Author:** @owner | **Created:** 2024-01-01 00:00:00 | **Status:** Closed

---

Main issue body.

---

## Comments (2)

### Comment #1 | 2024-01-02

First comment

### Comment #2 | 2024-01-03

Second comment

`,
		},
		{
			name: "issue with user links",
			data: &github.IssueData{
				Title:     "Test Issue",
				Body:      "Body mentioning @otheruser",
				Author:    github.Author{Login: "author"},
				CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				State:     "open",
				Number:    1,
				Comments:  []github.Comment{},
			},
			options: []ConverterOption{WithUserLinks(true)},
			want: `# Test Issue

**Author:** [author](https://github.com/author) | **Created:** 2024-01-01 00:00:00 | **Status:** Open

---

Body mentioning @otheruser

---

## Comments (0)

`,
		},
		{
			name: "closed issue",
			data: &github.IssueData{
				Title:     "Closed Issue",
				Body:      "Body",
				Author:    github.Author{Login: "user"},
				CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				State:     "closed",
				Number:    100,
				Comments:  []github.Comment{},
			},
			options: []ConverterOption{},
			want: `# Closed Issue

**Author:** @user | **Created:** 2024-01-01 00:00:00 | **Status:** Closed

---

Body

---

## Comments (0)

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := New(tt.options...)
			got, err := converter.ConvertIssue(tt.data)
			if err != nil {
				t.Fatalf("ConvertIssue() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("ConvertIssue() = %q, want %q", got, tt.want)
			}
		})
	}
}
