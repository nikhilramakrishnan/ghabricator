package phui

import (
	"strings"
	"testing"
)

// --- ObjectItemList tests ---

func TestEmptyListDefaultNoData(t *testing.T) {
	out := NewObjectItemList().Render()
	if !strings.Contains(out, `phui-oi-empty`) {
		t.Fatalf("expected empty class in %q", out)
	}
	if !strings.Contains(out, "No data.") {
		t.Fatalf("expected default no-data message in %q", out)
	}
}

func TestEmptyListAllowEmpty(t *testing.T) {
	out := NewObjectItemList().AllowEmpty(true).Render()
	if strings.Contains(out, `phui-oi-empty`) {
		t.Fatalf("should not have empty item when allowEmpty: %q", out)
	}
	if strings.Contains(out, "No data.") {
		t.Fatalf("should not have no-data message when allowEmpty: %q", out)
	}
}

func TestListWithItems(t *testing.T) {
	out := NewObjectItemList().
		AddItem(NewObjectItem("Alpha")).
		AddItem(NewObjectItem("Beta")).
		Render()
	if !strings.Contains(out, "Alpha") || !strings.Contains(out, "Beta") {
		t.Fatalf("expected both items in %q", out)
	}
	if strings.Contains(out, `phui-oi-empty`) {
		t.Fatalf("should not have empty item: %q", out)
	}
	// items get default itemClass
	if c := strings.Count(out, "phui-oi-standard"); c != 2 {
		t.Fatalf("expected 2 standard classes, got %d in %q", c, out)
	}
}

func TestFlushList(t *testing.T) {
	out := NewObjectItemList().Flush(true).AllowEmpty(true).Render()
	if !strings.Contains(out, "phui-oi-list-flush") {
		t.Fatalf("expected flush class in %q", out)
	}
}

func TestSimpleList(t *testing.T) {
	out := NewObjectItemList().Simple(true).AllowEmpty(true).Render()
	if !strings.Contains(out, "phui-oi-list-simple") {
		t.Fatalf("expected simple class in %q", out)
	}
}

func TestListHeader(t *testing.T) {
	out := NewObjectItemList().Header("Tasks").AllowEmpty(true).Render()
	if !strings.Contains(out, `<h1 class="phui-oi-list-header">Tasks</h1>`) {
		t.Fatalf("expected list header in %q", out)
	}
}

func TestCustomNoData(t *testing.T) {
	out := NewObjectItemList().NoDataString("Nothing here.").Render()
	if !strings.Contains(out, "Nothing here.") {
		t.Fatalf("expected custom no-data in %q", out)
	}
}

// --- ObjectItem tests ---

func TestItemHeaderOnly(t *testing.T) {
	out := NewObjectItem("My Task").Render()
	if !strings.Contains(out, `<div class="phui-oi-link">My Task</div>`) {
		t.Fatalf("expected header as div when no href: %q", out)
	}
	if !strings.Contains(out, "phui-oi-no-bar") {
		t.Fatalf("expected no-bar class: %q", out)
	}
	if !strings.Contains(out, "phui-oi-enabled") {
		t.Fatalf("expected enabled class: %q", out)
	}
}

func TestItemHeaderHref(t *testing.T) {
	out := NewObjectItem("My Task").Href("/T123").Render()
	if !strings.Contains(out, `<a class="phui-oi-link" href="/T123">My Task</a>`) {
		t.Fatalf("expected link in %q", out)
	}
}

func TestItemBarColor(t *testing.T) {
	out := NewObjectItem("X").BarColor("green").Render()
	if !strings.Contains(out, "phui-oi-bar-color-green") {
		t.Fatalf("expected bar color class in %q", out)
	}
	if strings.Contains(out, "phui-oi-no-bar") {
		t.Fatalf("should not have no-bar when bar color set: %q", out)
	}
}

func TestItemNoBar(t *testing.T) {
	out := NewObjectItem("X").Render()
	if !strings.Contains(out, "phui-oi-no-bar") {
		t.Fatalf("expected no-bar class in %q", out)
	}
}

func TestItemImageURI(t *testing.T) {
	out := NewObjectItem("X").ImageURI("/img/avatar.png").Render()
	if !strings.Contains(out, "phui-oi-with-image") {
		t.Fatalf("expected with-image class in %q", out)
	}
	if !strings.Contains(out, `background-image:url(/img/avatar.png)`) {
		t.Fatalf("expected background-image in %q", out)
	}
}

