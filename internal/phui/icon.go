package phui

// Icon renders a PHUI icon span (FontAwesome-based).
type Icon struct {
	name  string
	extra []string
	color string
}

// NewIcon creates an icon with the given FA class name.
func NewIcon(name string) *Icon {
	return &Icon{name: name}
}

// AddClass appends an extra CSS class.
func (i *Icon) AddClass(c string) *Icon {
	i.extra = append(i.extra, c)
	return i
}

// Color sets an inline color style.
func (i *Icon) Color(c string) *Icon {
	i.color = c
	return i
}

// Render produces the HTML span.
func (i *Icon) Render() string {
	cls := classes("phui-icon-view", "phui-font-fa", esc(i.name))
	for _, e := range i.extra {
		cls = classes(cls, esc(e))
	}
	style := ""
	if i.color != "" {
		style = ` style="color:` + esc(i.color) + `"`
	}
	return `<span` + attr("class", cls) + style + `></span>`
}
