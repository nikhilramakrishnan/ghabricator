package github

import (
	"context"
	"fmt"
	"io"
	"net/http"

	gh "github.com/google/go-github/v68/github"
)

// CreateGist creates a new GitHub Gist.
func CreateGist(ctx context.Context, client *gh.Client, description, filename, content string, public bool) (*Gist, error) {
	gist := &gh.Gist{
		Description: gh.Ptr(description),
		Public:      gh.Ptr(public),
		Files: map[gh.GistFilename]gh.GistFile{
			gh.GistFilename(filename): {Content: gh.Ptr(content)},
		},
	}
	created, _, err := client.Gists.Create(ctx, gist)
	if err != nil {
		return nil, fmt.Errorf("create gist: %w", err)
	}
	return convertGist(created), nil
}

// FetchGist fetches a single gist by ID.
func FetchGist(ctx context.Context, client *gh.Client, gistID string) (*Gist, error) {
	gist, _, err := client.Gists.Get(ctx, gistID)
	if err != nil {
		return nil, fmt.Errorf("fetch gist: %w", err)
	}
	return convertGist(gist), nil
}

// ListGists lists the authenticated user's recent gists.
func ListGists(ctx context.Context, client *gh.Client) ([]Gist, error) {
	opts := &gh.GistListOptions{
		ListOptions: gh.ListOptions{PerPage: 30},
	}
	gists, _, err := client.Gists.List(ctx, "", opts)
	if err != nil {
		return nil, fmt.Errorf("list gists: %w", err)
	}
	result := make([]Gist, 0, len(gists))
	for _, g := range gists {
		result = append(result, *convertGist(g))
	}
	return result, nil
}

func convertGist(g *gh.Gist) *Gist {
	gist := &Gist{
		ID:          g.GetID(),
		Description: g.GetDescription(),
		Public:      g.GetPublic(),
		HTMLURL:     g.GetHTMLURL(),
		CreatedAt:   g.GetCreatedAt().Time,
		UpdatedAt:   g.GetUpdatedAt().Time,
		Files:       make(map[string]GistFile),
	}
	if g.Owner != nil {
		gist.Owner = User{
			Login:     g.Owner.GetLogin(),
			AvatarURL: g.Owner.GetAvatarURL(),
		}
	}
	for name, f := range g.Files {
		gist.Files[string(name)] = GistFile{
			Filename: f.GetFilename(),
			Language: f.GetLanguage(),
			Content:  f.GetContent(),
			Size:     f.GetSize(),
		}
	}
	return gist
}

// FetchPR fetches pull request metadata.
func FetchPR(ctx context.Context, client *gh.Client, owner, repo string, number int) (*PullRequest, error) {
	pr, _, err := client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("fetch PR: %w", err)
	}

	result := &PullRequest{
		Number:       pr.GetNumber(),
		Title:        pr.GetTitle(),
		Body:         pr.GetBody(),
		State:        pr.GetState(),
		Draft:        pr.GetDraft(),
		Merged:       pr.GetMerged(),
		CreatedAt:    pr.GetCreatedAt().Time,
		UpdatedAt:    pr.GetUpdatedAt().Time,
		Additions:    pr.GetAdditions(),
		Deletions:    pr.GetDeletions(),
		ChangedFiles: pr.GetChangedFiles(),
		Author: User{
			Login:     pr.GetUser().GetLogin(),
			AvatarURL: pr.GetUser().GetAvatarURL(),
		},
		Head: Ref{
			Ref:  pr.GetHead().GetRef(),
			SHA:  pr.GetHead().GetSHA(),
			Repo: pr.GetHead().GetRepo().GetFullName(),
		},
		Base: Ref{
			Ref:  pr.GetBase().GetRef(),
			SHA:  pr.GetBase().GetSHA(),
			Repo: pr.GetBase().GetRepo().GetFullName(),
		},
	}

	for _, l := range pr.Labels {
		result.Labels = append(result.Labels, Label{
			Name:  l.GetName(),
			Color: l.GetColor(),
		})
	}

	reviewers, _, err := client.PullRequests.ListReviewers(ctx, owner, repo, number, nil)
	if err == nil && reviewers != nil {
		for _, u := range reviewers.Users {
			result.Reviewers = append(result.Reviewers, User{
				Login:     u.GetLogin(),
				AvatarURL: u.GetAvatarURL(),
			})
		}
	}

	return result, nil
}

