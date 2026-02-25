package phui

import (
	"fmt"
	"strings"
)

// ObjectItemIcon pairs an icon with a text label for display in the col2 icon list.
type ObjectItemIcon struct {
	Icon  *Icon
	Label string
}

// HandleIcon is a small circular avatar shown in the attributes row.
type HandleIcon struct {
	ImageURI string
	Label    string
}

// ObjectItem renders a single item in a PHUI object item list.
type ObjectItem struct {
	header      string
	href        string
	objectName  string
	subhead     string
	barColor    string
	imageURI    string
	imageIcon   *Icon
	imageHref   string
	disabled    bool
	effect      string // "highlighted", "selected", "visited"
	icons       []ObjectItemIcon
	attributes  []string // can contain raw HTML
	bylines     []string
	handleIcons []HandleIcon
	sideColumn  string // raw HTML
	statusIcon  *Icon
	epoch       string
	actions     []string // raw HTML
}

func NewObjectItem(header string) *ObjectItem {
	return &ObjectItem{header: header}
}

func (o *ObjectItem) Href(v string) *ObjectItem        { o.href = v; return o }
func (o *ObjectItem) ObjectName(v string) *ObjectItem   { o.objectName = v; return o }
func (o *ObjectItem) Subhead(v string) *ObjectItem      { o.subhead = v; return o }
func (o *ObjectItem) BarColor(v string) *ObjectItem     { o.barColor = v; return o }
func (o *ObjectItem) ImageURI(v string) *ObjectItem     { o.imageURI = v; return o }
func (o *ObjectItem) ImageIcon(v *Icon) *ObjectItem     { o.imageIcon = v; return o }
func (o *ObjectItem) ImageHref(v string) *ObjectItem    { o.imageHref = v; return o }
func (o *ObjectItem) Disabled(v bool) *ObjectItem       { o.disabled = v; return o }
func (o *ObjectItem) Effect(v string) *ObjectItem       { o.effect = v; return o }
func (o *ObjectItem) SideColumn(v string) *ObjectItem   { o.sideColumn = v; return o }
func (o *ObjectItem) StatusIcon(v *Icon) *ObjectItem    { o.statusIcon = v; return o }
func (o *ObjectItem) Epoch(v string) *ObjectItem        { o.epoch = v; return o }
func (o *ObjectItem) AddIcon(icon *Icon, label string) *ObjectItem {
	o.icons = append(o.icons, ObjectItemIcon{Icon: icon, Label: label})
	return o
}
func (o *ObjectItem) AddAttribute(v string) *ObjectItem {
	o.attributes = append(o.attributes, v)
	return o
}
func (o *ObjectItem) AddByline(v string) *ObjectItem {
	o.bylines = append(o.bylines, v)
	return o
}
func (o *ObjectItem) AddHandleIcon(uri, label string) *ObjectItem {
	o.handleIcons = append(o.handleIcons, HandleIcon{ImageURI: uri, Label: label})
	return o
}
func (o *ObjectItem) AddAction(v string) *ObjectItem {
	o.actions = append(o.actions, v)
	return o
}

func (o *ObjectItem) Render() string {
	return o.RenderWithClass("")
}

