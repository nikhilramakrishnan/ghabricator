package github

import (
	"context"
	"fmt"
	"time"
)

const prDetailQuery = `
query PRDetail($owner: String!, $repo: String!, $number: Int!) {
  repository(owner: $owner, name: $repo) {
    viewerPermission
    pullRequest(number: $number) {
      number
      title
      body
      state
      isDraft
      merged
      createdAt
      updatedAt
      additions
      deletions
      changedFiles
      author {
        login
        avatarUrl
      }
      headRef {
        name
        target { oid }
        repository { nameWithOwner }
      }
      baseRef {
        name
        target { oid }
        repository { nameWithOwner }
      }
      labels(first: 50) {
        nodes { name color }
      }
      reviewRequests(first: 50) {
        nodes {
          requestedReviewer {
            ... on User { login avatarUrl }
          }
        }
      }
      reviews(first: 100) {
        nodes {
          databaseId
          state
          body
          createdAt
          author { login avatarUrl }
        }
      }
      comments(first: 100) {
        nodes {
          databaseId
          body
          createdAt
          author { login avatarUrl }
          reactionGroups {
            content
            reactors(first: 0) { totalCount }
          }
        }
      }
      commits(first: 250) {
        nodes {
          commit {
            oid
            message
            author {
              user { login avatarUrl }
              date
            }
          }
        }
      }
    }
  }
}
`

// graphQL response types — intermediate structs for unmarshalling

type gqlAuthor struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatarUrl"`
}

type gqlReactionGroup struct {
	Content  string `json:"content"`
	Reactors struct {
		TotalCount int `json:"totalCount"`
	} `json:"reactors"`
}

type gqlPRResponse struct {
	Data struct {
		Repository struct {
			ViewerPermission string `json:"viewerPermission"`
			PullRequest      struct {
				Number       int       `json:"number"`
				Title        string    `json:"title"`
				Body         string    `json:"body"`
				State        string    `json:"state"`
				IsDraft      bool      `json:"isDraft"`
				Merged       bool      `json:"merged"`
				CreatedAt    time.Time `json:"createdAt"`
				UpdatedAt    time.Time `json:"updatedAt"`
				Additions    int       `json:"additions"`
				Deletions    int       `json:"deletions"`
				ChangedFiles int       `json:"changedFiles"`
				Author       gqlAuthor `json:"author"`
				HeadRef      *struct {
					Name       string `json:"name"`
					Target     struct{ Oid string } `json:"target"`
					Repository struct{ NameWithOwner string } `json:"repository"`
				} `json:"headRef"`
				BaseRef *struct {
					Name       string `json:"name"`
					Target     struct{ Oid string } `json:"target"`
					Repository struct{ NameWithOwner string } `json:"repository"`
				} `json:"baseRef"`
				Labels struct {
					Nodes []struct {
						Name  string `json:"name"`
						Color string `json:"color"`
					} `json:"nodes"`
				} `json:"labels"`
				ReviewRequests struct {
					Nodes []struct {
						RequestedReviewer *gqlAuthor `json:"requestedReviewer"`
					} `json:"nodes"`
				} `json:"reviewRequests"`
				Reviews struct {
					Nodes []struct {
						DatabaseId int64     `json:"databaseId"`
						State      string    `json:"state"`
						Body       string    `json:"body"`
						CreatedAt  time.Time `json:"createdAt"`
						Author     gqlAuthor `json:"author"`
					} `json:"nodes"`
				} `json:"reviews"`
				Comments struct {
					Nodes []struct {
						DatabaseId     int64              `json:"databaseId"`
						Body           string             `json:"body"`
						CreatedAt      time.Time          `json:"createdAt"`
						Author         gqlAuthor          `json:"author"`
						ReactionGroups []gqlReactionGroup `json:"reactionGroups"`
					} `json:"nodes"`
				} `json:"comments"`
				Commits struct {
					Nodes []struct {
						Commit struct {
							Oid     string `json:"oid"`
							Message string `json:"message"`
							Author  struct {
								User *gqlAuthor `json:"user"`
								Date time.Time  `json:"date"`
							} `json:"author"`
						} `json:"commit"`
					} `json:"nodes"`
				} `json:"commits"`
			} `json:"pullRequest"`
		} `json:"repository"`
	} `json:"data"`
}

// PRDetailGraphQL holds all data fetched from the single GraphQL query.
type PRDetailGraphQL struct {
	PR               *PullRequest
	Reviews          []Review
	IssueComments    []IssueComment
	Commits          []PRCommit
	CheckRuns        []CheckRun
	ViewerPermission string // ADMIN, MAINTAIN, WRITE, TRIAGE, READ, or ""
}

