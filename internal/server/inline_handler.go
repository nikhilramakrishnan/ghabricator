package server

import (
	"sync"
	"sync/atomic"
)

// Draft inline comment stored in memory between op=new and op=save.
type inlineDraft struct {
	Owner  string
	Repo   string
	Number int
	Path   string
	Line   int
	Side   string // LEFT or RIGHT
}

var (
	draftMu    sync.Mutex
	draftSeq   atomic.Int64
	draftStore = make(map[int64]*inlineDraft)
)

func nextDraftID() int64 {
	return draftSeq.Add(1)
}
