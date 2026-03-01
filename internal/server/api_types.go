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
	PR               APIPRDetail                    `json:"pr"`
	Changesets       []APIChangeset                 `json:"changesets"`
	CommentsByPath   map[string][]APIReviewComment  `json:"commentsByPath"`
	Reviews          []APIReview                    `json:"reviews"`
	IssueComments    []APIIssueComment              `json:"issueComments"`
	CheckRuns        []APICheckRun                  `json:"checkRuns"`
	Timeline         []APITimelineEvent             `json:"timeline"`
	HeraldMatches    []APIHeraldMatch               `json:"heraldMatches,omitempty"`
	Commits          []APICommit                    `json:"commits"`
	ViewerPermission string                         `json:"viewerPermission"`
}

type APICommit struct {
	SHA     string  `json:"sha"`
	Message string  `json:"message"`
	Author  APIUser `json:"author"`
	Date    string  `json:"date"`
}

type APIPRDetail struct {
	Number       int        `json:"number"`
	Title        string     `json:"title"`
	Body         string     `json:"body"`
	BodyRaw      string     `json:"bodyRaw,omitempty"`
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
	ID        int64              `json:"id"`
	Author    APIUser            `json:"author"`
	Body      string             `json:"body"`
	BodyRaw   string             `json:"bodyRaw,omitempty"`
	Path      string             `json:"path"`
	Line      int                `json:"line"`
	Side      string             `json:"side"`
	CreatedAt time.Time          `json:"createdAt"`
	InReplyTo int64              `json:"inReplyTo,omitempty"`
	Reactions []APIReaction `json:"reactions,omitempty"`
}

type APIReaction struct {
	Emoji string `json:"emoji"`
	Count int    `json:"count"`
}

type APIReview struct {
	ID        int64     `json:"id"`
	Author    APIUser   `json:"author"`
	State     string    `json:"state"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type APIIssueComment struct {
	ID        int64         `json:"id"`
	Author    APIUser       `json:"author"`
	Body      string        `json:"body"`
	CreatedAt time.Time     `json:"createdAt"`
	Reactions []APIReaction `json:"reactions,omitempty"`
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
	Author      APIUser       `json:"author"`
	Action      string        `json:"action"`
	Body        string        `json:"body,omitempty"`
	BodyRaw     string        `json:"bodyRaw,omitempty"`
	CreatedAt   time.Time     `json:"createdAt"`
	IconClass   string        `json:"iconClass"`
	IconColor   string        `json:"iconColor"`
	CommentID   int64         `json:"commentID,omitempty"`
	CommentType string        `json:"commentType,omitempty"`
	Reactions   []APIReaction `json:"reactions,omitempty"`
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

// --- Inline Comment API types ---

type APIInlineRequest struct {
	Operation string `json:"operation"` // new, save, edit, cancel, delete, done
	Owner     string `json:"owner"`
	Repo      string `json:"repo"`
	Number    int    `json:"number"`
	Path      string `json:"path"`
	Line      int    `json:"line"`
	Side      string `json:"side"` // LEFT or RIGHT
	Body      string `json:"body"`
	CommentID int64  `json:"commentID"`
	InReplyTo int64  `json:"inReplyTo,omitempty"`
}

type APIInlineComment struct {
	ID        int64   `json:"id"`
	Author    APIUser `json:"author"`
	Body      string  `json:"body"`
	Path      string  `json:"path"`
	Line      int     `json:"line"`
	Side      string  `json:"side"`
}

// --- Review/Merge/Close API types ---

type APIReviewRequest struct {
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	Number   int    `json:"number"`
	Action   string `json:"action"` // APPROVE, REQUEST_CHANGES, COMMENT
	Body     string `json:"body"`
}

type APIMergeRequest struct {
	Owner       string `json:"owner"`
	Repo        string `json:"repo"`
	Number      int    `json:"number"`
	MergeMethod string `json:"mergeMethod"` // merge, squash, rebase
}

type APICloseRequest struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Number int    `json:"number"`
	State  string `json:"state"` // closed, open
}

// --- Reaction API types ---

type APIReactionRequest struct {
	Owner       string `json:"owner"`
	Repo        string `json:"repo"`
	CommentID   int64  `json:"commentID"`
	Content     string `json:"content"`     // +1, -1, laugh, confused, heart, hooray, rocket, eyes
	CommentType string `json:"commentType"` // "issue" or "review"
}

// --- Repos API types ---

type APIRepoSummary struct {
	Name        string `json:"name"`
	FullName    string `json:"fullName"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	Private     bool   `json:"private"`
	Fork        bool   `json:"fork"`
	Archived    bool   `json:"archived"`
	AvatarURL   string `json:"avatarURL"`
	UpdatedAt   string `json:"updatedAt"`
}

type APIRepoInfo struct {
	FullName      string `json:"fullName"`
	Description   string `json:"description"`
	DefaultBranch string `json:"defaultBranch"`
	Private       bool   `json:"private"`
	HTMLURL       string `json:"htmlURL"`
	Stars         int    `json:"stars"`
	Forks         int    `json:"forks"`
}

type APIRepoEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"` // file or dir
	Size int    `json:"size"`
}

type APIRepoTreeResponse struct {
	Entries  []APIRepoEntry `json:"entries"`
	RepoInfo APIRepoInfo    `json:"repoInfo"`
}

type APIRepoFile struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Size    int      `json:"size"`
	Lines   []string `json:"lines,omitempty"`   // syntax-highlighted lines for code
	RawURL  string   `json:"rawURL,omitempty"`  // for images
	HTMLURL string   `json:"htmlURL"`
}

type APIRepoFileResponse struct {
	File     APIRepoFile `json:"file"`
	RepoInfo APIRepoInfo `json:"repoInfo"`
}

// --- Blame API types ---

type APIBlameRange struct {
	StartLine       int    `json:"startLine"`
	EndLine         int    `json:"endLine"`
	CommitOID       string `json:"commitOID"`
	CommitShort     string `json:"commitShort"`
	Message         string `json:"message"`
	AuthorLogin     string `json:"authorLogin"`
	AuthorAvatarURL string `json:"authorAvatarURL"`
	AuthorName      string `json:"authorName"`
	AuthoredDate    string `json:"authoredDate"`
}

type APIBlameResponse struct {
	Ranges []APIBlameRange `json:"ranges"`
}

// --- Paste API types ---

type APIPasteSummary struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Language    string `json:"language"`
	Public      bool   `json:"public"`
	CreatedAt   string `json:"createdAt"`
}

type APIPasteFile struct {
	Filename string   `json:"filename"`
	Language string   `json:"language"`
	Size     int      `json:"size"`
	Lines    []string `json:"lines"` // syntax-highlighted
}

type APIPasteDetail struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Public      bool           `json:"public"`
	Owner       APIUser        `json:"owner"`
	HTMLURL     string         `json:"htmlURL"`
	Files       []APIPasteFile `json:"files"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
}

type APIPasteCreateRequest struct {
	Title    string `json:"title"`
	Language string `json:"language"`
	Content  string `json:"content"`
	Public   bool   `json:"public"`
}

type APIPasteCreateResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// --- Search API types ---

type APISearchPR struct {
	Number        int        `json:"number"`
	Title         string     `json:"title"`
	Repo          string     `json:"repo"`
	State         string     `json:"state"`
	Author        string     `json:"author"`
	AvatarURL     string     `json:"avatarURL"`
	Labels        []APILabel `json:"labels,omitempty"`
	Body          string     `json:"body,omitempty"`
	CommentsCount int        `json:"commentsCount"`
	Draft         bool       `json:"draft"`
	CreatedAt     string     `json:"createdAt"`
	UpdatedAt     string     `json:"updatedAt"`
}

type APISearchIssue struct {
	Number        int        `json:"number"`
	Title         string     `json:"title"`
	Repo          string     `json:"repo"`
	State         string     `json:"state"`
	Author        string     `json:"author"`
	AvatarURL     string     `json:"avatarURL"`
	Labels        []APILabel `json:"labels,omitempty"`
	Body          string     `json:"body,omitempty"`
	CommentsCount int        `json:"commentsCount"`
	CreatedAt     string     `json:"createdAt"`
	UpdatedAt     string     `json:"updatedAt"`
}

type APISearchCodeResult struct {
	Repo       string `json:"repo"`
	Path       string `json:"path"`
	Fragment   string `json:"fragment"`
	Language   string `json:"language,omitempty"`
	MatchCount int    `json:"matchCount"`
	HTMLURL    string `json:"htmlURL,omitempty"`
}

type APISearchRepoResult struct {
	FullName    string   `json:"fullName"`
	Description string   `json:"description"`
	Stars       int      `json:"stars"`
	Forks       int      `json:"forks"`
	Language    string   `json:"language"`
	LangColor   string   `json:"langColor,omitempty"`
	AvatarURL   string   `json:"avatarURL"`
	UpdatedAt   string   `json:"updatedAt"`
	Topics      []string `json:"topics,omitempty"`
}

type APISearchResponse struct {
	Counts map[string]int        `json:"counts"`
	PRs    []APISearchPR         `json:"prs,omitempty"`
	Issues []APISearchIssue      `json:"issues,omitempty"`
	Code   []APISearchCodeResult `json:"code,omitempty"`
	Repos  []APISearchRepoResult `json:"repos,omitempty"`
}

// --- Workflow Runs API types ---

type APIWorkflowRunsResponse struct {
	Runs []APIWorkflowRun `json:"runs"`
}

type APIWorkflowRun struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	DisplayTitle string  `json:"displayTitle"`
	Status       string  `json:"status"`
	Conclusion   string  `json:"conclusion"`
	Branch       string  `json:"branch"`
	Event        string  `json:"event"`
	Actor        APIUser `json:"actor"`
	RepoOwner    string  `json:"repoOwner"`
	RepoName     string  `json:"repoName"`
	DurationMs   int64   `json:"durationMs"`
	HTMLURL      string  `json:"htmlURL"`
	CreatedAt    string  `json:"createdAt"`
}
