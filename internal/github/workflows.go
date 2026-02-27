package github

import (
	"context"
	"fmt"
	"time"

	gh "github.com/google/go-github/v68/github"
)

type WorkflowRun struct {
	ID              int64
	Name            string
	DisplayTitle    string
	Status          string // queued, in_progress, completed
	Conclusion      string // success, failure, cancelled, skipped, etc.
	Branch          string
	Event           string // push, pull_request, schedule, etc.
	ActorLogin      string
	ActorAvatarURL  string
	RepoOwner       string
	RepoName        string
	Duration        time.Duration
	HTMLURL         string
	CreatedAt       time.Time
}

func FetchWorkflowRuns(ctx context.Context, client *gh.Client, owner, repo string, perPage int) ([]WorkflowRun, error) {
	opts := &gh.ListWorkflowRunsOptions{
		ListOptions: gh.ListOptions{PerPage: perPage},
	}
	result, _, err := client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("list workflow runs for %s/%s: %w", owner, repo, err)
	}

	runs := make([]WorkflowRun, 0, len(result.WorkflowRuns))
	for _, r := range result.WorkflowRuns {
		run := WorkflowRun{
			ID:           r.GetID(),
			Name:         r.GetName(),
			DisplayTitle: r.GetDisplayTitle(),
			Status:       r.GetStatus(),
			Conclusion:   r.GetConclusion(),
			Branch:       r.GetHeadBranch(),
			Event:        r.GetEvent(),
			HTMLURL:      r.GetHTMLURL(),
			CreatedAt:    r.GetCreatedAt().Time,
			RepoOwner:    owner,
			RepoName:     repo,
		}
		if r.Actor != nil {
			run.ActorLogin = r.Actor.GetLogin()
			run.ActorAvatarURL = r.Actor.GetAvatarURL()
		}
		if !r.GetCreatedAt().IsZero() && !r.GetUpdatedAt().IsZero() {
			run.Duration = r.GetUpdatedAt().Time.Sub(r.GetCreatedAt().Time)
		}
		runs = append(runs, run)
	}
	return runs, nil
}
