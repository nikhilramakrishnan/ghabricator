package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/nikhilr/ghabricator/internal/auth"
	"github.com/nikhilr/ghabricator/internal/herald"
	"github.com/nikhilr/ghabricator/internal/templates"
)

// handleHeraldList renders the Herald rules list page.
func (s *Server) handleHeraldList(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	theme := templates.ThemeFromRequest(r)

	rules, err := s.herald.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("load rules: %v", err), http.StatusInternalServerError)
		return
	}

	var buf strings.Builder
	buf.WriteString(`<div class="phui-object-box">`)
	buf.WriteString(`<div class="phui-header-view">`)
	buf.WriteString(`<h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa fa-bullhorn"></span>Herald Rules</h1>`)
	buf.WriteString(`<div class="phui-header-action-links">`)
	buf.WriteString(`<a href="/herald/new" class="phui-button-view button-green" style="padding:4px 12px;border-radius:3px;background:#139543;color:#fff;text-decoration:none;font-size:13px;">`)
	buf.WriteString(`<span class="phui-icon-view phui-font-fa fa-plus" style="margin-right:4px;"></span>New Rule</a>`)
	buf.WriteString(`</div></div>`)

	if len(rules) == 0 {
		buf.WriteString(`<div class="phui-info-view phui-info-severity-nodata">No Herald rules have been created yet.</div>`)
	} else {
		buf.WriteString(`<ul class="phui-oi-list-view">`)
		for _, rule := range rules {
			statusClass := ""
			if rule.Disabled {
				statusClass = " phui-oi-status-draft"
			}
			buf.WriteString(fmt.Sprintf(`<li class="phui-oi phui-oi-standard%s">`, statusClass))
			buf.WriteString(`<div class="phui-oi-frame">`)

			// Icon
			buf.WriteString(`<div class="phui-oi-image-icon">`)
			buf.WriteString(`<span class="phui-icon-view phui-font-fa fa-bullhorn" style="font-size:16px;color:#6b748c;"></span>`)
			buf.WriteString(`</div>`)

			// Content
			buf.WriteString(`<div class="phui-oi-content-box">`)
			buf.WriteString(`<div class="phui-oi-name">`)
			buf.WriteString(fmt.Sprintf(`<a href="/herald/%s">%s</a>`,
				template.HTMLEscapeString(rule.ID),
				template.HTMLEscapeString(rule.Name)))
			buf.WriteString(`</div>`)

			// Attributes
			buf.WriteString(`<div class="phui-oi-attributes">`)
			matchMode := "all"
			if !rule.MustMatchAll {
				matchMode = "any"
			}
			buf.WriteString(fmt.Sprintf(`<span class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-filter"></span> %d conditions (%s)</span>`,
				len(rule.Conditions), matchMode))
			buf.WriteString(fmt.Sprintf(`<span class="phui-oi-attribute"><span class="phui-icon-view phui-font-fa fa-bolt"></span> %d actions</span>`,
				len(rule.Actions)))
			if rule.Disabled {
				buf.WriteString(`<span class="phui-oi-attribute"><span class="phui-tag-view phui-tag-shade-grey"><span class="phui-tag-core">Disabled</span></span></span>`)
			}
			buf.WriteString(`</div>`)

			buf.WriteString(`</div>`) // content-box
			buf.WriteString(`</div>`) // frame
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul>`)
	}
	buf.WriteString(`</div>`) // object-box

	templates.RenderPage(w, templates.PageData{
		Title:         "Herald Rules",
		Theme:         theme,
		HeaderTitle:   template.HTML("Herald"),
		HeaderIcon:    "fa-bullhorn",
		Content:       template.HTML(buf.String()),
		NavActive:     "herald",
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs: []templates.Crumb{
			{Name: "Home", Href: "/"},
			{Name: "Herald"},
		},
	})
}

// handleHeraldNew renders the create-rule form.
func (s *Server) handleHeraldNew(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	theme := templates.ThemeFromRequest(r)

	content := renderHeraldForm(nil)

	templates.RenderPage(w, templates.PageData{
		Title:         "New Herald Rule",
		Theme:         theme,
		HeaderTitle:   template.HTML("New Herald Rule"),
		HeaderIcon:    "fa-bullhorn",
		Content:       template.HTML(content),
		NavActive:     "herald",
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs: []templates.Crumb{
			{Name: "Home", Href: "/"},
			{Name: "Herald", Href: "/herald"},
			{Name: "New Rule"},
		},
	})
}

// handleHeraldView renders a single rule for viewing/editing.
func (s *Server) handleHeraldView(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	theme := templates.ThemeFromRequest(r)
	id := r.PathValue("id")

	rule, err := s.herald.Get(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("load rule: %v", err), http.StatusInternalServerError)
		return
	}
	if rule == nil {
		http.Error(w, "rule not found", http.StatusNotFound)
		return
	}

	content := renderHeraldForm(rule)

	templates.RenderPage(w, templates.PageData{
		Title:         "Edit Herald Rule",
		Theme:         theme,
		HeaderTitle:   template.HTML(template.HTMLEscapeString(rule.Name)),
		HeaderIcon:    "fa-bullhorn",
		Content:       template.HTML(content),
		NavActive:     "herald",
		UserLogin:     sess.Login,
		UserAvatarURL: sess.AvatarURL,
		Crumbs: []templates.Crumb{
			{Name: "Home", Href: "/"},
			{Name: "Herald", Href: "/herald"},
			{Name: rule.Name},
		},
	})
}

// handleHeraldSave processes the create/edit form submission.
func (s *Server) handleHeraldSave(w http.ResponseWriter, r *http.Request) {
	sess := auth.SessionFromContext(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	rule := herald.Rule{
		ID:           r.FormValue("id"),
		Name:         r.FormValue("name"),
		AuthorLogin:  sess.Login,
		MustMatchAll: r.FormValue("match_mode") == "all",
		Disabled:     r.FormValue("disabled") == "1",
	}

	// Parse conditions
	condTypes := r.Form["cond_type"]
	condValues := r.Form["cond_value"]
	for i := range condTypes {
		if i < len(condValues) && condTypes[i] != "" && condValues[i] != "" {
			rule.Conditions = append(rule.Conditions, herald.Condition{
				Type:  herald.ConditionType(condTypes[i]),
				Value: condValues[i],
			})
		}
	}

	// Parse actions
	actTypes := r.Form["action_type"]
	actValues := r.Form["action_value"]
	for i := range actTypes {
		if i < len(actValues) && actTypes[i] != "" && actValues[i] != "" {
			rule.Actions = append(rule.Actions, herald.Action{
				Type:  herald.ActionType(actTypes[i]),
				Value: actValues[i],
			})
		}
	}

	if rule.Name == "" {
		http.Error(w, "rule name is required", http.StatusBadRequest)
		return
	}

	if err := s.herald.Save(&rule); err != nil {
		http.Error(w, fmt.Sprintf("save rule: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/herald", http.StatusSeeOther)
}

// handleHeraldDelete deletes a rule.
func (s *Server) handleHeraldDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.herald.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("delete rule: %v", err), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/herald", http.StatusSeeOther)
}

// heraldFormCSS is a style block for Herald forms that handles both light and dark mode.
const heraldFormCSS = `<style>
.herald-rule-form input[type="text"],
.herald-rule-form select,
.herald-rule-form textarea {
  border: 1px solid #c7ccd9;
  border-radius: 3px;
  background: #fff;
  color: #292e36;
  box-sizing: border-box;
}
.herald-rule-form .herald-add-btn {
  border: 1px solid #c7ccd9;
  border-radius: 3px;
  background: #fff;
  color: #292e36;
}
.herald-rule-form .herald-delete-btn {
  border: 1px solid #c7ccd9;
  border-radius: 3px;
  color: #464c5c;
}
.herald-rule-form label.herald-label {
  color: #464c5c;
}
.herald-condition-table td,
.herald-action-table td {
  padding: 8px 4px;
}
.phui-theme-dark .herald-rule-form input[type="text"],
.phui-theme-dark .herald-rule-form select,
.phui-theme-dark .herald-rule-form textarea {
  border-color: #464C5C;
  background: #1B2028;
  color: #C8D1DB;
}
.phui-theme-dark .herald-rule-form .herald-add-btn {
  border-color: #464C5C;
  background: #1B2028;
  color: #C8D1DB;
}
.phui-theme-dark .herald-rule-form .herald-delete-btn {
  border-color: #464C5C;
  color: #C8D1DB;
}
.phui-theme-dark .herald-rule-form label.herald-label {
  color: #C8D1DB;
}
</style>`

// renderHeraldForm builds the HTML form for creating/editing a rule.
func renderHeraldForm(rule *herald.Rule) string {
	var b strings.Builder

	b.WriteString(heraldFormCSS)
	b.WriteString(`<div class="phui-object-box herald-rule-form">`)
	b.WriteString(`<div class="phui-header-view">`)
	if rule != nil {
		b.WriteString(`<h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa fa-pencil"></span>Edit Rule</h1>`)
	} else {
		b.WriteString(`<h1 class="phui-header-header"><span class="phui-header-icon phui-icon-view phui-font-fa fa-pencil"></span>Create Rule</h1>`)
	}
	b.WriteString(`</div>`)

	b.WriteString(`<div class="phui-form-view" style="padding:16px;">`)
	b.WriteString(`<form method="POST" action="/herald/save">`)

	id := ""
	name := ""
	matchAll := true
	disabled := false
	var conditions []herald.Condition
	var actions []herald.Action
	if rule != nil {
		id = rule.ID
		name = rule.Name
		matchAll = rule.MustMatchAll
		disabled = rule.Disabled
		conditions = rule.Conditions
		actions = rule.Actions
	}

	if id != "" {
		fmt.Fprintf(&b, `<input type="hidden" name="id" value="%s">`, template.HTMLEscapeString(id))
	}

	// Name field
	b.WriteString(`<div class="phui-form-item" style="margin-bottom:12px;">`)
	b.WriteString(`<label class="herald-label" style="display:block;font-weight:bold;margin-bottom:4px;">Rule Name</label>`)
	fmt.Fprintf(&b, `<input type="text" name="name" value="%s" required class="aphront-form-input" style="width:100%%;max-width:460px;padding:6px 8px;">`,
		template.HTMLEscapeString(name))
	b.WriteString(`</div>`)

	// Match mode
	b.WriteString(`<div class="phui-form-item" style="margin-bottom:12px;">`)
	b.WriteString(`<label class="herald-label" style="display:block;font-weight:bold;margin-bottom:4px;">Match Mode</label>`)
	b.WriteString(`<select name="match_mode" class="aphront-form-input" style="padding:4px 8px;">`)
	allSel, anySel := "", ""
	if matchAll {
		allSel = " selected"
	} else {
		anySel = " selected"
	}
	fmt.Fprintf(&b, `<option value="all"%s>All conditions must match</option>`, allSel)
	fmt.Fprintf(&b, `<option value="any"%s>Any condition must match</option>`, anySel)
	b.WriteString(`</select></div>`)

	// Conditions
	b.WriteString(`<div class="phui-form-item" style="margin-bottom:12px;">`)
	b.WriteString(`<label class="herald-label" style="display:block;font-weight:bold;margin-bottom:4px;">Conditions</label>`)
	b.WriteString(`<table class="herald-condition-table" id="herald-conditions">`)
	if len(conditions) == 0 {
		conditions = []herald.Condition{{}}
	}
	for _, c := range conditions {
		renderConditionRow(&b, c)
	}
	b.WriteString(`</table>`)
	b.WriteString(`<button type="button" onclick="addConditionRow()" class="phui-button-view herald-add-btn" style="margin-top:4px;padding:2px 8px;cursor:pointer;font-size:12px;">`)
	b.WriteString(`<span class="phui-icon-view phui-font-fa fa-plus" style="margin-right:4px;"></span>Add Condition</button>`)
	b.WriteString(`</div>`)

	// Actions
	b.WriteString(`<div class="phui-form-item" style="margin-bottom:12px;">`)
	b.WriteString(`<label class="herald-label" style="display:block;font-weight:bold;margin-bottom:4px;">Actions</label>`)
	b.WriteString(`<table class="herald-action-table" id="herald-actions">`)
	if len(actions) == 0 {
		actions = []herald.Action{{}}
	}
	for _, a := range actions {
		renderActionRow(&b, a)
	}
	b.WriteString(`</table>`)
	b.WriteString(`<button type="button" onclick="addActionRow()" class="phui-button-view herald-add-btn" style="margin-top:4px;padding:2px 8px;cursor:pointer;font-size:12px;">`)
	b.WriteString(`<span class="phui-icon-view phui-font-fa fa-plus" style="margin-right:4px;"></span>Add Action</button>`)
	b.WriteString(`</div>`)

	// Disabled toggle
	if rule != nil {
		checked := ""
		if disabled {
			checked = " checked"
		}
		b.WriteString(`<div class="phui-form-item" style="margin-bottom:12px;">`)
		fmt.Fprintf(&b, `<label><input type="checkbox" name="disabled" value="1"%s> Disable this rule</label>`, checked)
		b.WriteString(`</div>`)
	}

	// Submit
	b.WriteString(`<div class="phui-form-item">`)
	b.WriteString(`<button type="submit" class="phui-button-view button-green" style="padding:6px 20px;border-radius:3px;background:#139543;color:#fff;border:none;cursor:pointer;">`)
	b.WriteString(`<span class="phui-icon-view phui-font-fa fa-check" style="margin-right:4px;"></span>Save Rule</button>`)

	if rule != nil {
		fmt.Fprintf(&b, ` <a href="/herald/%s/delete" class="phui-button-view herald-delete-btn" style="padding:6px 20px;text-decoration:none;margin-left:8px;">`+
			`<span class="phui-icon-view phui-font-fa fa-trash" style="margin-right:4px;"></span>Delete</a>`,
			template.HTMLEscapeString(id))
	}
	b.WriteString(`</div>`)

	b.WriteString(`</form></div>`)
	b.WriteString(`</div>`) // object-box

	// Inline JS for add/remove rows
	b.WriteString(`<script>` + heraldFormJS + `</script>`)

	return b.String()
}