func TestItemImageIcon(t *testing.T) {
	out := NewObjectItem("X").ImageIcon(NewIcon("fa-user")).Render()
	if !strings.Contains(out, "phui-oi-with-image-icon") {
		t.Fatalf("expected with-image-icon class in %q", out)
	}
	if !strings.Contains(out, `phui-oi-image-icon`) {
		t.Fatalf("expected image-icon div in %q", out)
	}
	if !strings.Contains(out, "fa-user") {
		t.Fatalf("expected icon name in %q", out)
	}
}

func TestItemAttributes(t *testing.T) {
	out := NewObjectItem("X").
		AddAttribute("High").
		AddAttribute("Open").
		AddAttribute("Bug").
		Render()
	if !strings.Contains(out, "phui-oi-with-attrs") {
		t.Fatalf("expected with-attrs class in %q", out)
	}
	// first attribute: no spacer
	if !strings.Contains(out, `<li class="phui-oi-attribute">High</li>`) {
		t.Fatalf("expected first attr without spacer in %q", out)
	}
	// second attribute: has spacer
	if !strings.Contains(out, `<span class="phui-oi-attribute-spacer">`) {
		t.Fatalf("expected dot spacer in %q", out)
	}
	// count spacers: should be 2 for 3 attributes
	if c := strings.Count(out, "phui-oi-attribute-spacer"); c != 2 {
		t.Fatalf("expected 2 spacers, got %d in %q", c, out)
	}
}

func TestItemHandleIcons(t *testing.T) {
	out := NewObjectItem("X").
		AddHandleIcon("/img/a.png", "Alice").
		AddHandleIcon("/img/b.png", "Bob").
		Render()
	if !strings.Contains(out, "phui-oi-with-handle-icons") {
		t.Fatalf("expected with-handle-icons class in %q", out)
	}
	if !strings.Contains(out, `phui-oi-handle-icon`) {
		t.Fatalf("expected handle-icon span in %q", out)
	}
	if c := strings.Count(out, "phui-oi-handle-icon"); c != 3 { // 1 li + 2 spans
		t.Fatalf("expected 3 handle-icon occurrences (1 li class + 2 spans), got %d in %q", c, out)
	}
}

func TestItemIconsWithLabels(t *testing.T) {
	out := NewObjectItem("X").
		AddIcon(NewIcon("fa-calendar"), "Jan 1").
		AddIcon(NewIcon("fa-user"), "").
		Render()
	if !strings.Contains(out, "phui-oi-with-icons") {
		t.Fatalf("expected with-icons class in %q", out)
	}
	if !strings.Contains(out, `<span class="phui-oi-icon-label">Jan 1</span>`) {
		t.Fatalf("expected icon label in %q", out)
	}
	// second icon has no label
	if c := strings.Count(out, "phui-oi-icon-label"); c != 1 {
		t.Fatalf("expected 1 icon label, got %d in %q", c, out)
	}
}

func TestItemBylines(t *testing.T) {
	out := NewObjectItem("X").
		AddByline("Created by Alice").
		AddByline("Edited by Bob").
		Render()
	if !strings.Contains(out, `phui-oi-col2`) {
		t.Fatalf("expected col2 in %q", out)
	}
	if !strings.Contains(out, `phui-oi-bylines`) {
		t.Fatalf("expected bylines div in %q", out)
	}
	if c := strings.Count(out, `class="phui-oi-byline"`); c != 2 {
		t.Fatalf("expected 2 byline divs, got %d in %q", c, out)
	}
}

func TestItemSideColumn(t *testing.T) {
	out := NewObjectItem("X").SideColumn(`<div class="custom">stuff</div>`).Render()
	if !strings.Contains(out, `phui-oi-side-column`) {
		t.Fatalf("expected side-column class in %q", out)
	}
	if !strings.Contains(out, `<div class="custom">stuff</div>`) {
		t.Fatalf("expected raw side column HTML in %q", out)
	}
}

func TestItemStatusIcon(t *testing.T) {
	out := NewObjectItem("X").StatusIcon(NewIcon("fa-check").Color("green")).Render()
	if !strings.Contains(out, `phui-oi-col0`) {
		t.Fatalf("expected col0 in %q", out)
	}
	if !strings.Contains(out, `phui-oi-status-icon`) {
		t.Fatalf("expected status-icon div in %q", out)
	}
	if !strings.Contains(out, "fa-check") {
		t.Fatalf("expected icon in %q", out)
	}
}