// FetchPRDetailGraphQL fetches PR metadata, reviews, issue comments, commits,
// and check runs in a single GraphQL query.
func FetchPRDetailGraphQL(ctx context.Context, token, owner, repo string, number int) (*PRDetailGraphQL, error) {
	vars := map[string]interface{}{
		"owner":  owner,
		"repo":   repo,
		"number": number,
	}

	var resp gqlPRResponse
	if err := QueryGraphQL(ctx, token, prDetailQuery, vars, &resp); err != nil {
		return nil, fmt.Errorf("graphql PR detail: %w", err)
	}

	gpr := resp.Data.Repository.PullRequest
	if gpr.Number == 0 {
		return nil, fmt.Errorf("pull request %s/%s#%d not found", owner, repo, number)
	}

	// Map PR metadata
	pr := &PullRequest{
		Number:       gpr.Number,
		Title:        gpr.Title,
		Body:         gpr.Body,
		State:        mapPRState(gpr.State),
		Draft:        gpr.IsDraft,
		Merged:       gpr.Merged,
		CreatedAt:    gpr.CreatedAt,
		UpdatedAt:    gpr.UpdatedAt,
		Additions:    gpr.Additions,
		Deletions:    gpr.Deletions,
		ChangedFiles: gpr.ChangedFiles,
		Author:       User{Login: gpr.Author.Login, AvatarURL: gpr.Author.AvatarUrl},
	}

	if gpr.HeadRef != nil {
		pr.Head = Ref{
			Ref:  gpr.HeadRef.Name,
			SHA:  gpr.HeadRef.Target.Oid,
			Repo: gpr.HeadRef.Repository.NameWithOwner,
		}
	}
	if gpr.BaseRef != nil {
		pr.Base = Ref{
			Ref:  gpr.BaseRef.Name,
			SHA:  gpr.BaseRef.Target.Oid,
			Repo: gpr.BaseRef.Repository.NameWithOwner,
		}
	}

	for _, l := range gpr.Labels.Nodes {
		pr.Labels = append(pr.Labels, Label{Name: l.Name, Color: l.Color})
	}

	for _, rr := range gpr.ReviewRequests.Nodes {
		if rr.RequestedReviewer != nil {
			pr.Reviewers = append(pr.Reviewers, User{
				Login:     rr.RequestedReviewer.Login,
				AvatarURL: rr.RequestedReviewer.AvatarUrl,
			})
		}
	}

	// Map reviews
	reviews := make([]Review, 0, len(gpr.Reviews.Nodes))
	for _, r := range gpr.Reviews.Nodes {
		reviews = append(reviews, Review{
			ID:        r.DatabaseId,
			State:     r.State,
			Body:      r.Body,
			CreatedAt: r.CreatedAt,
			Author:    User{Login: r.Author.Login, AvatarURL: r.Author.AvatarUrl},
		})
	}

	// Map issue comments
	issueComments := make([]IssueComment, 0, len(gpr.Comments.Nodes))
	for _, c := range gpr.Comments.Nodes {
		ic := IssueComment{
			ID:        c.DatabaseId,
			Body:      c.Body,
			CreatedAt: c.CreatedAt,
			Author:    User{Login: c.Author.Login, AvatarURL: c.Author.AvatarUrl},
		}
		if rs := mapReactionGroups(c.ReactionGroups); rs != nil {
			ic.Reactions = rs
		}
		issueComments = append(issueComments, ic)
	}

	// Map commits
	commits := make([]PRCommit, 0, len(gpr.Commits.Nodes))
	for _, n := range gpr.Commits.Nodes {
		c := n.Commit
		commit := PRCommit{
			SHA:     c.Oid,
			Message: c.Message,
			Date:    c.Author.Date,
		}
		if c.Author.User != nil {
			commit.Author = User{Login: c.Author.User.Login, AvatarURL: c.Author.User.AvatarUrl}
		}
		commits = append(commits, commit)
	}

	return &PRDetailGraphQL{
		PR:               pr,
		Reviews:          reviews,
		IssueComments:    issueComments,
		Commits:          commits,
		ViewerPermission: resp.Data.Repository.ViewerPermission,
	}, nil
}

// mapPRState converts GraphQL UPPER_CASE state to REST lower_case.
func mapPRState(state string) string {
	switch state {
	case "OPEN":
		return "open"
	case "CLOSED":
		return "closed"
	case "MERGED":
		return "closed"
	default:
		return state
	}
}

// mapReactionGroups converts GraphQL reactionGroups to a ReactionSummary.
func mapReactionGroups(groups []gqlReactionGroup) *ReactionSummary {
	if len(groups) == 0 {
		return nil
	}
	rs := &ReactionSummary{}
	any := false
	for _, g := range groups {
		if g.Reactors.TotalCount == 0 {
			continue
		}
		any = true
		switch g.Content {
		case "THUMBS_UP":
			rs.PlusOne = g.Reactors.TotalCount
		case "THUMBS_DOWN":
			rs.MinusOne = g.Reactors.TotalCount
		case "LAUGH":
			rs.Laugh = g.Reactors.TotalCount
		case "CONFUSED":
			rs.Confused = g.Reactors.TotalCount
		case "HEART":
			rs.Heart = g.Reactors.TotalCount
		case "HOORAY":
			rs.Hooray = g.Reactors.TotalCount
		case "ROCKET":
			rs.Rocket = g.Reactors.TotalCount
		case "EYES":
			rs.Eyes = g.Reactors.TotalCount
		}
	}
	if !any {
		return nil
	}
	return rs
}
