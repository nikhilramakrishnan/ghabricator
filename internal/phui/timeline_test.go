package phui

import (
	"strings"
	"testing"
)

func TestMinorEvent(t *testing.T) {
	e := NewTimelineEvent().Title("changed status")
	html := e.Render()
	if !strings.Contains(html, "phui-timeline-minor-event") {
		t.Fatal("expected minor-event class")
	}
	if strings.Contains(html, "phui-timeline-major-event") {
		t.Fatal("unexpected major-event class")
	}
	if strings.Contains(html, "phui-timeline-core-content") {
		t.Fatal("minor event should not have core-content div")
	}
}

func TestMajorEvent(t *testing.T) {
	e := NewTimelineEvent().Title("added a comment").Body("<p>Hello</p>")
	html := e.Render()
	if !strings.Contains(html, "phui-timeline-major-event") {
		t.Fatal("expected major-event class")
	}
	if strings.Contains(html, "phui-timeline-minor-event") {
		t.Fatal("unexpected minor-event class")
	}
	if !strings.Contains(html, "phui-timeline-core-content") {
		t.Fatal("major event should have core-content div")
	}
	if !strings.Contains(html, "<p>Hello</p>") {
		t.Fatal("body not rendered")
	}
}

func TestEventWithIcon(t *testing.T) {
	icon := NewIcon("fa-check")
	e := NewTimelineEvent().Title("closed").SetIcon(icon).IconColor("green")
	html := e.Render()
	if !strings.Contains(html, "phui-timeline-icon-fill-green") {
		t.Fatal("expected icon fill color green")
	}
	if !strings.Contains(html, "fa-check") {
		t.Fatal("expected icon name")
	}
	if !strings.Contains(html, "phui-timeline-title-with-icon") {
		t.Fatal("expected title-with-icon class")
	}
}

func TestEventWithAvatar(t *testing.T) {
	e := NewTimelineEvent().AvatarURL("https://example.com/avatar.png")
	html := e.Render()
	if !strings.Contains(html, "phui-timeline-image") {
		t.Fatal("expected image div")
	}
	if !strings.Contains(html, "background-image:url(https://example.com/avatar.png)") {
		t.Fatal("expected background-image style")
	}
}

func TestEventWithoutAvatar(t *testing.T) {
	e := NewTimelineEvent().Title("test")
	html := e.Render()
	if strings.Contains(html, "phui-timeline-image") {
		t.Fatal("should not have image div when no avatar")
	}
}

func TestEventWithDate(t *testing.T) {
	e := NewTimelineEvent().Title("test").Date("Jan 5, 2025")
	html := e.Render()
	if !strings.Contains(html, `<span class="phui-timeline-extra">Jan 5, 2025</span>`) {
		t.Fatal("expected date in extra span")
	}
}

func TestEventWithAnchor(t *testing.T) {
	e := NewTimelineEvent().Title("test").Anchor("comment-3")
	html := e.Render()
	if !strings.Contains(html, `id="comment-3"`) {
		t.Fatal("expected anchor id")
	}
}

func TestTimelineMultipleEvents(t *testing.T) {
	tl := NewTimeline()
	tl.AddEvent(NewTimelineEvent().Title("one"))
	tl.AddEvent(NewTimelineEvent().Title("two"))
	html := tl.Render()
	// Should have 3 spacers: initial + after each event
	if strings.Count(html, "phui-timeline-spacer") != 3 {
		t.Fatalf("expected 3 spacers, got %d", strings.Count(html, "phui-timeline-spacer"))
	}
	if !strings.Contains(html, "one") || !strings.Contains(html, "two") {
		t.Fatal("expected both events")
	}
}

func TestTimelineTerminated(t *testing.T) {
	tl := NewTimeline().Terminate(true)
	html := tl.Render()
	if !strings.Contains(html, "the-worlds-end") {
		t.Fatal("expected the-worlds-end div")
	}
}

func TestTimelineEmpty(t *testing.T) {
	tl := NewTimeline()
	html := tl.Render()
	if strings.Count(html, "phui-timeline-spacer") != 1 {
		t.Fatal("empty timeline should have exactly 1 spacer")
	}
	if strings.Contains(html, "the-worlds-end") {
		t.Fatal("non-terminated timeline should not have the-worlds-end")
	}
}

func TestTitleIsRawHTML(t *testing.T) {
	e := NewTimelineEvent().Title(`<a href="/u/alice">alice</a> changed status`)
	html := e.Render()
	if !strings.Contains(html, `<a href="/u/alice">alice</a> changed status`) {
		t.Fatal("title should be raw HTML, not escaped")
	}
}

func TestBodyIsRawHTML(t *testing.T) {
	e := NewTimelineEvent().Title("comment").Body(`<em>bold</em>`)
	html := e.Render()
	if !strings.Contains(html, `<em>bold</em>`) {
		t.Fatal("body should be raw HTML, not escaped")
	}
}

func TestAvatarURLIsEscaped(t *testing.T) {
	e := NewTimelineEvent().AvatarURL(`https://example.com/a"b`)
	html := e.Render()
	if strings.Contains(html, `a"b`) {
		t.Fatal("avatar URL should be escaped")
	}
	if !strings.Contains(html, `a&#34;b`) {
		t.Fatal("expected escaped quote in avatar URL")
	}
}
