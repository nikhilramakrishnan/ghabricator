package phui

import "strings"

// TimelineEvent renders a single event in a PHUI timeline.
type TimelineEvent struct {
	author    string
	avatarURL string
	title     string // pre-rendered HTML (not escaped)
	body      string // pre-rendered HTML (makes event "major")
	icon      *Icon
	iconColor string
	dateStr   string
	anchor    string
	extra     []string
}

// NewTimelineEvent creates a new empty timeline event.
func NewTimelineEvent() *TimelineEvent {
	return &TimelineEvent{}
}

func (e *TimelineEvent) Author(a string) *TimelineEvent    { e.author = a; return e }
func (e *TimelineEvent) AvatarURL(url string) *TimelineEvent { e.avatarURL = url; return e }
func (e *TimelineEvent) Title(t string) *TimelineEvent      { e.title = t; return e }
func (e *TimelineEvent) Body(html string) *TimelineEvent    { e.body = html; return e }
func (e *TimelineEvent) SetIcon(i *Icon) *TimelineEvent     { e.icon = i; return e }
func (e *TimelineEvent) IconColor(c string) *TimelineEvent  { e.iconColor = c; return e }
func (e *TimelineEvent) Date(d string) *TimelineEvent       { e.dateStr = d; return e }
func (e *TimelineEvent) Anchor(a string) *TimelineEvent     { e.anchor = a; return e }
func (e *TimelineEvent) AddClass(c string) *TimelineEvent   { e.extra = append(e.extra, c); return e }
func (e *TimelineEvent) IsMajor() bool                      { return e.body != "" }

// Render produces the timeline event HTML.
func (e *TimelineEvent) Render() string {
	var b strings.Builder

	eventClass := "phui-timeline-minor-event"
	if e.IsMajor() {
		eventClass = "phui-timeline-major-event"
	}
	cls := classes("phui-timeline-event-view", eventClass)
	for _, c := range e.extra {
		cls = classes(cls, c)
	}

	b.WriteString(`<div class="phui-timeline-shell">`)
	b.WriteString(`<div` + attr("class", cls) + attr("id", e.anchor) + `>`)
	b.WriteString(`<div class="phui-timeline-content">`)

	// Avatar
	if e.avatarURL != "" {
		b.WriteString(`<div class="phui-timeline-image" style="background-image:url(` + esc(e.avatarURL) + `)"></div>`)
	}

	b.WriteString(`<div class="phui-timeline-wedge"></div>`)
	b.WriteString(`<div class="phui-timeline-group">`)

	if e.IsMajor() {
		b.WriteString(`<div class="phui-timeline-inner-content">`)
	}

	// Title
	titleClass := "phui-timeline-title"
	if e.icon != nil {
		titleClass += " phui-timeline-title-with-icon"
	}
	b.WriteString(`<div class="` + titleClass + `">`)

	// Icon
	if e.icon != nil {
		color := e.iconColor
		if color == "" {
			color = "grey"
		}
		b.WriteString(`<span class="phui-timeline-icon-fill fill-has-color phui-timeline-icon-fill-` + esc(color) + `">`)
		b.WriteString(`<span class="phui-icon-view phui-font-fa ` + esc(e.icon.name) + ` phui-timeline-icon"></span>`)
		b.WriteString(`</span>`)
	}

	// Title content (raw HTML)
	b.WriteString(e.title)

	// Date
	if e.dateStr != "" {
		b.WriteString(`<span class="phui-timeline-extra">` + esc(e.dateStr) + `</span>`)
	}

	b.WriteString(`</div>`) // title

	// Body (major only)
	if e.IsMajor() {
		b.WriteString(`<div class="phui-timeline-core-content">` + e.body + `</div>`)
		b.WriteString(`</div>`) // inner-content
	}

	b.WriteString(`</div>`) // group
	b.WriteString(`</div>`) // content
	b.WriteString(`</div>`) // event-view
	b.WriteString(`</div>`) // shell

	return b.String()
}

// Timeline renders a list of timeline events with spacers.
type Timeline struct {
	events    []*TimelineEvent
	terminate bool
}

// NewTimeline creates a new empty timeline.
func NewTimeline() *Timeline {
	return &Timeline{}
}

func (t *Timeline) AddEvent(e *TimelineEvent) *Timeline { t.events = append(t.events, e); return t }
func (t *Timeline) Terminate(b bool) *Timeline          { t.terminate = b; return t }

// Render produces the timeline HTML.
func (t *Timeline) Render() string {
	var b strings.Builder
	b.WriteString(`<div class="phui-timeline-view">`)
	b.WriteString(`<div class="phui-timeline-event-view phui-timeline-spacer"></div>`)
	for _, e := range t.events {
		b.WriteString(e.Render())
		b.WriteString(`<div class="phui-timeline-event-view phui-timeline-spacer"></div>`)
	}
	if t.terminate {
		b.WriteString(`<div class="phui-timeline-event-view the-worlds-end"></div>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}
