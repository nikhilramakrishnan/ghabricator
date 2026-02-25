package github

import "time"

type PullRequest struct {
	Number    int
	Title     string
	Body      string
	State     string // "open", "closed"
	Draft     bool
	Merged    bool
	Author    User
	CreatedAt time.Time
	UpdatedAt time.Time
	Labels    []Label
	Reviewers []User
	Head      Ref
	Base      Ref
	Additions int
	Deletions int
	ChangedFiles int
}

type User struct {
	Login     string
	AvatarURL string
}

type Label struct {
	Name  string
	Color string
}

type Ref struct {
	Ref  string
	SHA  string
	Repo string // "owner/repo"
}

type Review struct {
	ID        int64
	Author    User
	State     string // APPROVED, CHANGES_REQUESTED, COMMENTED, DISMISSED, PENDING
	Body      string
	CreatedAt time.Time
}

type ReviewComment struct {
	ID        int64
	Author    User
	Body      string
	Path      string
	Line      int
	Side      string // LEFT or RIGHT
	CreatedAt time.Time
	UpdatedAt time.Time
	InReplyTo int64
	DiffHunk  string
}

type InlineCommentRequest struct {
	Body string
	Path string
	Line int
	Side string // LEFT or RIGHT
}

type IssueComment struct {
	ID        int64
	Author    User
	Body      string
	CreatedAt time.Time
}

type Gist struct {
	ID          string
	Description string
	Public      bool
	Files       map[string]GistFile
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Owner       User
	HTMLURL     string
}

type GistFile struct {
	Filename string
	Language string
	Content  string
	Size     int
}
