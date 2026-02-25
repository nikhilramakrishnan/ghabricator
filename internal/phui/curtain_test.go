package phui

import (
	"strings"
	"testing"
)

func TestCurtainPanel_WithHeaderAndBody(t *testing.T) {
	html := NewCurtainPanel().
		Header("Details").
		Body("<p>Some content</p>").
		Render()

	if !strings.Contains(html, `class="mood-curtain-box"`) {
		t.Error("missing mood-curtain-box class")
	}
	if !strings.Contains(html, `<div class="mood-curtain-title">Details</div>`) {
		t.Error("missing header")
	}
	if !strings.Contains(html, "<p>Some content</p>") {
		t.Error("missing body")
	}
}

func TestCurtainPanel_WithoutHeader(t *testing.T) {
	html := NewCurtainPanel().
		Body("<span>no title</span>").
		Render()

	if strings.Contains(html, "mood-curtain-title") {
		t.Error("should not render title div when header is empty")
	}
	if !strings.Contains(html, "<span>no title</span>") {
		t.Error("missing body")
	}
}

func TestCurtainPanel_BodyNotEscaped(t *testing.T) {
	html := NewCurtainPanel().
		Body(`<a href="/link">click</a>`).
		Render()

	if !strings.Contains(html, `<a href="/link">click</a>`) {
		t.Error("body HTML should not be escaped")
	}
}

func TestCurtainPanel_HeaderEscaped(t *testing.T) {
	html := NewCurtainPanel().
		Header(`<script>alert("xss")</script>`).
		Render()

	if strings.Contains(html, "<script>") {
		t.Error("header should be HTML-escaped")
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Error("header should contain escaped script tag")
	}
}

func TestCurtain_PanelsSortedByOrder(t *testing.T) {
	c := NewCurtain()
	c.AddPanel(NewCurtainPanel().Header("Third").Order(30))
	c.AddPanel(NewCurtainPanel().Header("First").Order(10))
	c.AddPanel(NewCurtainPanel().Header("Second").Order(20))

	html := c.Render()

	i1 := strings.Index(html, "First")
	i2 := strings.Index(html, "Second")
	i3 := strings.Index(html, "Third")

	if i1 > i2 || i2 > i3 {
		t.Errorf("panels not sorted by order: First@%d Second@%d Third@%d", i1, i2, i3)
	}
}

func TestCurtain_WithActionList(t *testing.T) {
	al := NewActionList()
	al.AddAction(NewAction("Edit"))

	c := NewCurtain().SetActionList(al)
	c.AddPanel(NewCurtainPanel().Header("Info").Body("<p>details</p>"))

	html := c.Render()

	alIdx := strings.Index(html, "phabricator-action-list-view")
	panelIdx := strings.Index(html, "mood-curtain-box")

	if alIdx == -1 {
		t.Fatal("action list not rendered")
	}
	if panelIdx == -1 {
		t.Fatal("panel not rendered")
	}
	if alIdx > panelIdx {
		t.Error("action list should appear before panels")
	}
}

func TestCurtain_NewPanel(t *testing.T) {
	c := NewCurtain()
	p := c.NewPanel()
	p.Header("Created").Body("<p>via NewPanel</p>")

	html := c.Render()

	if !strings.Contains(html, "Created") {
		t.Error("panel created via NewPanel should be rendered")
	}
	if !strings.Contains(html, "<p>via NewPanel</p>") {
		t.Error("panel body should be rendered")
	}
}

func TestCurtain_Empty(t *testing.T) {
	html := NewCurtain().Render()
	if html != "" {
		t.Errorf("empty curtain should produce empty string, got: %q", html)
	}
}

func TestCurtain_InsertionOrderPreserved(t *testing.T) {
	c := NewCurtain()
	c.AddPanel(NewCurtainPanel().Header("Alpha"))
	c.AddPanel(NewCurtainPanel().Header("Beta"))
	c.AddPanel(NewCurtainPanel().Header("Gamma"))

	html := c.Render()

	iA := strings.Index(html, "Alpha")
	iB := strings.Index(html, "Beta")
	iG := strings.Index(html, "Gamma")

	if iA > iB || iB > iG {
		t.Error("panels with same order should preserve insertion order")
	}
}