func renderConditionRow(b *strings.Builder, c herald.Condition) {
	b.WriteString(`<tr>`)
	b.WriteString(`<td><select name="cond_type" class="aphront-form-input" style="width:160px;padding:2px 4px;">`)
	condOpts := []struct {
		val, label string
	}{
		{"", "-- select --"},
		{"file_path", "File path matches"},
		{"author", "Author is"},
		{"title", "Title contains"},
		{"label", "Label is"},
		{"base_branch", "Base branch is"},
	}
	for _, o := range condOpts {
		sel := ""
		if string(c.Type) == o.val {
			sel = " selected"
		}
		fmt.Fprintf(b, `<option value="%s"%s>%s</option>`, o.val, sel, o.label)
	}
	b.WriteString(`</select></td>`)
	b.WriteString(`<td class="value">`)
	fmt.Fprintf(b, `<input type="text" name="cond_value" value="%s" placeholder="value" class="aphront-form-input" style="width:95%%;max-width:460px;padding:2px 4px;">`,
		template.HTMLEscapeString(c.Value))
	b.WriteString(`</td>`)
	b.WriteString(`<td class="remove-column"><button type="button" onclick="this.closest('tr').remove()" style="border:none;background:none;cursor:pointer;color:#92969d;font-size:14px;">`)
	b.WriteString(`<span class="phui-icon-view phui-font-fa fa-times"></span></button></td>`)
	b.WriteString(`</tr>`)
}

