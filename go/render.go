package main

import (
	"html/template"
	"strings"
	"time"
	"unicode/utf8"
)

// EmployerGroup aggregates consecutive work entries sharing the same group key.
type EmployerGroup struct {
	DisplayName string
	FormerNames []string
	URL         string
	StartDate   string
	EndDate     string
	Positions   []WorkEntry
}

// groupWork groups consecutive WorkEntry items sharing the same group key.
// Work entries are expected in reverse-chronological order (most recent first).
func groupWork(entries []WorkEntry) []EmployerGroup {
	var groups []EmployerGroup
	for _, entry := range entries {
		key := entry.EmployerGroup
		if key == "" {
			key = entry.Employer
		}
		if len(groups) > 0 && groups[len(groups)-1].key() == key {
			g := &groups[len(groups)-1]
			// Track name changes (not subsidiaries — those show in job-meta).
			if entry.Employer != g.DisplayName && entry.EmployerGroup == "" {
				found := false
				for _, fn := range g.FormerNames {
					if fn == entry.Employer {
						found = true
						break
					}
				}
				if !found {
					g.FormerNames = append(g.FormerNames, entry.Employer)
				}
			}
			// Prefer URL from the first appended entry that provides one.
			if g.URL == "" && entry.URL != "" {
				g.URL = entry.URL
			}
			if entry.StartDate != "" && (g.StartDate == "" || entry.StartDate < g.StartDate) {
				g.StartDate = entry.StartDate
			}
			if entry.EndDate == "" || (g.EndDate != "" && entry.EndDate > g.EndDate) {
				g.EndDate = entry.EndDate
			}
			if entry.EndDate == "" {
				g.EndDate = ""
			}
			g.Positions = append(g.Positions, entry)
		} else {
			displayName := entry.Employer
			url := entry.URL
			if entry.EmployerGroup != "" {
				// Explicit group key: use it as the display name.
				// URL deferred until we find the entry whose employer matches.
				displayName = entry.EmployerGroup
				url = ""
			}
			groups = append(groups, EmployerGroup{
				DisplayName: displayName,
				URL:         url,
				StartDate:   entry.StartDate,
				EndDate:     entry.EndDate,
				Positions:   []WorkEntry{entry},
			})
		}
	}
	return groups
}

// key returns the group key (employerGroup if set, else employer of first position).
func (g EmployerGroup) key() string {
	if len(g.Positions) == 0 {
		return g.DisplayName
	}
	first := g.Positions[0]
	if first.EmployerGroup != "" {
		return first.EmployerGroup
	}
	return first.Employer
}

var dateLayouts = []string{
	"2006-01-02",
	"2006-01",
	"2006",
}

// formatDate converts ISO 8601 date strings to display form ("Jan 2006" or "2006").
func formatDate(iso string) string {
	if iso == "" {
		return "Present"
	}
	for _, layout := range dateLayouts {
		t, err := time.Parse(layout, iso)
		if err == nil {
			if layout == "2006" {
				return t.Format("2006")
			}
			return t.Format("Jan 2006")
		}
	}
	return iso
}

// nbspShortWords inserts &nbsp; before words of ≤4 chars, binding short
// connector words to their predecessor so they don't start a line alone.
func nbspShortWords(s string) template.HTML {
	words := strings.Split(s, " ")
	var parts []string
	for i, word := range words {
		escaped := template.HTMLEscapeString(word)
		if i == len(words)-1 {
			parts = append(parts, escaped)
		} else if utf8.RuneCountInString(words[i+1]) <= 4 {
			parts = append(parts, escaped+"&nbsp;")
		} else {
			parts = append(parts, escaped+" ")
		}
	}
	return template.HTML(strings.Join(parts, ""))
}

// levelClass returns the CSS modifier class for a skill level string.
func levelClass(level string) string {
	switch level {
	case "Advanced":
		return "adv"
	case "Intermediate":
		return "mid"
	default:
		return ""
	}
}

// skillByName looks up a SkillItem by name from the list.
func skillByName(list []SkillItem, name string) SkillItem {
	for _, item := range list {
		if item.Name == name {
			return item
		}
	}
	return SkillItem{Name: name}
}

// stripScheme removes https:// or http:// prefix for display.
func stripScheme(u string) string {
	u = strings.TrimPrefix(u, "https://")
	u = strings.TrimPrefix(u, "http://")
	return u
}

// certTitle builds the tooltip string for a certificate with an ID.
// Returns "" if the ID is already embedded in the cert URL (no need to repeat it).
func certTitle(c Certificate) string {
	if c.ID == "" || strings.Contains(c.URL, c.ID) {
		return ""
	}
	s := "ID: " + c.ID
	if c.VerificationCode != "" {
		s += " · Verification Code: " + c.VerificationCode
	}
	return s
}

// certPrintID builds the parenthetical print-only ID string.
// Returns "" if the ID is already embedded in the cert URL (no need to repeat it).
func certPrintID(c Certificate) string {
	if c.ID == "" || strings.Contains(c.URL, c.ID) {
		return ""
	}
	if c.VerificationCode != "" {
		return c.ID + " / " + c.VerificationCode
	}
	return c.ID
}

// TemplateData is the top-level data passed to the HTML template.
type TemplateData struct {
	Basics          Basics
	EmployerGroups  []EmployerGroup
	Projects        []Project
	SkillSets       []SkillSet
	SkillList       []SkillItem
	Certificates    []Certificate
	Education       []Education
	Languages       []Language
	Interests       []Interest
	Testimonials    []Testimonial
	GoogleFontsLink template.HTML
	NameFontCSS     template.CSS
}