func TestItemDisabled(t *testing.T) {
	out := NewObjectItem("X").Disabled(true).Render()
	if !strings.Contains(out, "phui-oi-disabled") {
		t.Fatalf("expected disabled class in %q", out)
	}
	if strings.Contains(out, "phui-oi-enabled") {
		t.Fatalf("should not have enabled when disabled: %q", out)
	}
}

func TestItemEffects(t *testing.T) {
	for _, effect := range []string{"highlighted", "selected", "visited"} {
		out := NewObjectItem("X").Effect(effect).Render()
		want := "phui-oi-" + effect
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q class in %q", want, out)
		}
	}
}

func TestItemObjectName(t *testing.T) {
	out := NewObjectItem("Fix the thing").ObjectName("T123").Href("/T123").Render()
	if !strings.Contains(out, `<span class="phui-oi-objname">T123</span>`) {
		t.Fatalf("expected objname span in %q", out)
	}
	// objname appears before the link
	nameIdx := strings.Index(out, "phui-oi-objname")
	linkIdx := strings.Index(out, "phui-oi-link")
	if nameIdx > linkIdx {
		t.Fatalf("objname should come before link in %q", out)
	}
}

func TestItemActions(t *testing.T) {
	out := NewObjectItem("X").
		AddAction(`<button>Edit</button>`).
		AddAction(`<button>Delete</button>`).
		Render()
	if !strings.Contains(out, "phui-oi-with-actions") {
		t.Fatalf("expected with-actions class in %q", out)
	}
	if !strings.Contains(out, "phui-oi-with-2-actions") {
		t.Fatalf("expected with-2-actions class in %q", out)
	}
	if !strings.Contains(out, `<li><button>Edit</button></li>`) {
		t.Fatalf("expected action in list in %q", out)
	}
}

func TestItemSubhead(t *testing.T) {
	out := NewObjectItem("X").Subhead("A subtitle").Render()
	if !strings.Contains(out, `<div class="phui-oi-subhead">A subtitle</div>`) {
		t.Fatalf("expected subhead div in %q", out)
	}
}

// --- XSS tests ---

func TestObjectItemXSSHeader(t *testing.T) {
	out := NewObjectItem(`<script>alert(1)</script>`).Render()
	if strings.Contains(out, "<script>") {
		t.Fatalf("unescaped script in header: %q", out)
	}
}

func TestObjectItemXSSHref(t *testing.T) {
	out := NewObjectItem("X").Href(`"><script>alert(1)</script>`).Render()
	if strings.Contains(out, "<script>") {
		t.Fatalf("unescaped script in href: %q", out)
	}
}

func TestObjectItemXSSImageURI(t *testing.T) {
	out := NewObjectItem("X").ImageURI(`"><script>alert(1)</script>`).Render()
	if strings.Contains(out, "<script>") {
		t.Fatalf("unescaped script in imageURI: %q", out)
	}
}

func TestObjectItemXSSObjectName(t *testing.T) {
	out := NewObjectItem("X").ObjectName(`<img onerror=alert(1)>`).Render()
	if strings.Contains(out, "<img") {
		t.Fatalf("unescaped tag in objectName: %q", out)
	}
}

func TestObjectItemXSSSubhead(t *testing.T) {
	out := NewObjectItem("X").Subhead(`<script>x</script>`).Render()
	if strings.Contains(out, "<script>") {
		t.Fatalf("unescaped script in subhead: %q", out)
	}
}

func TestObjectItemXSSByline(t *testing.T) {
	out := NewObjectItem("X").AddByline(`<script>x</script>`).Render()
	if strings.Contains(out, "<script>") {
		t.Fatalf("unescaped script in byline: %q", out)
	}
}

func TestObjectItemXSSHandleIconURI(t *testing.T) {
	out := NewObjectItem("X").AddHandleIcon(`"><script>`, "Hax").Render()
	if strings.Contains(out, "<script>") {
		t.Fatalf("unescaped script in handle icon URI: %q", out)
	}
}

func TestObjectItemXSSListHeader(t *testing.T) {
	out := NewObjectItemList().Header(`<script>x</script>`).AllowEmpty(true).Render()
	if strings.Contains(out, "<script>") {
		t.Fatalf("unescaped script in list header: %q", out)
	}
}

// --- RenderWithClass ---

func TestRenderWithClass(t *testing.T) {
	out := NewObjectItem("X").RenderWithClass("my-custom-class")
	if !strings.Contains(out, "my-custom-class") {
		t.Fatalf("expected custom class in %q", out)
	}
}
