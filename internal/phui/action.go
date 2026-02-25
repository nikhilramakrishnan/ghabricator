package phui

import "strings"

// ActionView renders a single action item in an action list.
type ActionView struct {
	name     string
	icon     *Icon
	href     string
	disabled bool
	extra    []string
}

// NewAction creates an action with the given label.
func NewAction(name string) *ActionView {
	return &ActionView{name: name}
}

// SetIcon sets the icon for this action.
func (a *ActionView) SetIcon(i *Icon) *ActionView {
	a.icon = i
	return a
}

// Href sets the link target.
func (a *ActionView) Href(h string) *ActionView {
	a.href = h
	return a
}

// Disabled marks the action as disabled.
func (a *ActionView) Disabled(d bool) *ActionView {
	a.disabled = d
	return a
}

// AddClass appends an extra CSS class.
func (a *ActionView) AddClass(c string) *ActionView {
	a.extra = append(a.extra, c)
	return a
}

// Render produces the HTML for this action item.
func (a *ActionView) Render() string {
	var b strings.Builder

	liClass := classes(
		"phabricator-action-view",
		"action-has-icon",
		cond(a.disabled, "phabricator-action-view-disabled", ""),
		strings.Join(a.extra, " "),
	)
	b.WriteString(`<li` + attr("class", liClass) + `>`)

	// Icon span — always rendered for alignment.
	iconClass := "phabricator-action-view-icon phui-icon-view phui-font-fa"
	if a.icon != nil {
		iconClass += " " + esc(a.icon.name)
	}
	iconSpan := `<span` + attr("class", iconClass) + `></span>`

	// Inner element: <a> with href, or <span> if no href / disabled.
	if a.href != "" && !a.disabled {
		b.WriteString(`<a href="` + esc(a.href) + `"` +
			attr("class", "phabricator-action-view-item") + `>`)
		b.WriteString(iconSpan)
		b.WriteString(esc(a.name))
		b.WriteString(`</a>`)
	} else {
		attrs := attr("class", "phabricator-action-view-item")
		if a.disabled {
			attrs += ` disabled`
		}
		b.WriteString(`<span` + attrs + `>`)
		b.WriteString(iconSpan)
		b.WriteString(esc(a.name))
		b.WriteString(`</span>`)
	}

	b.WriteString(`</li>`)
	return b.String()
}

// ActionList renders an unordered list of ActionViews.
type ActionList struct {
	actions []*ActionView
}

// NewActionList creates an empty action list.
func NewActionList() *ActionList {
	return &ActionList{}
}

// AddAction appends an action to the list.
func (l *ActionList) AddAction(a *ActionView) *ActionList {
	l.actions = append(l.actions, a)
	return l
}

// Render produces the HTML for the action list.
func (l *ActionList) Render() string {
	var b strings.Builder
	b.WriteString(`<ul class="phabricator-action-list-view">`)
	for _, a := range l.actions {
		b.WriteString(a.Render())
	}
	b.WriteString(`</ul>`)
	return b.String()
}

// cond returns t if condition is true, f otherwise.
func cond(c bool, t, f string) string {
	if c {
		return t
	}
	return f
}
