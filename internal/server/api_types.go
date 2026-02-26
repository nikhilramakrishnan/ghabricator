package server

import "time"

// --- Dashboard API types ---

type APIDashboardResponse struct {
	Authored        []APIPRSummary `json:"authored"`
	ReviewRequested []APIPRSummary `json:"reviewRequested"`
}

type APIPRSummary struct {
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	Owner     string     `json:"owner"`
	Repo      string     `json:"repo"`
	Author    APIUser    `json:"author"`
	Draft     bool       `json:"draft"`
	Labels    []APILabel `json:"labels,omitempty"`
	Reviewers []APIUser  `json:"reviewers,omitempty"`
	UpdatedAt string     `json:"updatedAt"`
}

type APIUser struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatarURL"`
}

type APILabel struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// --- PR Detail API types ---

type APIPRDetailResponse struct {
	PR             APIPRDetail                    `json:"pr"`
	Changesets     []APIChangeset                 `json:"changesets"`
	CommentsByPath map[string][]APIReviewComment  `json:"commentsByPath"`
	Reviews        []APIReview                    `json:"reviews"`
	IssueComments  []APIIssueComment              `json:"issueComments"`
	CheckRuns      []APICheckRun                  `json:"checkRuns"`
	Timeline       []APITimelineEvent             `json:"timeline"`
	HeraldMatches  []APIHeraldMatch               `json:"heraldMatches,omitempty"`
}

type APIPRDetail struct {
	Number       int        `json:"number"`
	Title        string     `json:"title"`
	Body         string     `json:"body"`
	State        string     `json:"state"`
	Draft        bool       `json:"draft"`
	Merged       bool       `json:"merged"`
	Author       APIUser    `json:"author"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	Labels       []APILabel `json:"labels"`
	Reviewers    []APIUser  `json:"reviewers"`
	Head         APIRef     `json:"head"`
	Base         APIRef     `json:"base"`
	Additions    int        `json:"additions"`
	Deletions    int        `json:"deletions"`
	ChangedFiles int        `json:"changedFiles"`
}

type APIRef struct {
	Ref  string `json:"ref"`
	SHA  string `json:"sha"`
	Repo string `json:"repo"`
}

type APIChangeset struct {
	ID           int          `json:"id"`
	OldName      string       `json:"oldName"`
	NewName      string       `json:"newName"`
	DisplayPath  string       `json:"displayPath"`
	LinesAdded   int          `json:"linesAdded"`
	LinesRemoved int          `json:"linesRemoved"`
	IsNew        bool         `json:"isNew"`
	IsDeleted    bool         `json:"isDeleted"`
	IsRenamed    bool         `json:"isRenamed"`
	IsBinary     bool         `json:"isBinary"`
	Rows         []APIDiffRow `json:"rows"`
}

type APIDiffRow struct {
	OldNum     int    `json:"oldNum"`
	NewNum     int    `json:"newNum"`
	OldClass   string `json:"oldClass"`
	NewClass   string `json:"newClass"`
	OldContent string `json:"oldContent"`
	NewContent string `json:"newContent"`
	IsContext  bool   `json:"isContext"`
}

type APIReviewComment struct {
	ID        int64     `json:"id"`
	Author    APIUser   `json:"author"`
	Body      string    `json:"body"`
	Path      string    `json:"path"`
	Line      int       `json:"line"`
	Side      string    `json:"side"`
	CreatedAt time.Time `json:"createdAt"`
	InReplyTo int64     `json:"inReplyTo,omitempty"`
}

type APIReview struct {
	ID        int64     `json:"id"`
	Author    APIUser   `json:"author"`
	State     string    `json:"state"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type APIIssueComment struct {
	ID        int64     `json:"id"`
	Author    APIUser   `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type APICheckRun struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	DetailsURL  string    `json:"detailsURL"`
	AppName     string    `json:"appName"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
}

type APITimelineEvent struct {
	Author    APIUser   `json:"author"`
	Action    string    `json:"action"`
	Body      string    `json:"body,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	IconClass string    `json:"iconClass"`
	IconColor string    `json:"iconColor"`
}

type APIHeraldMatch struct {
	RuleID   string            `json:"ruleId"`
	RuleName string            `json:"ruleName"`
	Actions  []APIHeraldAction `json:"actions"`
}

type APIHeraldAction struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
