package phui

import (
	"sort"
	"strings"
)

// CurtainPanel renders a single panel box in a curtain sidebar.
type CurtainPanel struct {
	header string
	body   string // pre-rendered HTML
	order  int
}

// NewCurtainPanel creates an empty curtain panel.
func NewCurtainPanel() *CurtainPanel {
	return &CurtainPanel{}
}

// Header sets the panel title.
func (p *CurtainPanel) Header(h string) *CurtainPanel {
	p.header = h
	return p
}

// Body sets the panel body (raw HTML, not escaped).
func (p *CurtainPanel) Body(html string) *CurtainPanel {
	p.body = html
	return p
}

// Order sets the sort order (lower = first).
func (p *CurtainPanel) Order(o int) *CurtainPanel {
	p.order = o
	return p
}

// Render produces the HTML for this panel.
func (p *CurtainPanel) Render() string {
	var b strings.Builder
	b.WriteString(`<div class="mood-curtain-box">`)
	if p.header != "" {
		b.WriteString(`<div class="mood-curtain-title">`)
		b.WriteString(esc(p.header))
		b.WriteString(`</div>`)
	}
	b.WriteString(p.body)
	b.WriteString(`</div>`)
	return b.String()
}

// Curtain renders a sidebar column with action list and panels.
type Curtain struct {
	panels     []*CurtainPanel
	actionList *ActionList
}

// NewCurtain creates an empty curtain.
func NewCurtain() *Curtain {
	return &Curtain{}
}

// AddPanel appends a panel to the curtain.
func (c *Curtain) AddPanel(p *CurtainPanel) *Curtain {
	c.panels = append(c.panels, p)
	return c
}

// NewPanel creates a panel, adds it, and returns it for configuration.
func (c *Curtain) NewPanel() *CurtainPanel {
	p := NewCurtainPanel()
	c.panels = append(c.panels, p)
	return p
}

// SetActionList sets the action list rendered before panels.
func (c *Curtain) SetActionList(l *ActionList) *Curtain {
	c.actionList = l
	return c
}

// Render produces the HTML for the curtain.
func (c *Curtain) Render() string {
	var b strings.Builder

	if c.actionList != nil {
		b.WriteString(c.actionList.Render())
	}

	sorted := make([]*CurtainPanel, len(c.panels))
	copy(sorted, c.panels)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].order < sorted[j].order
	})

	for _, p := range sorted {
		b.WriteString(p.Render())
	}

	return b.String()
}