func renderActionRow(b *strings.Builder, a herald.Action) {
	b.WriteString(`<tr>`)
	b.WriteString(`<td><select name="action_type" class="aphront-form-input" style="width:160px;padding:2px 4px;">`)
	actOpts := []struct {
		val, label string
	}{
		{"", "-- select --"},
		{"add_reviewer", "Add reviewer"},
		{"add_label", "Add label"},
		{"post_comment", "Post comment"},
	}
	for _, o := range actOpts {
		sel := ""
		if string(a.Type) == o.val {
			sel = " selected"
		}
		fmt.Fprintf(b, `<option value="%s"%s>%s</option>`, o.val, sel, o.label)
	}
	b.WriteString(`</select></td>`)
	b.WriteString(`<td class="value">`)
	fmt.Fprintf(b, `<input type="text" name="action_value" value="%s" placeholder="value" class="aphront-form-input" style="width:95%%;max-width:460px;padding:2px 4px;">`,
		template.HTMLEscapeString(a.Value))
	b.WriteString(`</td>`)
	b.WriteString(`<td class="remove-column"><button type="button" onclick="this.closest('tr').remove()" style="border:none;background:none;cursor:pointer;color:#92969d;font-size:14px;">`)
	b.WriteString(`<span class="phui-icon-view phui-font-fa fa-times"></span></button></td>`)
	b.WriteString(`</tr>`)
}

