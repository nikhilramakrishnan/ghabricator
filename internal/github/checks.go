package github

import (
	"context"
	"fmt"
	"time"

	gh "github.com/google/go-github/v68/github"
)

// CheckRun represents a single CI check run (GitHub Actions, etc.).
type CheckRun struct {
	Name        string
	Status      string // queued, in_progress, completed
	Conclusion  string // success, failure, neutral, cancelled, skipped, timed_out, action_required
	DetailsURL  string
	AppName     string
	StartedAt   time.Time
	CompletedAt time.Time
}

// FetchCheckRuns returns check runs for a given commit SHA.
func FetchCheckRuns(ctx context.Context, client *gh.Client, owner, repo, ref string) ([]CheckRun, error) {
	opts := &gh.ListCheckRunsOptions{
		ListOptions: gh.ListOptions{PerPage: 100},
	}
	result, _, err := client.Checks.ListCheckRunsForRef(ctx, owner, repo, ref, opts)
	if err != nil {
		return nil, fmt.Errorf("fetch check runs: %w", err)
	}
	runs := make([]CheckRun, 0, len(result.CheckRuns))
	for _, cr := range result.CheckRuns {
		run := CheckRun{
			Name:       cr.GetName(),
			Status:     cr.GetStatus(),
			Conclusion: cr.GetConclusion(),
			DetailsURL: cr.GetDetailsURL(),
		}
		if cr.App != nil {
			run.AppName = cr.App.GetName()
		}
		if cr.StartedAt != nil {
			run.StartedAt = cr.StartedAt.Time
		}
		if cr.CompletedAt != nil {
			run.CompletedAt = cr.CompletedAt.Time
		}
		runs = append(runs, run)
	}
	return runs, nil
}
