package herald

import (
	"path/filepath"
	"strings"
)

// Evaluate runs all enabled rules against the given PR context and returns matches.
func Evaluate(rules []Rule, ctx *PRContext) []RuleMatch {
	var matches []RuleMatch
	for i := range rules {
		r := &rules[i]
		if r.Disabled {
			continue
		}
		if matchRule(r, ctx) {
			matches = append(matches, RuleMatch{
				Rule:    r,
				Actions: r.Actions,
			})
		}
	}
	return matches
}

func matchRule(r *Rule, ctx *PRContext) bool {
	if len(r.Conditions) == 0 {
		return true
	}
	if r.MustMatchAll {
		for _, c := range r.Conditions {
			if !matchCondition(&c, ctx) {
				return false
			}
		}
		return true
	}
	// OR mode
	for _, c := range r.Conditions {
		if matchCondition(&c, ctx) {
			return true
		}
	}
	return false
}

func matchCondition(c *Condition, ctx *PRContext) bool {
	switch c.Type {
	case CondAuthor:
		return strings.EqualFold(ctx.Author, c.Value)
	case CondTitle:
		return strings.Contains(strings.ToLower(ctx.Title), strings.ToLower(c.Value))
	case CondLabel:
		for _, l := range ctx.Labels {
			if strings.EqualFold(l, c.Value) {
				return true
			}
		}
		return false
	case CondBaseBranch:
		return ctx.BaseBranch == c.Value
	case CondFilePath:
		for _, f := range ctx.ChangedFiles {
			if matched, _ := filepath.Match(c.Value, f); matched {
				return true
			}
			// Also try matching against just the filename
			if matched, _ := filepath.Match(c.Value, filepath.Base(f)); matched {
				return true
			}
		}
		return false
	}
	return false
}