// FetchDiff fetches the raw unified diff for a pull request.
func FetchDiff(ctx context.Context, client *gh.Client, owner, repo string, number int) (string, error) {
	// Use the raw diff endpoint with Accept header
	url := fmt.Sprintf("repos/%s/%s/pulls/%d", owner, repo, number)
	req, err := client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create diff request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.diff")

	resp, err := client.BareDo(ctx, req)
	if err != nil {
		return "", fmt.Errorf("fetch diff: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read diff body: %w", err)
	}
	return string(body), nil
}

// FetchReviewComment fetches a single inline review comment by ID.
func FetchReviewComment(ctx context.Context, client *gh.Client, owner, repo string, commentID int64) (*ReviewComment, error) {
	c, _, err := client.PullRequests.GetComment(ctx, owner, repo, commentID)
	if err != nil {
		return nil, fmt.Errorf("fetch review comment: %w", err)
	}
	return &ReviewComment{
		ID:   c.GetID(),
		Body: c.GetBody(),
		Path: c.GetPath(),
		Line: c.GetLine(),
		Side: c.GetSide(),
		Author: User{
			Login:     c.GetUser().GetLogin(),
			AvatarURL: c.GetUser().GetAvatarURL(),
		},
		CreatedAt: c.GetCreatedAt().Time,
		UpdatedAt: c.GetUpdatedAt().Time,
	}, nil
}

