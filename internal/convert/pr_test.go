package convert

import (
	"testing"
	"time"

	"github.com/pzhtstv/issue2md/internal/github"
)

func TestConvertPullRequest(t *testing.T) {
	tests := []struct {
		name    string
		data    *github.PullRequestData
		options []ConverterOption
		want    string
	}{
		{
			name: "basic PR",
			data: &github.PullRequestData{
				Title:     "Feature: Add new feature",
				Body:      "This PR adds a new feature.",
				Author:    github.Author{Login: "developer"},
				CreatedAt: time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
				State:     "open",
				Number:    100,
				Comments:  []github.Comment{},
			},
			options: []ConverterOption{},
			want: `# Feature: Add new feature

**Author:** @developer | **Created:** 2024-01-10 09:00:00 | **Status:** Open

---

This PR adds a new feature.

---

## Comments (0)

`,
		},
		{
			name: "merged PR",
			data: &github.PullRequestData{
				Title:     "PR: Fix bug",
				Body:      "Bug fix PR.",
				Author:    github.Author{Login: "bugfixer"},
				CreatedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				State:     "closed",
				Merged:    true,
				Number:    50,
				Comments: []github.Comment{
					{
						ID:        1,
						Body:      "LGTM!",
						Author:    github.Author{Login: "reviewer"},
						CreatedAt: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			options: []ConverterOption{},
			want: `# PR: Fix bug

**Author:** @bugfixer | **Created:** 2024-01-05 00:00:00 | **Status:** Merged

---

Bug fix PR.

---

## Comments (1)

### Comment #1 | 2024-01-06

LGTM!

`,
		},
		{
			name: "PR with user links",
			data: &github.PullRequestData{
				Title:     "PR Title",
				Body:      "Body by @contributor",
				Author:    github.Author{Login: "author"},
				CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				State:     "open",
				Number:    1,
				Comments:  []github.Comment{},
			},
			options: []ConverterOption{WithUserLinks(true)},
			want: `# PR Title

**Author:** [author](https://github.com/author) | **Created:** 2024-01-01 00:00:00 | **Status:** Open

---

Body by @contributor

---

## Comments (0)

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := New(tt.options...)
			got, err := converter.ConvertPullRequest(tt.data)
			if err != nil {
				t.Fatalf("ConvertPullRequest() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("ConvertPullRequest() = %q, want %q", got, tt.want)
			}
		})
	}
}