func (o *ObjectItem) RenderWithClass(itemClass string) string {
	var b strings.Builder

	// --- li classes ---
	liCls := []string{"phui-oi"}
	if itemClass != "" {
		liCls = append(liCls, itemClass)
	}
	if o.barColor != "" {
		liCls = append(liCls, "phui-oi-bar-color-"+o.barColor)
	} else {
		liCls = append(liCls, "phui-oi-no-bar")
	}
	if o.disabled {
		liCls = append(liCls, "phui-oi-disabled")
	} else {
		liCls = append(liCls, "phui-oi-enabled")
	}
	if o.effect != "" {
		liCls = append(liCls, "phui-oi-"+o.effect)
	}
	if o.imageURI != "" {
		liCls = append(liCls, "phui-oi-with-image")
	}
	if o.imageIcon != nil {
		liCls = append(liCls, "phui-oi-with-image-icon")
	}
	if len(o.icons) > 0 {
		liCls = append(liCls, "phui-oi-with-icons")
	}
	if len(o.attributes) > 0 {
		liCls = append(liCls, "phui-oi-with-attrs")
	}
	if len(o.handleIcons) > 0 {
		liCls = append(liCls, "phui-oi-with-handle-icons")
	}
	if len(o.actions) > 0 {
		liCls = append(liCls, "phui-oi-with-actions",
			fmt.Sprintf("phui-oi-with-%d-actions", len(o.actions)))
	}

	b.WriteString(`<li` + attr("class", classes(liCls...)) + `>`)
	b.WriteString(`<div class="phui-oi-frame">`)
	b.WriteString(`<div class="phui-oi-frame-content">`)

	// actions
	if len(o.actions) > 0 {
		b.WriteString(`<ul class="phui-oi-actions">`)
		for _, a := range o.actions {
			b.WriteString(`<li>`)
			b.WriteString(a)
			b.WriteString(`</li>`)
		}
		b.WriteString(`</ul>`)
	}

	// image
	if o.imageURI != "" {
		b.WriteString(`<div class="phui-oi-image" style="background-image:url(` + esc(o.imageURI) + `)"></div>`)
	}

	// image icon
	if o.imageIcon != nil {
		b.WriteString(`<div class="phui-oi-image-icon">` + o.imageIcon.Render() + `</div>`)
	}

	// content box
	b.WriteString(`<div class="phui-oi-content-box">`)
	b.WriteString(`<div class="phui-oi-table">`)
	b.WriteString(`<div class="phui-oi-table-row">`)

	// col0 — status icon
	if o.statusIcon != nil {
		b.WriteString(`<div class="phui-oi-col0"><div class="phui-oi-status-icon">` +
			o.statusIcon.Render() + `</div></div>`)
	}

	// col1
	b.WriteString(`<div class="phui-oi-col1">`)

	// name row
	b.WriteString(`<div class="phui-oi-name">`)
	if o.objectName != "" {
		b.WriteString(`<span class="phui-oi-objname">` + esc(o.objectName) + `</span>`)
	}
	if o.href != "" {
		b.WriteString(`<a class="phui-oi-link"` + attr("href", o.href) + `>` + esc(o.header) + `</a>`)
	} else {
		b.WriteString(`<div class="phui-oi-link">` + esc(o.header) + `</div>`)
	}
	b.WriteString(`</div>`) // phui-oi-name

	// content (subhead + attributes)
	b.WriteString(`<div class="phui-oi-content">`)
	if o.subhead != "" {
		b.WriteString(`<div class="phui-oi-subhead">` + esc(o.subhead) + `</div>`)
	}
	if len(o.handleIcons) > 0 || len(o.attributes) > 0 {
		b.WriteString(`<ul class="phui-oi-attributes">`)
		if len(o.handleIcons) > 0 {
			b.WriteString(`<li class="phui-oi-handle-icons">`)
			for _, h := range o.handleIcons {
				b.WriteString(`<span class="phui-oi-handle-icon" style="background-image:url(` + esc(h.ImageURI) + `)"></span>`)
			}
			b.WriteString(`</li>`)
		}
		for i, a := range o.attributes {
			b.WriteString(`<li class="phui-oi-attribute">`)
			if i > 0 {
				b.WriteString(`<span class="phui-oi-attribute-spacer">` + "\xc2\xb7" + `</span>`)
			}
			b.WriteString(a) // raw HTML
			b.WriteString(`</li>`)
		}
		b.WriteString(`</ul>`)
	}
	b.WriteString(`</div>`) // phui-oi-content
	b.WriteString(`</div>`) // phui-oi-col1

	// col2 — icons / bylines
	if len(o.icons) > 0 || len(o.bylines) > 0 {
		b.WriteString(`<div class="phui-oi-col2">`)
		if len(o.icons) > 0 {
			b.WriteString(`<ul class="phui-oi-icons">`)
			for _, ic := range o.icons {
				b.WriteString(`<li class="phui-oi-icon"><div class="phui-oi-icon-image">`)
				b.WriteString(ic.Icon.Render())
				b.WriteString(`</div>`)
				if ic.Label != "" {
					b.WriteString(`<span class="phui-oi-icon-label">` + esc(ic.Label) + `</span>`)
				}
				b.WriteString(`</li>`)
			}
			b.WriteString(`</ul>`)
		}
		if len(o.bylines) > 0 {
			b.WriteString(`<div class="phui-oi-bylines">`)
			for _, bl := range o.bylines {
				b.WriteString(`<div class="phui-oi-byline">` + esc(bl) + `</div>`)
			}
			b.WriteString(`</div>`)
		}
		b.WriteString(`</div>`) // phui-oi-col2
	}

	// side column
	if o.sideColumn != "" {
		b.WriteString(`<div class="phui-oi-col2 phui-oi-side-column">` + o.sideColumn + `</div>`)
	}

	b.WriteString(`</div>`) // phui-oi-table-row
	b.WriteString(`</div>`) // phui-oi-table
	b.WriteString(`</div>`) // phui-oi-content-box
	b.WriteString(`</div>`) // phui-oi-frame-content
	b.WriteString(`</div>`) // phui-oi-frame
	b.WriteString(`</li>`)

	return b.String()
}

// ObjectItemList renders a collection of ObjectItems.
type ObjectItemList struct {
	header     string
	items      []*ObjectItem
	itemClass  string
	flush      bool
	simple     bool
	noData     string
	allowEmpty bool
}

func NewObjectItemList() *ObjectItemList {
	return &ObjectItemList{
		itemClass: "phui-oi-standard",
		noData:    "No data.",
	}
}

func (l *ObjectItemList) Header(v string) *ObjectItemList     { l.header = v; return l }
func (l *ObjectItemList) AddItem(v *ObjectItem) *ObjectItemList { l.items = append(l.items, v); return l }
func (l *ObjectItemList) ItemClass(v string) *ObjectItemList   { l.itemClass = v; return l }
func (l *ObjectItemList) Flush(v bool) *ObjectItemList         { l.flush = v; return l }
func (l *ObjectItemList) Simple(v bool) *ObjectItemList        { l.simple = v; return l }
func (l *ObjectItemList) NoDataString(v string) *ObjectItemList { l.noData = v; return l }
func (l *ObjectItemList) AllowEmpty(v bool) *ObjectItemList    { l.allowEmpty = v; return l }

func (l *ObjectItemList) Render() string {
	var b strings.Builder

	ulCls := "phui-oi-list-view"
	if l.flush {
		ulCls += " phui-oi-list-flush"
	}
	if l.simple {
		ulCls += " phui-oi-list-simple"
	}

	b.WriteString(`<ul` + attr("class", ulCls) + `>`)

	if l.header != "" {
		b.WriteString(`<h1 class="phui-oi-list-header">` + esc(l.header) + `</h1>`)
	}

	if len(l.items) > 0 {
		for _, item := range l.items {
			b.WriteString(item.RenderWithClass(l.itemClass))
		}
	} else if !l.allowEmpty {
		b.WriteString(`<li class="phui-oi-empty">` + l.noData + `</li>`)
	}

	b.WriteString(`</ul>`)
	return b.String()
}