// renderHeraldCurtainPanel renders a Herald section for the PR curtain sidebar.
func renderHeraldCurtainPanel(b *strings.Builder, matches []herald.RuleMatch) {
	b.WriteString(`<div class="mood-curtain-box">`)
	b.WriteString(`<div class="mood-curtain-title"><span class="phui-icon-view phui-font-fa fa-bullhorn mrs"></span>Herald</div>`)

	if len(matches) == 0 {
		b.WriteString(`<div style="font-size:12px;color:#6b748c">No rules matched.</div>`)
	} else {
		for _, m := range matches {
			b.WriteString(`<div style="font-size:12px;margin-bottom:6px">`)
			fmt.Fprintf(b, `<span class="phui-icon-view phui-font-fa fa-check mrs" style="color:#139543"></span>`)
			fmt.Fprintf(b, `<a href="/herald/%s"><strong>%s</strong></a>`,
				template.HTMLEscapeString(m.Rule.ID),
				template.HTMLEscapeString(m.Rule.Name))

			if len(m.Actions) > 0 {
				b.WriteString(`<div style="font-size:11px;color:#6b748c;margin-left:20px">`)
				for _, a := range m.Actions {
					icon := actionIcon(a.Type)
					fmt.Fprintf(b, `<div><span class="phui-icon-view phui-font-fa %s mrs"></span>%s: %s</div>`,
						icon, actionLabel(a.Type), template.HTMLEscapeString(a.Value))
				}
				b.WriteString(`</div>`)
			}
			b.WriteString(`</div>`)
		}
	}

	b.WriteString(`<div style="margin-top:8px;font-size:12px">`)
	b.WriteString(`<a href="/herald"><span class="phui-icon-view phui-font-fa fa-list mrs"></span>Manage Rules</a>`)
	b.WriteString(`</div>`)
	b.WriteString(`</div>`)
}

