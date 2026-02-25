package phui

import (
	"strings"
	"testing"
)

func TestInfoSeverityClasses(t *testing.T) {
	for _, sev := range []InfoSeverity{InfoNoData, InfoWarning, InfoError, InfoSuccess} {
		got := NewInfoView(sev).Render()
		want := "phui-info-severity-" + string(sev)
		if !strings.Contains(got, want) {
			t.Errorf("severity %s: expected class %q in %q", sev, want, got)
		}
	}
}

func TestInfoTitleOnly(t *testing.T) {
	got := NewInfoView(InfoWarning).Title("heads up").Render()
	if !strings.Contains(got, `<p class="phui-info-view-head">heads up</p>`) {
		t.Errorf("title not rendered: %s", got)
	}
	if strings.Contains(got, "phui-info-view-body") {
		t.Error("body should not be present")
	}
}

func TestInfoBodyOnly(t *testing.T) {
	got := NewInfoView(InfoError).Body("something broke").Render()
	if !strings.Contains(got, `<p class="phui-info-view-body">something broke</p>`) {
		t.Errorf("body not rendered: %s", got)
	}
	if strings.Contains(got, "phui-info-view-head") {
		t.Error("title should not be present")
	}
}

func TestInfoTitleAndBody(t *testing.T) {
	got := NewInfoView(InfoSuccess).Title("Done").Body("All good").Render()
	if !strings.Contains(got, `phui-info-view-head`) {
		t.Error("missing title")
	}
	if !strings.Contains(got, `phui-info-view-body`) {
		t.Error("missing body")
	}
}

func TestInfoEmpty(t *testing.T) {
	got := NewInfoView(InfoNoData).Render()
	if !strings.HasPrefix(got, `<div class="phui-info-view`) {
		t.Errorf("expected outer div, got: %s", got)
	}
	if !strings.HasSuffix(got, `</div>`) {
		t.Errorf("expected closing div, got: %s", got)
	}
}

func TestInfoTitleEscaped(t *testing.T) {
	got := NewInfoView(InfoWarning).Title(`<script>alert("xss")</script>`).Render()
	if strings.Contains(got, "<script>") {
		t.Error("title was not escaped")
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("expected escaped title, got: %s", got)
	}
}

func TestInfoBodyRaw(t *testing.T) {
	html := `<strong>bold</strong>`
	got := NewInfoView(InfoError).Body(html).Render()
	if !strings.Contains(got, html) {
		t.Errorf("body HTML should be raw, got: %s", got)
	}
}
