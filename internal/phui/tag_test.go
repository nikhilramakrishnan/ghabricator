package phui

import (
	"strings"
	"testing"
)

func TestDefaultShadeTag(t *testing.T) {
	html := NewTag("Draft").Color("blue").Render()
	expect := []string{
		"phui-tag-view",
		"phui-tag-type-shade",
		"phui-tag-shade",
		"phui-tag-blue",
		"phui-tag-core",
		">Draft<",
	}
	for _, s := range expect {
		if !strings.Contains(html, s) {
			t.Errorf("missing %q in:\n%s", s, html)
		}
	}
	if !strings.HasPrefix(html, "<span") {
		t.Errorf("expected <span root, got:\n%s", html)
	}
}

func TestTagTypes(t *testing.T) {
	tests := []struct {
		tt       TagType
		color    string
		rootHas  []string
		rootNot  []string
		coreHas  []string
	}{
		{TagShade, "green", []string{"phui-tag-shade", "phui-tag-green"}, nil, nil},
		{TagOutline, "red", []string{"phui-tag-red"}, []string{"phui-tag-shade"}, nil},
		{TagState, "blue", []string{"phui-tag-blue"}, nil, []string{"phui-tag-color-blue"}},
		{TagObject, "violet", nil, []string{"phui-tag-violet"}, []string{"phui-tag-color-violet"}},
		{TagPerson, "indigo", nil, []string{"phui-tag-indigo"}, nil},
	}
	for _, tc := range tests {
		html := NewTag("X").Type(tc.tt).Color(tc.color).Render()
		must := "phui-tag-type-" + string(tc.tt)
		if !strings.Contains(html, must) {
			t.Errorf("[%s] missing %q", tc.tt, must)
		}
		for _, s := range tc.rootHas {
			if !strings.Contains(html, s) {
				t.Errorf("[%s] missing %q in:\n%s", tc.tt, s, html)
			}
		}
		for _, s := range tc.rootNot {
			if strings.Contains(html, s) {
				t.Errorf("[%s] unexpected %q in:\n%s", tc.tt, s, html)
			}
		}
		for _, s := range tc.coreHas {
			if !strings.Contains(html, s) {
				t.Errorf("[%s] missing core %q in:\n%s", tc.tt, s, html)
			}
		}
	}
}

func TestHrefSwitchesToAnchor(t *testing.T) {
	html := NewTag("Link").Href("/foo").Render()
	if !strings.HasPrefix(html, "<a ") {
		t.Errorf("expected <a root, got:\n%s", html)
	}
	if !strings.Contains(html, `href="/foo"`) {
		t.Errorf("missing href in:\n%s", html)
	}
	if !strings.HasSuffix(html, "</a>") {
		t.Errorf("expected </a> closing, got:\n%s", html)
	}
}

func TestClosedWraps(t *testing.T) {
	html := NewTag("Done").Closed(true).Render()
	if !strings.Contains(html, "phui-tag-core-closed") {
		t.Errorf("missing closed wrapper in:\n%s", html)
	}
	if !strings.Contains(html, "phui-tag-core") {
		t.Errorf("missing inner core in:\n%s", html)
	}
}

func TestSlimAddsClass(t *testing.T) {
	html := NewTag("S").Slim(true).Render()
	if !strings.Contains(html, "phui-tag-slim") {
		t.Errorf("missing slim class in:\n%s", html)
	}
}

func TestDotColor(t *testing.T) {
	html := NewTag("Status").DotColor("green").Render()
	if !strings.Contains(html, "phui-tag-dot") {
		t.Errorf("missing dot span in:\n%s", html)
	}
	if !strings.Contains(html, "phui-tag-color-green") {
		t.Errorf("missing dot color in:\n%s", html)
	}
}

func TestIconRendersInside(t *testing.T) {
	icon := NewIcon("fa-check")
	html := NewTag("OK").SetIcon(icon).Render()
	if !strings.Contains(html, "phui-icon-view") {
		t.Errorf("missing icon in:\n%s", html)
	}
	if !strings.Contains(html, "fa-check") {
		t.Errorf("missing icon name in:\n%s", html)
	}
}

func TestBorderAddsClass(t *testing.T) {
	html := NewTag("X").Border("border-none").Render()
	if !strings.Contains(html, "phui-tag-border-none") {
		t.Errorf("missing border class in:\n%s", html)
	}
}

func TestStateColorOnCore(t *testing.T) {
	html := NewTag("Active").Type(TagState).Color("green").Render()
	// Core should have color class
	idx := strings.Index(html, "phui-tag-core")
	if idx == -1 {
		t.Fatalf("no core span in:\n%s", html)
	}
	coreOnward := html[idx:]
	if !strings.Contains(coreOnward, "phui-tag-color-green") {
		t.Errorf("state color not on core in:\n%s", html)
	}
}

func TestXSSEscaping(t *testing.T) {
	html := NewTag(`<script>alert(1)</script>`).Href(`" onclick="alert(1)`).Render()
	if strings.Contains(html, "<script>") {
		t.Errorf("XSS in name not escaped:\n%s", html)
	}
	// href must escape the quote so attribute injection can't break out
	if strings.Contains(html, `" onclick=`) {
		t.Errorf("XSS in href not escaped:\n%s", html)
	}
	if !strings.Contains(html, `&#34;`) {
		t.Errorf("expected escaped quote in href:\n%s", html)
	}
}