// FetchReviewComments fetches inline review comments on a PR.
func FetchReviewComments(ctx context.Context, client *gh.Client, owner, repo string, number int) ([]ReviewComment, error) {
	opts := &gh.PullRequestListCommentsOptions{
		ListOptions: gh.ListOptions{PerPage: 100},
	}
	var result []ReviewComment
	for {
		comments, resp, err := client.PullRequests.ListComments(ctx, owner, repo, number, opts)
		if err != nil {
			return nil, fmt.Errorf("fetch review comments: %w", err)
		}
		for _, c := range comments {
			rc := ReviewComment{
				ID:        c.GetID(),
				Body:      c.GetBody(),
				Path:      c.GetPath(),
				Line:      c.GetLine(),
				Side:      c.GetSide(),
				CreatedAt: c.GetCreatedAt().Time,
				UpdatedAt: c.GetUpdatedAt().Time,
				DiffHunk:  c.GetDiffHunk(),
				Author: User{
					Login:     c.GetUser().GetLogin(),
					AvatarURL: c.GetUser().GetAvatarURL(),
				},
			}
			if c.InReplyTo != nil {
				rc.InReplyTo = c.GetInReplyTo()
			}
			if c.Reactions != nil {
				rc.Reactions = &ReactionSummary{
					PlusOne:  c.Reactions.GetPlusOne(),
					MinusOne: c.Reactions.GetMinusOne(),
					Laugh:    c.Reactions.GetLaugh(),
					Confused: c.Reactions.GetConfused(),
					Heart:    c.Reactions.GetHeart(),
					Hooray:   c.Reactions.GetHooray(),
					Rocket:   c.Reactions.GetRocket(),
					Eyes:     c.Reactions.GetEyes(),
				}
			}
			result = append(result, rc)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return result, nil
}

// FetchReviews fetches review summaries (approve, request changes, etc).
func FetchReviews(ctx context.Context, client *gh.Client, owner, repo string, number int) ([]Review, error) {
	opts := &gh.ListOptions{PerPage: 100}
	var result []Review
	for {
		reviews, resp, err := client.PullRequests.ListReviews(ctx, owner, repo, number, opts)
		if err != nil {
			return nil, fmt.Errorf("fetch reviews: %w", err)
		}
		for _, r := range reviews {
			result = append(result, Review{
				ID:        r.GetID(),
				State:     r.GetState(),
				Body:      r.GetBody(),
				CreatedAt: r.GetSubmittedAt().Time,
				Author: User{
					Login:     r.GetUser().GetLogin(),
					AvatarURL: r.GetUser().GetAvatarURL(),
				},
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return result, nil
}

// CreateReviewComment creates an inline review comment.
func CreateReviewComment(ctx context.Context, client *gh.Client, owner, repo string, number int, body, path string, line int, side string) (*ReviewComment, error) {
	comment := &gh.PullRequestComment{
		Body: gh.Ptr(body),
		Path: gh.Ptr(path),
		Line: gh.Ptr(line),
		Side: gh.Ptr(side),
	}
	created, _, err := client.PullRequests.CreateComment(ctx, owner, repo, number, comment)
	if err != nil {
		return nil, fmt.Errorf("create review comment: %w", err)
	}
	return &ReviewComment{
		ID:   created.GetID(),
		Body: created.GetBody(),
		Path: created.GetPath(),
		Line: created.GetLine(),
		Side: created.GetSide(),
		Author: User{
			Login:     created.GetUser().GetLogin(),
			AvatarURL: created.GetUser().GetAvatarURL(),
		},
		CreatedAt: created.GetCreatedAt().Time,
		UpdatedAt: created.GetUpdatedAt().Time,
	}, nil
}

// CreateReplyComment creates an inline comment as a reply to an existing comment.
func CreateReplyComment(ctx context.Context, client *gh.Client, owner, repo string, number int, body string, inReplyTo int64) (*ReviewComment, error) {
	created, _, err := client.PullRequests.CreateCommentInReplyTo(ctx, owner, repo, number, body, inReplyTo)
	if err != nil {
		return nil, fmt.Errorf("create reply comment: %w", err)
	}
	return &ReviewComment{
		ID:   created.GetID(),
		Body: created.GetBody(),
		Path: created.GetPath(),
		Line: created.GetLine(),
		Side: created.GetSide(),
		Author: User{
			Login:     created.GetUser().GetLogin(),
			AvatarURL: created.GetUser().GetAvatarURL(),
		},
		InReplyTo: created.GetInReplyTo(),
		CreatedAt: created.GetCreatedAt().Time,
		UpdatedAt: created.GetUpdatedAt().Time,
	}, nil
}

// AddCommentReaction adds a reaction to a pull request review comment.
func AddCommentReaction(ctx context.Context, client *gh.Client, owner, repo string, commentID int64, content string) error {
	_, _, err := client.Reactions.CreatePullRequestCommentReaction(ctx, owner, repo, commentID, content)
	if err != nil {
		return fmt.Errorf("add comment reaction: %w", err)
	}
	return nil
}

// AddIssueCommentReaction adds a reaction to an issue comment.
func AddIssueCommentReaction(ctx context.Context, client *gh.Client, owner, repo string, commentID int64, content string) error {
	_, _, err := client.Reactions.CreateIssueCommentReaction(ctx, owner, repo, commentID, content)
	return err
}

// UpdateReviewComment updates an existing inline review comment.
func UpdateReviewComment(ctx context.Context, client *gh.Client, owner, repo string, commentID int64, body string) (*ReviewComment, error) {
	comment := &gh.PullRequestComment{
		Body: gh.Ptr(body),
	}
	updated, _, err := client.PullRequests.EditComment(ctx, owner, repo, commentID, comment)
	if err != nil {
		return nil, fmt.Errorf("update review comment: %w", err)
	}
	return &ReviewComment{
		ID:   updated.GetID(),
		Body: updated.GetBody(),
		Path: updated.GetPath(),
		Line: updated.GetLine(),
		Side: updated.GetSide(),
		Author: User{
			Login:     updated.GetUser().GetLogin(),
			AvatarURL: updated.GetUser().GetAvatarURL(),
		},
		CreatedAt: updated.GetCreatedAt().Time,
		UpdatedAt: updated.GetUpdatedAt().Time,
	}, nil
}

// DeleteReviewComment deletes an inline review comment.
func DeleteReviewComment(ctx context.Context, client *gh.Client, owner, repo string, commentID int64) error {
	_, err := client.PullRequests.DeleteComment(ctx, owner, repo, commentID)
	if err != nil {
		return fmt.Errorf("delete review comment: %w", err)
	}
	return nil
}

// FetchIssueComments fetches top-level issue comments on a PR (not inline review comments).
func FetchIssueComments(ctx context.Context, client *gh.Client, owner, repo string, number int) ([]IssueComment, error) {
	opts := &gh.IssueListCommentsOptions{
		ListOptions: gh.ListOptions{PerPage: 100},
	}
	var result []IssueComment
	for {
		comments, resp, err := client.Issues.ListComments(ctx, owner, repo, number, opts)
		if err != nil {
			return nil, fmt.Errorf("fetch issue comments: %w", err)
		}
		for _, c := range comments {
			ic := IssueComment{
				ID:        c.GetID(),
				Body:      c.GetBody(),
				CreatedAt: c.GetCreatedAt().Time,
				Author: User{
					Login:     c.GetUser().GetLogin(),
					AvatarURL: c.GetUser().GetAvatarURL(),
				},
			}
			if c.Reactions != nil {
				ic.Reactions = &ReactionSummary{
					PlusOne:  c.Reactions.GetPlusOne(),
					MinusOne: c.Reactions.GetMinusOne(),
					Laugh:    c.Reactions.GetLaugh(),
					Confused: c.Reactions.GetConfused(),
					Heart:    c.Reactions.GetHeart(),
					Hooray:   c.Reactions.GetHooray(),
					Rocket:   c.Reactions.GetRocket(),
					Eyes:     c.Reactions.GetEyes(),
				}
			}
			result = append(result, ic)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return result, nil
}

// FetchPRCommits fetches all commits on a pull request with pagination.
func FetchPRCommits(ctx context.Context, client *gh.Client, owner, repo string, number int) ([]PRCommit, error) {
	opts := &gh.ListOptions{PerPage: 100}
	var result []PRCommit
	for {
		commits, resp, err := client.PullRequests.ListCommits(ctx, owner, repo, number, opts)
		if err != nil {
			return nil, fmt.Errorf("fetch PR commits: %w", err)
		}
		for _, c := range commits {
			commit := PRCommit{
				SHA:     c.GetSHA(),
				Message: c.GetCommit().GetMessage(),
			}
			if c.GetAuthor() != nil {
				commit.Author = User{
					Login:     c.GetAuthor().GetLogin(),
					AvatarURL: c.GetAuthor().GetAvatarURL(),
				}
			}
			if c.GetCommit().GetAuthor() != nil {
				commit.Date = c.GetCommit().GetAuthor().GetDate().Time
			}
			result = append(result, commit)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return result, nil
}

// FetchCompare fetches the raw diff between two refs via the compare API.
func FetchCompare(ctx context.Context, client *gh.Client, owner, repo, base, head string) (string, error) {
	url := fmt.Sprintf("repos/%s/%s/compare/%s...%s", owner, repo, base, head)
	req, err := client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create compare request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3.diff")

	resp, err := client.BareDo(ctx, req)
	if err != nil {
		return "", fmt.Errorf("fetch compare diff: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read compare body: %w", err)
	}
	return string(body), nil
}

// SubmitReview submits a pull request review.
func SubmitReview(ctx context.Context, client *gh.Client, owner, repo string, number int, event, body string, comments []InlineCommentRequest) (*Review, error) {
	var reviewComments []*gh.DraftReviewComment
	for _, c := range comments {
		reviewComments = append(reviewComments, &gh.DraftReviewComment{
			Path: gh.Ptr(c.Path),
			Line: gh.Ptr(c.Line),
			Side: gh.Ptr(c.Side),
			Body: gh.Ptr(c.Body),
		})
	}
	review := &gh.PullRequestReviewRequest{
		Event:    gh.Ptr(event),
		Body:     gh.Ptr(body),
		Comments: reviewComments,
	}
	created, _, err := client.PullRequests.CreateReview(ctx, owner, repo, number, review)
	if err != nil {
		return nil, fmt.Errorf("submit review: %w", err)
	}
	return &Review{
		ID:        created.GetID(),
		State:     created.GetState(),
		Body:      created.GetBody(),
		CreatedAt: created.GetSubmittedAt().Time,
		Author: User{
			Login:     created.GetUser().GetLogin(),
			AvatarURL: created.GetUser().GetAvatarURL(),
		},
	}, nil
}
