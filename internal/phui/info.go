package phui

type InfoSeverity string

const (
	InfoNoData  InfoSeverity = "nodata"
	InfoWarning InfoSeverity = "warning"
	InfoError   InfoSeverity = "error"
	InfoSuccess InfoSeverity = "success"
)

type InfoView struct {
	severity InfoSeverity
	title    string
	body     string
}

func NewInfoView(severity InfoSeverity) *InfoView {
	return &InfoView{severity: severity}
}

func (v *InfoView) Title(t string) *InfoView {
	v.title = t
	return v
}

func (v *InfoView) Body(b string) *InfoView {
	v.body = b
	return v
}

func (v *InfoView) Render() string {
	cls := classes("phui-info-view", "phui-info-severity-"+string(v.severity))
	out := `<div class="` + cls + `">`
	if v.title != "" {
		out += `<p class="phui-info-view-head">` + esc(v.title) + `</p>`
	}
	if v.body != "" {
		out += `<p class="phui-info-view-body">` + v.body + `</p>`
	}
	out += `</div>`
	return out
}
