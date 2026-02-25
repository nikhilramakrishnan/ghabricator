package phui

import (
	"strings"
	"testing"
)

func TestHeaderBasic(t *testing.T) {
	out := NewHeader("My Header").Render()
	if !strings.Contains(out, "phui-header-shell") {
		t.Fatalf("missing shell class: %q", out)
	}
	if !strings.Contains(out, "My Header") {
		t.Fatalf("missing text: %q", out)
	}
	if !strings.Contains(out, "phui-header-header") {
		t.Fatalf("missing header class: %q", out)
	}
}

func TestHeaderWithIcon(t *testing.T) {
	out := NewHeader("Tasks").SetIcon(NewIcon("fa-anchor")).Render()
	if !strings.Contains(out, "phui-header-icon") {
		t.Fatalf("missing icon wrapper: %q", out)
	}
	if !strings.Contains(out, "fa-anchor") {
		t.Fatalf("missing icon class: %q", out)
	}
}

func TestHeaderSubheader(t *testing.T) {
	out := NewHeader("Title").Subheader("Subtitle here").Render()
	if !strings.Contains(out, "phui-header-subheader") {
		t.Fatalf("missing subheader class: %q", out)
	}
	if !strings.Contains(out, "Subtitle here") {
		t.Fatalf("missing subheader text: %q", out)
	}
}

func TestHeaderWithTags(t *testing.T) {
	out := NewHeader("Review").
		AddTag(NewTag("Accepted").Color("green")).
		AddTag(NewTag("v2")).
		Render()
	if !strings.Contains(out, "Accepted") {
		t.Fatalf("missing first tag: %q", out)
	}
	if !strings.Contains(out, "v2") {
		t.Fatalf("missing second tag: %q", out)
	}
	if !strings.Contains(out, "phui-tag-view") {
		t.Fatalf("tag Render not called: %q", out)
	}
}

func TestHeaderActionLinks(t *testing.T) {
	out := NewHeader("Config").
		AddActionLink(NewButton("Edit").Color("green")).
		Render()
	if !strings.Contains(out, "phui-header-action-links") {
		t.Fatalf("missing action links wrapper: %q", out)
	}
	if !strings.Contains(out, "Edit") {
		t.Fatalf("missing button text: %q", out)
	}
}

func TestHeaderImage(t *testing.T) {
	out := NewHeader("User").Image("/avatar.png").Render()
	if !strings.Contains(out, "phui-header-image") {
		t.Fatalf("missing image class: %q", out)
	}
	if !strings.Contains(out, "background-image:url(/avatar.png)") {
		t.Fatalf("missing background-image: %q", out)
	}
}

func TestHeaderImageWithHref(t *testing.T) {
	out := NewHeader("User").Image("/avatar.png").ImageHref("/profile").Render()
	if !strings.Contains(out, `href="/profile"`) {
		t.Fatalf("missing image href: %q", out)
	}
}

func TestHeaderTall(t *testing.T) {
	out := NewHeader("Big").Tall(true).Render()
	if !strings.Contains(out, "phui-header-tall") {
		t.Fatalf("missing tall class: %q", out)
	}
}

func TestHeaderNoBackground(t *testing.T) {
	out := NewHeader("Clean").NoBackground(true).Render()
	if !strings.Contains(out, "phui-header-no-background") {
		t.Fatalf("missing no-background class: %q", out)
	}
}

func TestHeaderHref(t *testing.T) {
	out := NewHeader("Linked").Href("/detail").Render()
	if !strings.Contains(out, `<a href="/detail">Linked</a>`) {
		t.Fatalf("missing linked header: %q", out)
	}
}

func TestHeaderXSS(t *testing.T) {
	out := NewHeader(`<script>alert("xss")</script>`).
		Subheader(`<img onerror="hack()">`).
		Image(`" onload="hack()`).
		Render()
	if strings.Contains(out, "<script>") {
		t.Fatalf("unescaped script in text: %q", out)
	}
	if strings.Contains(out, "<img") {
		t.Fatalf("unescaped img tag in subheader: %q", out)
	}
	if strings.Contains(out, `" onload`) {
		t.Fatalf("unescaped onload in image URL: %q", out)
	}
}
