package convert

import (
	"testing"
	"time"

	"github.com/pzhtstv/issue2md/internal/github"
)

func TestConvertDiscussion(t *testing.T) {
	tests := []struct {
		name    string
		data    *github.DiscussionData
		options []ConverterOption
		want    string
	}{
		{
			name: "basic discussion",
			data: &github.DiscussionData{
				Title:     "How to use this feature?",
				Body:      "I have a question about using this feature.",
				Author:    github.Author{Login: "questioner"},
				CreatedAt: time.Date(2024, 2, 1, 12, 0, 0, 0, time.UTC),
				Category:  "Q&A",
				Number:    10,
				Answers:   []github.Comment{},
				Comments:  []github.Comment{},
			},
			options: []ConverterOption{},
			want: `# How to use this feature?

**Author:** @questioner | **Created:** 2024-02-01 12:00:00 | **Category:** Q&A

---

I have a question about using this feature.

---

## Answers (0)

## Comments (0)

`,
		},
		{
			name: "discussion with answers",
			data: &github.DiscussionData{
				Title:     "Discussion with Answers",
				Body:      "Main discussion body.",
				Author:    github.Author{Login: "author"},
				CreatedAt: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				Category:  "General",
				Number:    20,
				Answers: []github.Comment{
					{
						ID:        1,
						Body:      "This is the answer.",
						Author:    github.Author{Login: "helper"},
						CreatedAt: time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
					},
				},
				Comments: []github.Comment{
					{
						ID:        2,
						Body:      "Follow-up comment",
						Author:    github.Author{Login: "commenter"},
						CreatedAt: time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			options: []ConverterOption{},
			want: `# Discussion with Answers

**Author:** @author | **Created:** 2024-02-01 00:00:00 | **Category:** General

---

Main discussion body.

---

## Answers (1)

### Answer by @helper | 2024-02-02

This is the answer.

## Comments (1)

### Comment #1 | 2024-02-03

Follow-up comment

`,
		},
		{
			name: "discussion with user links",
			data: &github.DiscussionData{
				Title:     "Question",
				Body:      "Asked by @user1",
				Author:    github.Author{Login: "author"},
				CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Category:  "Help",
				Number:    1,
				Answers:   []github.Comment{},
				Comments:  []github.Comment{},
			},
			options: []ConverterOption{WithUserLinks(true)},
			want: `# Question

**Author:** [author](https://github.com/author) | **Created:** 2024-01-01 00:00:00 | **Category:** Help

---

Asked by @user1

---

## Answers (0)

## Comments (0)

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := New(tt.options...)
			got, err := converter.ConvertDiscussion(tt.data)
			if err != nil {
				t.Fatalf("ConvertDiscussion() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("ConvertDiscussion() = %q, want %q", got, tt.want)
			}
		})
	}
}
