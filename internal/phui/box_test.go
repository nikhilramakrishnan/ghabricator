package phui

import (
	"strings"
	"testing"
)

func TestObjectBoxBodyOnly(t *testing.T) {
	html := NewObjectBox().Body("<p>hello</p>").Render()
	if !strings.Contains(html, "phui-object-box") {
		t.Fatal("missing phui-object-box class")
	}
	if !strings.Contains(html, "<p>hello</p>") {
		t.Fatal("missing body content")
	}
}

func TestObjectBoxWithHeader(t *testing.T) {
	h := NewHeader("My Title")
	html := NewObjectBox().SetHeader(h).Body("<p>body</p>").Render()
	if !strings.Contains(html, "phui-header-shell") {
		t.Fatal("missing header render")
	}
	if !strings.Contains(html, "My Title") {
		t.Fatal("missing header text")
	}
}

func TestObjectBoxWithHeaderText(t *testing.T) {
	html := NewObjectBox().HeaderText("Auto Header").Body("<p>body</p>").Render()
	if !strings.Contains(html, "phui-header-shell") {
		t.Fatal("headerText should auto-create header")
	}
	if !strings.Contains(html, "Auto Header") {
		t.Fatal("missing headerText content")
	}
}

func TestObjectBoxHeaderPrecedence(t *testing.T) {
	h := NewHeader("Explicit")
	html := NewObjectBox().SetHeader(h).HeaderText("Fallback").Body("x").Render()
	if !strings.Contains(html, "Explicit") {
		t.Fatal("explicit header should be present")
	}
	if strings.Contains(html, "Fallback") {
		t.Fatal("headerText should not render when header is set")
	}
}

func TestObjectBoxColors(t *testing.T) {
	for _, color := range []string{"red", "blue", "green", "yellow"} {
		html := NewObjectBox().Color(color).Render()
		want := "phui-object-box-" + color
		if !strings.Contains(html, want) {
			t.Fatalf("color %q: missing class %q", color, want)
		}
	}
}

func TestObjectBoxFlush(t *testing.T) {
	html := NewObjectBox().Flush(true).Render()
	if !strings.Contains(html, "phui-object-box-flush") {
		t.Fatal("missing flush class")
	}
}

func TestObjectBoxExtraClasses(t *testing.T) {
	html := NewObjectBox().AddClass("my-custom").AddClass("another").Render()
	if !strings.Contains(html, "my-custom") {
		t.Fatal("missing custom class")
	}
	if !strings.Contains(html, "another") {
		t.Fatal("missing second custom class")
	}
}

func TestObjectBoxRawHTML(t *testing.T) {
	html := NewObjectBox().Body(`<div class="inner"><b>bold</b></div>`).Render()
	if !strings.Contains(html, `<div class="inner"><b>bold</b></div>`) {
		t.Fatal("body HTML should not be escaped")
	}
}

func TestObjectBoxEmpty(t *testing.T) {
	html := NewObjectBox().Render()
	if !strings.Contains(html, "phui-object-box") {
		t.Fatal("empty box should still have base classes")
	}
	if strings.Contains(html, "phui-header-shell") {
		t.Fatal("empty box should not have a header")
	}
}

func TestObjectBoxHeaderTextEscaping(t *testing.T) {
	html := NewObjectBox().HeaderText(`<script>alert("xss")</script>`).Render()
	if strings.Contains(html, "<script>") {
		t.Fatal("headerText should be escaped by Header")
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Fatal("headerText should contain escaped script tag")
	}
}