func actionIcon(t herald.ActionType) string {
	switch t {
	case herald.ActionAddReviewer:
		return "fa-user-plus"
	case herald.ActionAddLabel:
		return "fa-tag"
	case herald.ActionPostComment:
		return "fa-comment"
	default:
		return "fa-bolt"
	}
}

func actionLabel(t herald.ActionType) string {
	switch t {
	case herald.ActionAddReviewer:
		return "Add reviewer"
	case herald.ActionAddLabel:
		return "Add label"
	case herald.ActionPostComment:
		return "Post comment"
	default:
		return string(t)
	}
}

const heraldFormJS = `
function addConditionRow() {
  var tbody = document.getElementById('herald-conditions');
  var tr = document.createElement('tr');
  tr.innerHTML = '<td><select name="cond_type" class="aphront-form-input" style="width:160px;padding:2px 4px;">' +
    '<option value="">-- select --</option>' +
    '<option value="file_path">File path matches</option>' +
    '<option value="author">Author is</option>' +
    '<option value="title">Title contains</option>' +
    '<option value="label">Label is</option>' +
    '<option value="base_branch">Base branch is</option>' +
    '</select></td>' +
    '<td class="value"><input type="text" name="cond_value" placeholder="value" class="aphront-form-input" style="width:95%;max-width:460px;padding:2px 4px;"></td>' +
    '<td class="remove-column"><button type="button" onclick="this.closest(\'tr\').remove()" style="border:none;background:none;cursor:pointer;color:#92969d;font-size:14px;"><span class="phui-icon-view phui-font-fa fa-times"></span></button></td>';
  tbody.appendChild(tr);
}
function addActionRow() {
  var tbody = document.getElementById('herald-actions');
  var tr = document.createElement('tr');
  tr.innerHTML = '<td><select name="action_type" class="aphront-form-input" style="width:160px;padding:2px 4px;">' +
    '<option value="">-- select --</option>' +
    '<option value="add_reviewer">Add reviewer</option>' +
    '<option value="add_label">Add label</option>' +
    '<option value="post_comment">Post comment</option>' +
    '</select></td>' +
    '<td class="value"><input type="text" name="action_value" placeholder="value" class="aphront-form-input" style="width:95%;max-width:460px;padding:2px 4px;"></td>' +
    '<td class="remove-column"><button type="button" onclick="this.closest(\'tr\').remove()" style="border:none;background:none;cursor:pointer;color:#92969d;font-size:14px;"><span class="phui-icon-view phui-font-fa fa-times"></span></button></td>';
  tbody.appendChild(tr);
}
`
