package phui

import "testing"

func TestActionWithIconAndHref(t *testing.T) {
	html := NewAction("Edit").SetIcon(NewIcon("fa-pencil")).Href("/edit").Render()
	expect(t, html, `phabricator-action-view`)
	expect(t, html, `action-has-icon`)
	expect(t, html, `fa-pencil`)
	expect(t, html, `<a href="/edit"`)
	expect(t, html, `phabricator-action-view-item`)
	expect(t, html, `Edit`)
}

func TestActionWithoutIcon(t *testing.T) {
	html := NewAction("Plain").Href("/go").Render()
	// Icon span still rendered for alignment, but no FA class.
	expect(t, html, `phabricator-action-view-icon phui-icon-view phui-font-fa"`)
	expect(t, html, `<a href="/go"`)
}

func TestActionWithoutHref(t *testing.T) {
	html := NewAction("Label").SetIcon(NewIcon("fa-tag")).Render()
	// Should use <span> not <a>.
	expect(t, html, `<span class="phabricator-action-view-item"`)
	expectNot(t, html, `<a `)
}

func TestActionDisabled(t *testing.T) {
	html := NewAction("Nope").SetIcon(NewIcon("fa-ban")).Href("/no").Disabled(true).Render()
	expect(t, html, `phabricator-action-view-disabled`)
	expect(t, html, `<span`)
	expect(t, html, ` disabled`)
	// Disabled actions should not have <a>.
	expectNot(t, html, `<a `)
}

func TestActionListMultiple(t *testing.T) {
	list := NewActionList().
		AddAction(NewAction("One").Href("/1")).
		AddAction(NewAction("Two").Href("/2"))
	html := list.Render()
	expect(t, html, `phabricator-action-list-view`)
	expect(t, html, `One`)
	expect(t, html, `Two`)
	expect(t, html, `<ul`)
	expect(t, html, `</ul>`)
}

func TestActionListEmpty(t *testing.T) {
	html := NewActionList().Render()
	expect(t, html, `<ul class="phabricator-action-list-view"></ul>`)
}

func TestActionExtraClasses(t *testing.T) {
	html := NewAction("X").AddClass("foo").AddClass("bar").Render()
	expect(t, html, `foo`)
	expect(t, html, `bar`)
}

func TestActionXSS(t *testing.T) {
	html := NewAction(`<script>alert(1)</script>`).
		Href(`" onclick="alert(2)`).
		Render()
	expectNot(t, html, `<script>`)
	expect(t, html, `&lt;script&gt;`)
	expectNot(t, html, `" onclick=`)
}
