package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nikhilr/ghabricator/internal/auth"

	gh "github.com/google/go-github/v68/github"
)

type paletteResult struct {
	Icon string `json:"icon"`
	Text string `json:"text"`
	Hint string `json:"hint"`
	Href string `json:"href"`
}

func (s *Server) handlePalette(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	client := auth.GitHubClientFromContext(r.Context())
	q := r.URL.Query().Get("q")
	if q == "" {
		json.NewEncoder(w).Encode([]paletteResult{})
		return
	}

	var results []paletteResult

	// Static navigation matches
	navItems := []struct {
		name, icon, hint, href string
		keywords               []string
	}{
		{"Dashboard", "fa-th-list", "Page", "/dashboard", []string{"dashboard", "home", "revisions"}},
		{"Paste", "fa-clipboard", "Page", "/paste", []string{"paste", "gist", "snippet"}},
		{"Create Paste", "fa-plus", "Action", "/paste/new", []string{"paste", "new", "create", "gist"}},
		{"Herald", "fa-shield", "Page", "/herald", []string{"herald", "rules", "automation"}},
		{"Create Herald Rule", "fa-plus", "Action", "/herald/new", []string{"herald", "new", "create", "rule"}},
		{"Search", "fa-search", "Page", "/search", []string{"search", "find"}},
	}

	lq := toLowerCase(q)
	for _, nav := range navItems {
		if fuzzyMatch(lq, nav.name, nav.keywords) {
			results = append(results, paletteResult{
				Icon: nav.icon,
				Text: nav.name,
				Hint: nav.hint,
				Href: nav.href,
			})
		}
	}

	// User's own repos (prioritized over global search)
	userRepos, _, err := client.Repositories.List(r.Context(), sess.Login, &gh.RepositoryListOptions{
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: gh.ListOptions{PerPage: 100},
	})
	if err != nil {
		log.Printf("palette user repos error: %v", err)
	} else {
		count := 0
		for _, repo := range userRepos {
			if count >= 5 {
				break
			}
			name := repo.GetFullName()
			if fuzzyMatch(lq, name, nil) {
				results = append(results, paletteResult{
					Icon: "fa-database",
					Text: name,
					Hint: "Your repo",
					Href: "/repo/" + name,
				})
				count++
			}
		}
	}

	// Repos the user contributes to / has starred
	starredRepos, _, err := client.Activity.ListStarred(r.Context(), sess.Login, &gh.ActivityListStarredOptions{
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: gh.ListOptions{PerPage: 50},
	})
	if err != nil {
		log.Printf("palette starred repos error: %v", err)
	} else {
		count := 0
		for _, sr := range starredRepos {
			if count >= 3 {
				break
			}
			repo := sr.GetRepository()
			if repo == nil {
				continue
			}
			name := repo.GetFullName()
			if fuzzyMatch(lq, name, nil) {
				// Avoid duplicates from own repos
				dup := false
				for _, r := range results {
					if r.Href == "/repo/"+name {
						dup = true
						break
					}
				}
				if !dup {
					results = append(results, paletteResult{
						Icon: "fa-star",
						Text: name,
						Hint: "Starred",
						Href: "/repo/" + name,
					})
					count++
				}
			}
		}
	}

	// Search GitHub for PRs (quick, 5 results)
	prs, _, err := client.Search.Issues(r.Context(), "is:pr "+q, &gh.SearchOptions{
		Sort:        "updated",
		Order:       "desc",
		ListOptions: gh.ListOptions{PerPage: 5},
	})
	if err != nil {
		log.Printf("palette PR search error: %v", err)
	} else {
		for _, issue := range prs.Issues {
			if issue.PullRequestLinks == nil {
				continue
			}
			repo := ""
			if issue.Repository != nil {
				repo = issue.Repository.GetFullName()
			} else if issue.RepositoryURL != nil {
				repo = extractRepoFromURL(issue.GetRepositoryURL())
			}
			href := "#"
			if repo != "" {
				href = "/pr/" + repo + "/" + itoa(issue.GetNumber())
			}
			results = append(results, paletteResult{
				Icon: "fa-code-fork",
				Text: issue.GetTitle(),
				Hint: repo + "#" + itoa(issue.GetNumber()),
				Href: href,
			})
		}
	}

	// Global repo search (only if no user repos matched)
	hasRepoResults := false
	for _, r := range results {
		if r.Hint == "Your repo" || r.Hint == "Starred" {
			hasRepoResults = true
			break
		}
	}
	if !hasRepoResults {
		repos, _, err := client.Search.Repositories(r.Context(), q, &gh.SearchOptions{
			Sort:        "stars",
			Order:       "desc",
			ListOptions: gh.ListOptions{PerPage: 3},
		})
		if err != nil {
			log.Printf("palette repo search error: %v", err)
		} else {
			for _, repo := range repos.Repositories {
				results = append(results, paletteResult{
					Icon: "fa-database",
					Text: repo.GetFullName(),
					Hint: "Repository",
					Href: "/repo/" + repo.GetFullName(),
				})
			}
		}
	}

	// Always add a "Search for X" fallback
	results = append(results, paletteResult{
		Icon: "fa-search",
		Text: "Search for \"" + q + "\"",
		Hint: "Full search",
		Href: "/search?q=" + q,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func toLowerCase(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}

func fuzzyMatch(query, name string, keywords []string) bool {
	ln := toLowerCase(name)
	if contains(ln, query) {
		return true
	}
	for _, kw := range keywords {
		if contains(kw, query) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(substr) <= len(s) && (s == substr || len(substr) == 0 || indexString(s, substr) >= 0)
}

func indexString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
