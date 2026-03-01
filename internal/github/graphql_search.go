package github

import (
	"context"
	"strings"
	"time"
)

const searchQuery = `
query($prQuery: String!, $issueQuery: String!, $repoQuery: String!,
      $prFirst: Int!, $issueFirst: Int!, $repoFirst: Int!) {
  prs: search(query: $prQuery, type: ISSUE, first: $prFirst) {
    issueCount
    nodes {
      ... on PullRequest {
        number
        title
        state
        body
        isDraft
        repository { nameWithOwner }
        author { login avatarUrl }
        labels(first: 10) { nodes { name color } }
        createdAt
        updatedAt
        comments { totalCount }
      }
    }
  }
  issues: search(query: $issueQuery, type: ISSUE, first: $issueFirst) {
    issueCount
    nodes {
      ... on Issue {
        number
        title
        state
        body
        repository { nameWithOwner }
        author { login avatarUrl }
        labels(first: 10) { nodes { name color } }
        createdAt
        updatedAt
        comments { totalCount }
      }
    }
  }
  repos: search(query: $repoQuery, type: REPOSITORY, first: $repoFirst) {
    repositoryCount
    nodes {
      ... on Repository {
        nameWithOwner
        description
        stargazerCount
        forkCount
        primaryLanguage { name color }
        updatedAt
        repositoryTopics(first: 10) { nodes { topic { name } } }
        owner { avatarUrl }
      }
    }
  }
}
`

// GraphQL response types for search

type searchGQLResponse struct {
	Data struct {
		PRs    searchPRResult   `json:"prs"`
		Issues searchIssueResult `json:"issues"`
		Repos  searchRepoResult  `json:"repos"`
	} `json:"data"`
}

type searchPRResult struct {
	IssueCount int             `json:"issueCount"`
	Nodes      []searchPRNode  `json:"nodes"`
}

type searchIssueResult struct {
	IssueCount int               `json:"issueCount"`
	Nodes      []searchIssueNode `json:"nodes"`
}

type searchRepoResult struct {
	RepositoryCount int              `json:"repositoryCount"`
	Nodes           []searchRepoNode `json:"nodes"`
}

type searchPRNode struct {
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	State      string    `json:"state"`
	Body       string    `json:"body"`
	IsDraft    bool      `json:"isDraft"`
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
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Comments  *struct {
		TotalCount int `json:"totalCount"`
	} `json:"comments"`
}

type searchIssueNode struct {
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	State      string    `json:"state"`
	Body       string    `json:"body"`
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
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Comments  *struct {
		TotalCount int `json:"totalCount"`
	} `json:"comments"`
}

