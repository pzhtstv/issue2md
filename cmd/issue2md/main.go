package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pzhtstv/issue2md/internal/convert"
	"github.com/pzhtstv/issue2md/internal/github"
	"github.com/pzhtstv/issue2md/internal/parser"
)

var (
	version   = "dev"
	commit    = ""
	date      = ""
)

func main() {
	// Define flags
	output := flag.String("o", "", "output file path")
	outputAlt := flag.String("output", "", "output file path")
	includeReactions := flag.Bool("include-reactions", false, "include reactions statistics")
	userLinks := flag.Bool("user-links", false, "render username as GitHub link")
	token := flag.String("token", "", "GitHub token (optional, can also use GITHUB_TOKEN env)")
	showVersion := flag.Bool("v", false, "show version")
	showVersionAlt := flag.Bool("version", false, "show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "issue2md - Convert GitHub Issue/PR/Discussion to Markdown\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s <github-url> [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Handle version flag
	if *showVersion || *showVersionAlt {
		fmt.Fprintf(os.Stdout, "issue2md version %s", version)
		if commit != "" {
			fmt.Fprintf(os.Stdout, " (commit: %s)", commit)
		}
		if date != "" {
			fmt.Fprintf(os.Stdout, " (date: %s)", date)
		}
		fmt.Fprintln(os.Stdout)
		os.Exit(0)
	}

	// Get URL from arguments
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "error: please provide a GitHub URL")
		flag.Usage()
		os.Exit(1)
	}

	url := args[0]

	// Use output flag (support both -o and --output)
	outFile := *output
	if outFile == "" {
		outFile = *outputAlt
	}

	// Get token from flag or environment variable
	githubToken := *token
	if githubToken == "" {
		githubToken = os.Getenv("GITHUB_TOKEN")
	}

	// Parse URL
	p := parser.New()
	parsed, err := p.Parse(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	// Create GitHub client
	client := github.New(github.WithToken(githubToken))

	// Fetch data based on URL type
	var content string
	switch strings.ToLower(string(parsed.Type)) {
	case "issue":
		data, err := client.FetchIssue(parsed.Owner, parsed.Repo, parsed.Number)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		md := convert.New(convert.WithUserLinks(*userLinks), convert.WithIncludeReactions(*includeReactions))
		content, err = md.ConvertIssue(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "pull":
		data, err := client.FetchPullRequest(parsed.Owner, parsed.Repo, parsed.Number)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		md := convert.New(convert.WithUserLinks(*userLinks), convert.WithIncludeReactions(*includeReactions))
		content, err = md.ConvertPullRequest(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "discussion":
		data, err := client.FetchDiscussion(parsed.Owner, parsed.Repo, parsed.Number)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		md := convert.New(convert.WithUserLinks(*userLinks), convert.WithIncludeReactions(*includeReactions))
		content, err = md.ConvertDiscussion(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "error: unsupported URL type: %s\n", parsed.Type)
		os.Exit(2)
	}

	// Output
	if outFile != "" {
		err = os.WriteFile(outFile, []byte(content), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to write file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stdout, content)
	}
}
