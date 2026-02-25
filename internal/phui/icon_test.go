package phui

import (
	"strings"
	"testing"
)

func TestIconBasic(t *testing.T) {
	out := NewIcon("fa-check").Render()
	if !strings.Contains(out, "fa-check") {
		t.Fatalf("expected fa-check in %q", out)
	}
	if !strings.Contains(out, "phui-icon-view") {
		t.Fatalf("expected phui-icon-view in %q", out)
	}
	if !strings.Contains(out, "phui-font-fa") {
		t.Fatalf("expected phui-font-fa in %q", out)
	}
}

func TestIconExtraClass(t *testing.T) {
	out := NewIcon("fa-check").AddClass("mrs").Render()
	if !strings.Contains(out, "fa-check") {
		t.Fatalf("expected fa-check in %q", out)
	}
	if !strings.Contains(out, "mrs") {
		t.Fatalf("expected mrs in %q", out)
	}
}

func TestIconMultipleExtraClasses(t *testing.T) {
	out := NewIcon("fa-github").AddClass("mrs").AddClass("phui-timeline-icon").Render()
	for _, want := range []string{"fa-github", "mrs", "phui-timeline-icon"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in %q", want, out)
		}
	}
}

func TestIconColor(t *testing.T) {
	out := NewIcon("fa-star").Color("gold").Render()
	if !strings.Contains(out, `style="color:gold"`) {
		t.Fatalf("expected color style in %q", out)
	}
}

func TestIconEmptyName(t *testing.T) {
	out := NewIcon("").Render()
	if strings.Contains(out, "  ") {
		t.Fatalf("double space in class attr: %q", out)
	}
	if !strings.HasPrefix(out, "<span") {
		t.Fatalf("expected valid span: %q", out)
	}
}

func TestIconXSS(t *testing.T) {
	out := NewIcon("<script>").Render()
	if strings.Contains(out, "<script>") {
		t.Fatalf("unescaped script tag in %q", out)
	}
}
