package phui

// Header renders a PHUI header component.
type Header struct {
	text         string
	subheader    string
	icon         *Icon
	image        string
	imageHref    string
	tags         []*Tag
	actionLinks  []*Button
	href         string
	tall         bool
	noBackground bool
}

func NewHeader(text string) *Header              { return &Header{text: text} }
func (h *Header) Subheader(s string) *Header     { h.subheader = s; return h }
func (h *Header) SetIcon(i *Icon) *Header        { h.icon = i; return h }
func (h *Header) Image(url string) *Header       { h.image = url; return h }
func (h *Header) ImageHref(url string) *Header   { h.imageHref = url; return h }
func (h *Header) AddTag(t *Tag) *Header          { h.tags = append(h.tags, t); return h }
func (h *Header) AddActionLink(b *Button) *Header { h.actionLinks = append(h.actionLinks, b); return h }
func (h *Header) Href(url string) *Header        { h.href = url; return h }
func (h *Header) Tall(b bool) *Header            { h.tall = b; return h }
func (h *Header) NoBackground(b bool) *Header    { h.noBackground = b; return h }

func (h *Header) Render() string {
	shellCls := "phui-header-shell"
	if h.tall {
		shellCls = classes(shellCls, "phui-header-tall")
	}
	if h.noBackground {
		shellCls = classes(shellCls, "phui-header-no-background")
	}

	s := `<div` + attr("class", shellCls) + `>`
	s += `<div class="phui-header-row">`

	// Col1: image/avatar
	s += `<div class="phui-header-col1">`
	if h.image != "" {
		imgSpan := `<span class="phui-header-image" style="background-image:url(` + esc(h.image) + `)"></span>`
		if h.imageHref != "" {
			imgSpan = `<a` + attr("href", h.imageHref) + `>` + imgSpan + `</a>`
		}
		s += imgSpan
	}
	s += `</div>`

	// Col2: header text, icon, tags, subheader
	s += `<div class="phui-header-col2">`
	s += `<span class="phui-header-header">`
	if h.icon != nil {
		s += `<span class="phui-header-icon">` + h.icon.Render() + `</span>`
	}
	if h.href != "" {
		s += `<a` + attr("href", h.href) + `>` + esc(h.text) + `</a>`
	} else {
		s += esc(h.text)
	}
	s += `</span>`
	for _, tag := range h.tags {
		s += tag.Render()
	}
	if h.subheader != "" {
		s += `<div class="phui-header-subheader">` + esc(h.subheader) + `</div>`
	}
	s += `</div>`

	// Col3: action links
	s += `<div class="phui-header-col3">`
	if len(h.actionLinks) > 0 {
		s += `<div class="phui-header-action-links">`
		for _, btn := range h.actionLinks {
			s += btn.Render()
		}
		s += `</div>`
	}
	s += `</div>`

	s += `</div></div>`
	return s
}
