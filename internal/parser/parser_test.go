package parser

import (
	"testing"
)

func TestURLParser(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantType URLType
		wantOwner string
		wantRepo  string
		wantNum  int
		wantErr  bool
	}{
		{
			name:      "valid issue URL",
			url:       "https://github.com/golang/go/issues/123",
			wantType:  TypeIssue,
			wantOwner: "golang",
			wantRepo:  "go",
			wantNum:   123,
			wantErr:   false,
		},
		{
			name:      "valid pull request URL",
			url:       "https://github.com/golang/go/pull/456",
			wantType:  TypePullRequest,
			wantOwner: "golang",
			wantRepo:  "go",
			wantNum:   456,
			wantErr:   false,
		},
		{
			name:      "valid discussion URL",
			url:       "https://github.com/golang/go/discussions/789",
			wantType:  TypeDiscussion,
			wantOwner: "golang",
			wantRepo:  "go",
			wantNum:   789,
			wantErr:   false,
		},
		{
			name:      "invalid URL - random string",
			url:       "not-a-valid-url",
			wantType:  TypeUnknown,
			wantOwner: "",
			wantRepo:  "",
			wantNum:   0,
			wantErr:   true,
		},
		{
			name:      "invalid URL - repository home",
			url:       "https://github.com/golang/go",
			wantType:  TypeUnknown,
			wantOwner: "",
			wantRepo:  "",
			wantNum:   0,
			wantErr:   true,
		},
		{
			name:      "invalid URL - blob",
			url:       "https://github.com/golang/go/blob/main/README.md",
			wantType:  TypeUnknown,
			wantOwner: "",
			wantRepo:  "",
			wantNum:   0,
			wantErr:   true,
		},
		{
			name:      "invalid URL - tree",
			url:       "https://github.com/golang/go/tree/main",
			wantType:  TypeUnknown,
			wantOwner: "",
			wantRepo:  "",
			wantNum:   0,
			wantErr:   true,
		},
		{
			name:      "valid issue URL with www",
			url:       "https://www.github.com/golang/go/issues/100",
			wantType:  TypeIssue,
			wantOwner: "golang",
			wantRepo:  "go",
			wantNum:   100,
			wantErr:   false,
		},
		{
			name:      "valid PR URL without https",
			url:       "github.com/golang/go/pull/200",
			wantType:  TypePullRequest,
			wantOwner: "golang",
			wantRepo:  "go",
			wantNum:   200,
			wantErr:   false,
		},
		{
			name:      "invalid - negative number",
			url:       "https://github.com/golang/go/issues/-1",
			wantType:  TypeUnknown,
			wantOwner: "",
			wantRepo:  "",
			wantNum:   0,
			wantErr:   true,
		},
		{
			name:      "invalid - zero number",
			url:       "https://github.com/golang/go/issues/0",
			wantType:  TypeUnknown,
			wantOwner: "",
			wantRepo:  "",
			wantNum:   0,
			wantErr:   true,
		},
	}

	p := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.Type != tt.wantType {
				t.Errorf("Parse().Type = %v, want %v", got.Type, tt.wantType)
			}
			if got.Owner != tt.wantOwner {
				t.Errorf("Parse().Owner = %v, want %v", got.Owner, tt.wantOwner)
			}
			if got.Repo != tt.wantRepo {
				t.Errorf("Parse().Repo = %v, want %v", got.Repo, tt.wantRepo)
			}
			if got.Number != tt.wantNum {
				t.Errorf("Parse().Number = %v, want %v", got.Number, tt.wantNum)
			}
		})
	}
}
