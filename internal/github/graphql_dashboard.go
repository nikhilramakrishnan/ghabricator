package github

import (
	"context"
	"strings"
	"time"
)

const dashboardQuery = `
query($authoredQuery: String!, $reviewQuery: String!) {
  authored: search(query: $authoredQuery, type: ISSUE, first: 25) {
    nodes {
      ... on PullRequest {
        number
        title
        isDraft
        updatedAt
        repository { nameWithOwner }
        author { login avatarUrl }
        labels(first: 10) { nodes { name color } }
        assignees(first: 10) { nodes { login avatarUrl } }
      }
    }
  }
  reviewing: search(query: $reviewQuery, type: ISSUE, first: 25) {
    nodes {
      ... on PullRequest {
        number
        title
        isDraft
        updatedAt
        repository { nameWithOwner }
        author { login avatarUrl }
        labels(first: 10) { nodes { name color } }
        assignees(first: 10) { nodes { login avatarUrl } }
      }
    }
  }
}
`

type dashboardGQLResponse struct {
	Data struct {
		Authored  dashboardSearchResult `json:"authored"`
		Reviewing dashboardSearchResult `json:"reviewing"`
	} `json:"data"`
}

type dashboardSearchResult struct {
	Nodes []dashboardPRNode `json:"nodes"`
}

type dashboardPRNode struct {
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	IsDraft    bool      `json:"isDraft"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Repository *struct {
		NameWithOwner string `json:"nameWithOwner"`
	} `json:"repository"`
	Author *struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatarUrl"`
	} `json:"author"`
	Labels *struct {
		Nodes []struct {
			Name  string `json:"name"`
			Color string `json:"color"`
		} `json:"nodes"`
	} `json:"labels"`
	Assignees *struct {
		Nodes []struct {
			Login     string `json:"login"`
			AvatarURL string `json:"avatarUrl"`
		} `json:"nodes"`
	} `json:"assignees"`
}

// DashboardPR is the GraphQL-sourced equivalent of the REST dashboardPR.
type DashboardPR struct {
	Number    int
	Title     string
	Repo      string // "owner/repo"
	Author    string
	AvatarURL string
	UpdatedAt time.Time
	Draft     bool
	Labels    []Label
	Assignees []User
}

// FetchDashboardGraphQL fetches authored and review-requested PRs in a single
// GraphQL call using search aliases.
func FetchDashboardGraphQL(ctx context.Context, token, login string) (authored, reviewing []DashboardPR, err error) {
	vars := map[string]interface{}{
		"authoredQuery": "is:open is:pr author:" + login,
		"reviewQuery":   "is:open is:pr review-requested:" + login,
	}

	var resp dashboardGQLResponse
	if err := QueryGraphQL(ctx, token, dashboardQuery, vars, &resp); err != nil {
		return nil, nil, err
	}

	authored = convertNodes(resp.Data.Authored.Nodes)
	reviewing = convertNodes(resp.Data.Reviewing.Nodes)
	return authored, reviewing, nil
}

func convertNodes(nodes []dashboardPRNode) []DashboardPR {
	prs := make([]DashboardPR, 0, len(nodes))
	for _, n := range nodes {
		if n.Number == 0 {
			continue // skip non-PR search results
		}
		pr := DashboardPR{
			Number:    n.Number,
			Title:     n.Title,
			Draft:     n.IsDraft,
			UpdatedAt: n.UpdatedAt,
		}
		if n.Repository != nil {
			pr.Repo = n.Repository.NameWithOwner
		}
		if n.Author != nil {
			pr.Author = n.Author.Login
			pr.AvatarURL = n.Author.AvatarURL
		}
		if n.Labels != nil {
			for _, l := range n.Labels.Nodes {
				pr.Labels = append(pr.Labels, Label{
					Name:  l.Name,
					Color: strings.TrimPrefix(l.Color, "#"),
				})
			}
		}
		if n.Assignees != nil {
			for _, a := range n.Assignees.Nodes {
				pr.Assignees = append(pr.Assignees, User{
					Login:     a.Login,
					AvatarURL: a.AvatarURL,
				})
			}
		}
		prs = append(prs, pr)
	}
	return prs
}