type searchRepoNode struct {
	NameWithOwner string `json:"nameWithOwner"`
	Description   string `json:"description"`
	StargazerCount int   `json:"stargazerCount"`
	ForkCount     int    `json:"forkCount"`
	PrimaryLanguage *struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"primaryLanguage"`
	UpdatedAt        time.Time `json:"updatedAt"`
	RepositoryTopics *struct {
		Nodes []struct {
			Topic struct {
				Name string `json:"name"`
			} `json:"topic"`
		} `json:"nodes"`
	} `json:"repositoryTopics"`
	Owner *struct {
		AvatarURL string `json:"avatarUrl"`
	} `json:"owner"`
}

// SearchResult holds counts for all search types plus full results for the selected type.
type SearchResult struct {
	Counts map[string]int
	PRs    []SearchPR
	Issues []SearchIssue
	Repos  []SearchRepo
}

type SearchPR struct {
	Number        int
	Title         string
	State         string
	Body          string
	Draft         bool
	Repo          string
	Author        string
	AvatarURL     string
	Labels        []Label
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CommentsCount int
}

type SearchIssue struct {
	Number        int
	Title         string
	State         string
	Body          string
	Repo          string
	Author        string
	AvatarURL     string
	Labels        []Label
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CommentsCount int
}

type SearchRepo struct {
	FullName    string
	Description string
	Stars       int
	Forks       int
	Language    string
	LangColor   string
	AvatarURL   string
	UpdatedAt   time.Time
	Topics      []string
}

// SearchGraphQL executes a single GraphQL query to get counts for PRs, issues,
// and repos, plus full results for the selected type.
func SearchGraphQL(ctx context.Context, token, query, selectedType string) (*SearchResult, error) {
	prFirst, issueFirst, repoFirst := 0, 0, 0
	switch selectedType {
	case "prs":
		prFirst = 25
	case "issues":
		issueFirst = 25
	case "repos":
		repoFirst = 25
	default:
		prFirst = 25
	}

	vars := map[string]interface{}{
		"prQuery":    "is:pr " + query,
		"issueQuery": "is:issue " + query,
		"repoQuery":  query,
		"prFirst":    prFirst,
		"issueFirst": issueFirst,
		"repoFirst":  repoFirst,
	}

	var resp searchGQLResponse
	if err := QueryGraphQL(ctx, token, searchQuery, vars, &resp); err != nil {
		return nil, err
	}

	result := &SearchResult{
		Counts: map[string]int{
			"prs":    resp.Data.PRs.IssueCount,
			"issues": resp.Data.Issues.IssueCount,
			"repos":  resp.Data.Repos.RepositoryCount,
		},
	}

	// Convert PR nodes
	for _, n := range resp.Data.PRs.Nodes {
		if n.Number == 0 {
			continue
		}
		pr := SearchPR{
			Number:    n.Number,
			Title:     n.Title,
			State:     mapSearchState(n.State),
			Body:      truncateBody(n.Body, 200),
			Draft:     n.IsDraft,
			CreatedAt: n.CreatedAt,
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
		if n.Comments != nil {
			pr.CommentsCount = n.Comments.TotalCount
		}
		result.PRs = append(result.PRs, pr)
	}

	// Convert issue nodes
	for _, n := range resp.Data.Issues.Nodes {
		if n.Number == 0 {
			continue
		}
		issue := SearchIssue{
			Number:    n.Number,
			Title:     n.Title,
			State:     mapSearchState(n.State),
			Body:      truncateBody(n.Body, 200),
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
		}
		if n.Repository != nil {
			issue.Repo = n.Repository.NameWithOwner
		}
		if n.Author != nil {
			issue.Author = n.Author.Login
			issue.AvatarURL = n.Author.AvatarURL
		}
		if n.Labels != nil {
			for _, l := range n.Labels.Nodes {
				issue.Labels = append(issue.Labels, Label{
					Name:  l.Name,
					Color: strings.TrimPrefix(l.Color, "#"),
				})
			}
		}
		if n.Comments != nil {
			issue.CommentsCount = n.Comments.TotalCount
		}
		result.Issues = append(result.Issues, issue)
	}

	// Convert repo nodes
	for _, n := range resp.Data.Repos.Nodes {
		if n.NameWithOwner == "" {
			continue
		}
		repo := SearchRepo{
			FullName:    n.NameWithOwner,
			Description: n.Description,
			Stars:       n.StargazerCount,
			Forks:       n.ForkCount,
			UpdatedAt:   n.UpdatedAt,
		}
		if n.PrimaryLanguage != nil {
			repo.Language = n.PrimaryLanguage.Name
			repo.LangColor = n.PrimaryLanguage.Color
		}
		if n.Owner != nil {
			repo.AvatarURL = n.Owner.AvatarURL
		}
		if n.RepositoryTopics != nil {
			for _, t := range n.RepositoryTopics.Nodes {
				repo.Topics = append(repo.Topics, t.Topic.Name)
			}
		}
		result.Repos = append(result.Repos, repo)
	}

	return result, nil
}

func mapSearchState(state string) string {
	switch state {
	case "OPEN":
		return "open"
	case "CLOSED":
		return "closed"
	case "MERGED":
		return "merged"
	default:
		return strings.ToLower(state)
	}
}

func truncateBody(body string, maxLen int) string {
	// Strip newlines for excerpt
	body = strings.ReplaceAll(body, "\r\n", " ")
	body = strings.ReplaceAll(body, "\n", " ")
	body = strings.TrimSpace(body)
	if len(body) <= maxLen {
		return body
	}
	return body[:maxLen] + "..."
}
