package phui

// Button renders a PHUI button (button, a, or input tag).
type Button struct {
	text      string
	tag       string
	color     string
	size      string
	href      string
	icon      *Icon
	iconFirst bool
	disabled  bool
	selected  bool
	dropdown  bool
	name      string
	btnType   string
	extra     []string
}

// NewButton creates a button with the given text.
func NewButton(text string) *Button {
	return &Button{text: text, tag: "button", iconFirst: true}
}

func (b *Button) Tag(t string) *Button       { b.tag = t; return b }
func (b *Button) Color(c string) *Button     { b.color = c; return b }
func (b *Button) Size(s string) *Button      { b.size = s; return b }
func (b *Button) Href(h string) *Button      { b.href = h; return b }
func (b *Button) SetIcon(i *Icon) *Button    { b.icon = i; return b }
func (b *Button) IconFirst(f bool) *Button   { b.iconFirst = f; return b }
func (b *Button) Disabled(d bool) *Button    { b.disabled = d; return b }
func (b *Button) Selected(s bool) *Button    { b.selected = s; return b }
func (b *Button) Dropdown(d bool) *Button    { b.dropdown = d; return b }
func (b *Button) Name(n string) *Button      { b.name = n; return b }
func (b *Button) Type(t string) *Button      { b.btnType = t; return b }
func (b *Button) AddClass(c string) *Button  { b.extra = append(b.extra, c); return b }

// Render produces the HTML for this button.
func (b *Button) Render() string {
	if b.tag == "input" {
		return b.renderInput()
	}
	return b.renderTag()
}

func (b *Button) renderInput() string {
	cls := b.baseClasses()
	s := `<input type="submit"` + attr("value", esc(b.text)) + attr("class", cls)
	s += attr("name", b.name)
	if b.disabled {
		s += ` disabled`
	}
	s += ` />`
	return s
}

func (b *Button) renderTag() string {
	cls := b.baseClasses()
	if b.icon != nil {
		cls = classes(cls, "has-icon")
	}
	if b.text != "" {
		cls = classes(cls, "has-text")
	}
	if b.icon != nil && !b.iconFirst {
		cls = classes(cls, "icon-last")
	}
	if b.disabled {
		cls = classes(cls, "disabled")
	}
	if b.selected {
		cls = classes(cls, "selected")
	}
	if b.dropdown {
		cls = classes(cls, "dropdown")
	}

	tag := b.tag
	s := `<` + tag + attr("class", cls)
	if tag == "a" && b.href != "" {
		s += attr("href", b.href)
	}
	s += attr("name", b.name)
	s += attr("type", b.btnType)
	if b.disabled && tag == "button" {
		s += ` disabled`
	}
	s += `>`

	if b.iconFirst && b.icon != nil {
		s += b.icon.Render()
	}
	s += `<div class="phui-button-text">` + esc(b.text) + `</div>`
	if !b.iconFirst && b.icon != nil {
		s += b.icon.Render()
	}
	if b.dropdown {
		s += `<span class="caret"></span>`
	}

	s += `</` + tag + `>`
	return s
}

func (b *Button) baseClasses() string {
	cc := []string{"button"}
	if b.color != "" {
		cc = append(cc, "button-"+b.color)
	}
	if b.size != "" {
		cc = append(cc, b.size)
	}
	for _, e := range b.extra {
		cc = append(cc, e)
	}
	return classes(cc...)
}
