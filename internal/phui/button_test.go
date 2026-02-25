package phui

import (
	"strings"
	"testing"
)

func TestDefaultGreenButton(t *testing.T) {
	html := NewButton("Save").Color("green").Render()
	expect(t, html, `<button`)
	expect(t, html, `class="button button-green has-text"`)
	expect(t, html, `<div class="phui-button-text">Save</div>`)
	expect(t, html, `</button>`)
}

func TestLinkButton(t *testing.T) {
	html := NewButton("View").Tag("a").Href("/foo").Color("blue").Render()
	expect(t, html, `<a`)
	expect(t, html, `href="/foo"`)
	expect(t, html, `</a>`)
	expectNot(t, html, `</button>`)
}

func TestSubmitInput(t *testing.T) {
	html := NewButton("Go").Tag("input").Color("green").Name("action").Render()
	expect(t, html, `<input type="submit"`)
	expect(t, html, `value="Go"`)
	expect(t, html, `class="button button-green"`)
	expect(t, html, `name="action"`)
	expect(t, html, `/>`)
	expectNot(t, html, `</input>`)
}

func TestIconFirst(t *testing.T) {
	icon := NewIcon("fa-plus")
	html := NewButton("Add").SetIcon(icon).Render()
	expect(t, html, `has-icon`)
	iconIdx := strings.Index(html, `phui-icon-view`)
	textIdx := strings.Index(html, `phui-button-text`)
	if iconIdx > textIdx {
		t.Errorf("icon should appear before text, icon at %d, text at %d", iconIdx, textIdx)
	}
}

func TestIconLast(t *testing.T) {
	icon := NewIcon("fa-caret-down")
	html := NewButton("Menu").SetIcon(icon).IconFirst(false).Render()
	expect(t, html, `icon-last`)
	iconIdx := strings.Index(html, `phui-icon-view`)
	textIdx := strings.Index(html, `phui-button-text`)
	if iconIdx < textIdx {
		t.Errorf("icon should appear after text, icon at %d, text at %d", iconIdx, textIdx)
	}
}

func TestDisabledButton(t *testing.T) {
	html := NewButton("Nope").Disabled(true).Render()
	expect(t, html, `class="button has-text disabled"`)
	// disabled attribute on button tag
	if strings.Count(html, "disabled") < 2 {
		t.Error("expected disabled class AND disabled attribute")
	}
}

func TestDisabledInput(t *testing.T) {
	html := NewButton("Nope").Tag("input").Disabled(true).Render()
	expect(t, html, ` disabled`)
}

func TestDropdown(t *testing.T) {
	html := NewButton("Actions").Dropdown(true).Render()
	expect(t, html, `dropdown`)
	expect(t, html, `<span class="caret"></span>`)
}

func TestSelected(t *testing.T) {
	html := NewButton("Active").Selected(true).Render()
	expect(t, html, `selected`)
}

func TestSizeSmall(t *testing.T) {
	html := NewButton("Tiny").Size("small").Render()
	expect(t, html, `small`)
}

func TestSizeBig(t *testing.T) {
	html := NewButton("Large").Size("big").Render()
	expect(t, html, `big`)
}

func TestExtraClasses(t *testing.T) {
	html := NewButton("X").AddClass("custom-one").AddClass("custom-two").Render()
	expect(t, html, `custom-one`)
	expect(t, html, `custom-two`)
}

func TestXSSText(t *testing.T) {
	html := NewButton(`<script>alert(1)</script>`).Render()
	expectNot(t, html, `<script>`)
	expect(t, html, `&lt;script&gt;`)
}

func TestXSSHref(t *testing.T) {
	html := NewButton("Click").Tag("a").Href(`" onclick="alert(1)`).Render()
	// The double quote must be escaped so the attribute can't break out
	expect(t, html, `&#34;`)
	expectNot(t, html, `" onclick=`)
}

func expect(t *testing.T, html, sub string) {
	t.Helper()
	if !strings.Contains(html, sub) {
		t.Errorf("expected %q in:\n%s", sub, html)
	}
}

func expectNot(t *testing.T, html, sub string) {
	t.Helper()
	if strings.Contains(html, sub) {
		t.Errorf("did not expect %q in:\n%s", sub, html)
	}
}
