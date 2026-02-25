package phui

// ObjectBox renders a PHUI object box container.
type ObjectBox struct {
	header     *Header
	headerText string
	body       string
	color      string
	flush      bool
	extra      []string
}

func NewObjectBox() *ObjectBox                    { return &ObjectBox{} }
func (b *ObjectBox) SetHeader(h *Header) *ObjectBox { b.header = h; return b }
func (b *ObjectBox) HeaderText(t string) *ObjectBox { b.headerText = t; return b }
func (b *ObjectBox) Body(html string) *ObjectBox    { b.body = html; return b }
func (b *ObjectBox) Color(c string) *ObjectBox      { b.color = c; return b }
func (b *ObjectBox) Flush(f bool) *ObjectBox        { b.flush = f; return b }
func (b *ObjectBox) AddClass(c string) *ObjectBox   { b.extra = append(b.extra, c); return b }

func (b *ObjectBox) Render() string {
	cls := "phui-box phui-box-border phui-object-box mlt mlr"
	if b.color != "" {
		cls = classes(cls, "phui-object-box-"+b.color)
	}
	if b.flush {
		cls = classes(cls, "phui-object-box-flush")
	}
	for _, c := range b.extra {
		cls = classes(cls, c)
	}

	s := `<div` + attr("class", cls) + `>`
	if b.header != nil {
		s += b.header.Render()
	} else if b.headerText != "" {
		s += NewHeader(b.headerText).Render()
	}
	s += b.body
	s += `</div>`
	return s
}
