package herald

import (
	"time"
)

// ConditionType identifies what a condition checks.
type ConditionType string

const (
	CondFilePath   ConditionType = "file_path"   // glob match on changed file paths
	CondAuthor     ConditionType = "author"       // PR author login matches
	CondTitle      ConditionType = "title"        // PR title contains substring
	CondLabel      ConditionType = "label"        // PR has label
	CondBaseBranch ConditionType = "base_branch"  // PR base branch matches
)

// ActionType identifies what an action does.
type ActionType string

const (
	ActionAddReviewer ActionType = "add_reviewer" // request review from user
	ActionAddLabel    ActionType = "add_label"    // add label to PR
	ActionPostComment ActionType = "post_comment" // post a comment on PR
)

// Condition is a single predicate in a rule.
type Condition struct {
	Type  ConditionType `json:"type"`
	Value string        `json:"value"`
}

// Action is a single effect triggered by a rule.
type Action struct {
	Type  ActionType `json:"type"`
	Value string     `json:"value"`
}

// Rule is a Herald automation rule.
type Rule struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	AuthorLogin string      `json:"author_login"`
	Conditions  []Condition `json:"conditions"`
	Actions     []Action    `json:"actions"`
	MustMatchAll bool       `json:"must_match_all"` // true=AND, false=OR
	Disabled    bool        `json:"disabled"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// RuleMatch records that a rule fired and which actions it produced.
type RuleMatch struct {
	Rule    *Rule
	Actions []Action
}

// PRContext contains the PR metadata needed for rule evaluation.
type PRContext struct {
	Author       string
	Title        string
	Labels       []string
	BaseBranch   string
	ChangedFiles []string
}
