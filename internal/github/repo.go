package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/v68/github"
)

// RepoInfo holds basic repository metadata.
type RepoInfo struct {
	FullName      string
	Description   string
	DefaultBranch string
	Private       bool
	HTMLURL       string
	Stars         int
	Forks         int
}

// RepoEntry represents a file or directory in a repository tree listing.
type RepoEntry struct {
	Name    string
	Path    string
	Type    string // "file" or "dir"
	Size    int
	HTMLURL string
}

// RepoFile holds the content of a single file fetched from GitHub.
type RepoFile struct {
	Name     string
	Path     string
	Size     int
	Content  string
	HTMLURL  string
	Encoding string
}

// Branch holds branch metadata.
type Branch struct {
	Name      string
	Protected bool
	SHA       string
}

// FetchRepoTree returns the directory listing at the given path and ref.
func FetchRepoTree(ctx context.Context, client *gh.Client, owner, repo, ref, path string) ([]RepoEntry, error) {
	opts := &gh.RepositoryContentGetOptions{Ref: ref}
	_, dirContents, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return nil, fmt.Errorf("fetch repo tree: %w", err)
	}
	entries := make([]RepoEntry, 0, len(dirContents))
	for _, c := range dirContents {
		entries = append(entries, RepoEntry{
			Name:    c.GetName(),
			Path:    c.GetPath(),
			Type:    c.GetType(),
			Size:    c.GetSize(),
			HTMLURL: c.GetHTMLURL(),
		})
	}
	return entries, nil
}

// FetchFileContent returns the decoded content of a single file.
func FetchFileContent(ctx context.Context, client *gh.Client, owner, repo, ref, path string) (*RepoFile, error) {
	opts := &gh.RepositoryContentGetOptions{Ref: ref}
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return nil, fmt.Errorf("fetch file content: %w", err)
	}
	if fileContent == nil {
		return nil, fmt.Errorf("path is a directory, not a file")
	}
	content, err := fileContent.GetContent()
	if err != nil {
		return nil, fmt.Errorf("decode file content: %w", err)
	}
	return &RepoFile{
		Name:     fileContent.GetName(),
		Path:     fileContent.GetPath(),
		Size:     fileContent.GetSize(),
		Content:  content,
		HTMLURL:  fileContent.GetHTMLURL(),
		Encoding: fileContent.GetEncoding(),
	}, nil
}

// FetchBranches lists branches for a repository.
func FetchBranches(ctx context.Context, client *gh.Client, owner, repo string) ([]Branch, error) {
	opts := &gh.BranchListOptions{
		ListOptions: gh.ListOptions{PerPage: 100},
	}
	ghBranches, _, err := client.Repositories.ListBranches(ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("fetch branches: %w", err)
	}
	branches := make([]Branch, 0, len(ghBranches))
	for _, b := range ghBranches {
		branches = append(branches, Branch{
			Name:      b.GetName(),
			Protected: b.GetProtected(),
			SHA:       b.GetCommit().GetSHA(),
		})
	}
	return branches, nil
}

// FetchRepoInfo returns basic metadata for a repository.
func FetchRepoInfo(ctx context.Context, client *gh.Client, owner, repo string) (*RepoInfo, error) {
	r, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("fetch repo info: %w", err)
	}
	return &RepoInfo{
		FullName:      r.GetFullName(),
		Description:   r.GetDescription(),
		DefaultBranch: r.GetDefaultBranch(),
		Private:       r.GetPrivate(),
		HTMLURL:       r.GetHTMLURL(),
		Stars:         r.GetStargazersCount(),
		Forks:         r.GetForksCount(),
	}, nil
}
