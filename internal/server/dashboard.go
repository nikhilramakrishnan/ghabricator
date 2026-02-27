package server

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/nikhilr/ghabricator/internal/auth"

	gh "github.com/google/go-github/v68/github"
)

type dashboardPR struct {
	Title     string
	Repo      string
	Number    int
	URL       string
	Author    string
	AvatarURL string
	UpdatedAt string
	Draft     bool
	Labels    []string
	Reviewers []dashboardUser
}

type dashboardUser struct {
	Login     string
	AvatarURL string
}

func (s *Server) searchPRs(r *http.Request, client *gh.Client, query string) []dashboardPR {
	result, _, err := client.Search.Issues(r.Context(), query, &gh.SearchOptions{
		Sort:  "updated",
		Order: "desc",
		ListOptions: gh.ListOptions{PerPage: 25},
	})
	if err != nil {
		log.Printf("dashboard search error: %v", err)
		return nil
	}

	var prs []dashboardPR
	for _, issue := range result.Issues {
		if issue.PullRequestLinks == nil {
			continue
		}
		repo := ""
		if issue.Repository != nil {
			repo = issue.Repository.GetFullName()
		} else if issue.RepositoryURL != nil {
			// Parse owner/repo from URL like https://api.github.com/repos/owner/repo
			repo = extractRepoFromURL(issue.GetRepositoryURL())
		}

		pr := dashboardPR{
			Title:     issue.GetTitle(),
			Repo:      repo,
			Number:    issue.GetNumber(),
			Author:    issue.GetUser().GetLogin(),
			AvatarURL: issue.GetUser().GetAvatarURL(),
			UpdatedAt: timeAgo(issue.GetUpdatedAt().Time),
			Draft:     issue.IsPullRequest() && issue.GetDraft(),
		}

		// Build link
		if repo != "" {
			pr.URL = fmt.Sprintf("/pr/%s/%d", repo, issue.GetNumber())
		}

		for _, l := range issue.Labels {
			pr.Labels = append(pr.Labels, l.GetName())
		}
		for _, a := range issue.Assignees {
			pr.Reviewers = append(pr.Reviewers, dashboardUser{
				Login:     a.GetLogin(),
				AvatarURL: a.GetAvatarURL(),
			})
		}

		prs = append(prs, pr)
	}
	return prs
}

func (s *Server) handleAPIDashboard(w http.ResponseWriter, r *http.Request) {
	client := auth.GitHubClientFromContext(r.Context())
	sess := auth.SessionFromContext(r.Context())
	login := sess.Login

	authored := s.searchPRs(r, client, fmt.Sprintf("is:open is:pr author:%s", login))
	reviewing := s.searchPRs(r, client, fmt.Sprintf("is:open is:pr review-requested:%s", login))

	resp := APIDashboardResponse{
		Authored:        dashboardPRsToAPI(authored),
		ReviewRequested: dashboardPRsToAPI(reviewing),
	}
	jsonOK(w, resp)
}

func dashboardPRsToAPI(prs []dashboardPR) []APIPRSummary {
	result := make([]APIPRSummary, 0, len(prs))
	for _, pr := range prs {
		owner, repo := splitRepo(pr.Repo)
		s := APIPRSummary{
			Number: pr.Number,
			Title:  pr.Title,
			Owner:  owner,
			Repo:   repo,
			Author: APIUser{
				Login:     pr.Author,
				AvatarURL: pr.AvatarURL,
			},
			Draft:     pr.Draft,
			UpdatedAt: pr.UpdatedAt,
		}
		for _, l := range pr.Labels {
			s.Labels = append(s.Labels, APILabel{Name: l})
		}
		for _, u := range pr.Reviewers {
			s.Reviewers = append(s.Reviewers, APIUser{
				Login:     u.Login,
				AvatarURL: u.AvatarURL,
			})
		}
		result = append(result, s)
	}
	return result
}

func splitRepo(fullName string) (owner, repo string) {
	for i := 0; i < len(fullName); i++ {
		if fullName[i] == '/' {
			return fullName[:i], fullName[i+1:]
		}
	}
	return fullName, ""
}

func extractRepoFromURL(url string) string {
	// https://api.github.com/repos/owner/repo → owner/repo
	const prefix = "/repos/"
	idx := 0
	for i := 0; i <= len(url)-len(prefix); i++ {
		if url[i:i+len(prefix)] == prefix {
			idx = i + len(prefix)
			break
		}
	}
	if idx > 0 {
		return url[idx:]
	}
	return ""
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(math.Round(d.Minutes()))
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(math.Round(d.Hours()))
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	default:
		days := int(math.Round(d.Hours() / 24))
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
