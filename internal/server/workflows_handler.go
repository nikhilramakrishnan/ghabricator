package server

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"

	"github.com/nikhilr/ghabricator/internal/auth"
	ghapi "github.com/nikhilr/ghabricator/internal/github"

	gh "github.com/google/go-github/v68/github"
)

func (s *Server) handleAPIWorkflowRuns(w http.ResponseWriter, r *http.Request) {
	client := auth.GitHubClientFromContext(r.Context())
	sess := auth.SessionFromContext(r.Context())
	login := sess.Login

	// Discover repos from user's open PRs (same approach as dashboard)
	repos := discoverReposFromPRs(r, client, login)
	if len(repos) == 0 {
		jsonOK(w, APIWorkflowRunsResponse{Runs: []APIWorkflowRun{}})
		return
	}

	// Fan out fetches, bounded at 5 goroutines
	var keys []repoKey
	for k := range repos {
		keys = append(keys, k)
	}

	var (
		mu      sync.Mutex
		allRuns []ghapi.WorkflowRun
		wg      sync.WaitGroup
		sem     = make(chan struct{}, 5)
	)

	for _, k := range keys {
		wg.Add(1)
		sem <- struct{}{}
		go func(owner, repo string) {
			defer wg.Done()
			defer func() { <-sem }()
			runs, err := ghapi.FetchWorkflowRuns(r.Context(), client, owner, repo, 10)
			if err != nil {
				log.Printf("workflow runs for %s/%s: %v", owner, repo, err)
				return
			}
			mu.Lock()
			allRuns = append(allRuns, runs...)
			mu.Unlock()
		}(k.owner, k.repo)
	}
	wg.Wait()

	// Sort by CreatedAt desc
	sort.Slice(allRuns, func(i, j int) bool {
		return allRuns[i].CreatedAt.After(allRuns[j].CreatedAt)
	})

	// Cap at 50
	if len(allRuns) > 50 {
		allRuns = allRuns[:50]
	}

	apiRuns := make([]APIWorkflowRun, 0, len(allRuns))
	for _, r := range allRuns {
		apiRuns = append(apiRuns, APIWorkflowRun{
			ID:           r.ID,
			Name:         r.Name,
			DisplayTitle: r.DisplayTitle,
			Status:       r.Status,
			Conclusion:   r.Conclusion,
			Branch:       r.Branch,
			Event:        r.Event,
			Actor: APIUser{
				Login:     r.ActorLogin,
				AvatarURL: r.ActorAvatarURL,
			},
			RepoOwner:  r.RepoOwner,
			RepoName:   r.RepoName,
			DurationMs: r.Duration.Milliseconds(),
			HTMLURL:    r.HTMLURL,
			CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	jsonOK(w, APIWorkflowRunsResponse{Runs: apiRuns})
}

type repoKey struct{ owner, repo string }

// discoverReposFromPRs finds unique repos from the user's open PRs.
func discoverReposFromPRs(r *http.Request, client *gh.Client, login string) map[repoKey]bool {
	repos := make(map[repoKey]bool)

	queries := []string{
		fmt.Sprintf("is:open is:pr author:%s", login),
		fmt.Sprintf("is:open is:pr review-requested:%s", login),
	}

	for _, q := range queries {
		result, _, err := client.Search.Issues(r.Context(), q, &gh.SearchOptions{
			Sort:        "updated",
			Order:       "desc",
			ListOptions: gh.ListOptions{PerPage: 25},
		})
		if err != nil {
			log.Printf("discover repos search error: %v", err)
			continue
		}
		for _, issue := range result.Issues {
			fullName := ""
			if issue.Repository != nil {
				fullName = issue.Repository.GetFullName()
			} else if issue.RepositoryURL != nil {
				fullName = extractRepoFromURL(issue.GetRepositoryURL())
			}
			if fullName == "" {
				continue
			}
			owner, repo := splitRepo(fullName)
			if owner != "" && repo != "" {
				repos[repoKey{owner, repo}] = true
			}
		}
	}
	return repos
}
