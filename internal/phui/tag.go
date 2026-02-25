package phui

// TagType controls the visual style of a Tag.
type TagType string

const (
	TagShade   TagType = "shade"
	TagOutline TagType = "outline"
	TagState   TagType = "state"
	TagObject  TagType = "object"
	TagPerson  TagType = "person"
)

// Tag renders a PHUI tag/badge component.
type Tag struct {
	name     string
	tagType  TagType
	color    string
	icon     *Icon
	href     string
	closed   bool
	slim     bool
	dotColor string
	border   string
}

func NewTag(name string) *Tag {
	return &Tag{name: name, tagType: TagShade}
}

func (t *Tag) Type(tt TagType) *Tag    { t.tagType = tt; return t }
func (t *Tag) Color(c string) *Tag     { t.color = c; return t }
func (t *Tag) SetIcon(icon *Icon) *Tag { t.icon = icon; return t }
func (t *Tag) Href(h string) *Tag      { t.href = h; return t }
func (t *Tag) Closed(b bool) *Tag      { t.closed = b; return t }
func (t *Tag) Slim(b bool) *Tag        { t.slim = b; return t }
func (t *Tag) DotColor(c string) *Tag  { t.dotColor = c; return t }
func (t *Tag) Border(b string) *Tag    { t.border = b; return t }

func (t *Tag) Render() string {
	// Root classes
	rootCls := []string{"phui-tag-view", "phui-tag-type-" + string(t.tagType)}

	switch t.tagType {
	case TagShade:
		rootCls = append(rootCls, "phui-tag-shade")
		if t.color != "" {
			rootCls = append(rootCls, "phui-tag-"+t.color)
		}
	case TagOutline, TagState:
		if t.color != "" {
			rootCls = append(rootCls, "phui-tag-"+t.color)
		}
	}

	if t.slim {
		rootCls = append(rootCls, "phui-tag-slim")
	}
	if t.border != "" {
		rootCls = append(rootCls, "phui-tag-"+t.border)
	}

	// Core classes
	coreCls := "phui-tag-core"
	if t.tagType == TagState || t.tagType == TagObject {
		if t.color != "" {
			coreCls = classes(coreCls, "phui-tag-color-"+t.color)
		}
	}

	// Build inner content
	inner := t.renderInner(coreCls)

	// Root element
	tag := "span"
	hrefAttr := ""
	if t.href != "" {
		tag = "a"
		hrefAttr = attr("href", t.href)
	}

	return "<" + tag + attr("class", classes(rootCls...)) + hrefAttr + ">" +
		inner +
		"</" + tag + ">"
}

func (t *Tag) renderInner(coreCls string) string {
	iconHTML := ""
	if t.icon != nil {
		iconHTML = t.icon.Render()
	}

	nameHTML := esc(t.name)

	if t.closed {
		return `<span class="phui-tag-core-closed">` +
			iconHTML +
			`<span` + attr("class", coreCls) + `>` + nameHTML + `</span>` +
			`</span>`
	}

	dotHTML := ""
	if t.dotColor != "" {
		dotHTML = `<span` + attr("class", classes("phui-tag-dot", "phui-tag-color-"+t.dotColor)) + `></span>`
	}

	return `<span` + attr("class", coreCls) + `>` +
		dotHTML + iconHTML + nameHTML +
		`</span>`
}
