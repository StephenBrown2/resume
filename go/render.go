package main

import (
	"html/template"
	"math/rand/v2"
	"slices"
	"sort"
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

// filterWorkSince returns only entries whose endDate is absent (current) or
// falls on or after the cutoff. The cutoff string may be YYYY, YYYY-MM, or
// YYYY-MM-DD; shorter forms are right-padded to YYYY-MM-DD for comparison.
func filterWorkSince(entries []WorkEntry, since string) []WorkEntry {
	cutoff := padDate(since)
	var out []WorkEntry
	for _, e := range entries {
		if e.EndDate == "" || padDate(e.EndDate) >= cutoff {
			out = append(out, e)
		}
	}
	return out
}

// padDate right-pads a partial ISO date to YYYY-MM-DD for lexicographic comparison.
func padDate(s string) string {
	switch len(s) {
	case 4: // YYYY
		return s + "-01-01"
	case 7: // YYYY-MM
		return s + "-01"
	default:
		return s
	}
}

// filterProjectsSince returns only projects whose endDate is absent (current)
// or falls on or after the cutoff. Uses the same padDate logic as filterWorkSince.
func filterProjectsSince(projects []Project, since string) []Project {
	cutoff := padDate(since)
	var out []Project
	for _, p := range projects {
		if p.EndDate == "" || padDate(p.EndDate) >= cutoff {
			out = append(out, p)
		}
	}
	return out
}

// shuffleKeywords randomises keyword order for work entries and projects.
func shuffleKeywords(resume *Resume) {
	for i := range resume.Work {
		rand.Shuffle(len(resume.Work[i].Keywords), func(a, b int) {
			resume.Work[i].Keywords[a], resume.Work[i].Keywords[b] = resume.Work[i].Keywords[b], resume.Work[i].Keywords[a]
		})
	}
	for i := range resume.Projects {
		rand.Shuffle(len(resume.Projects[i].Keywords), func(a, b int) {
			resume.Projects[i].Keywords[a], resume.Projects[i].Keywords[b] = resume.Projects[i].Keywords[b], resume.Projects[i].Keywords[a]
		})
	}
}

// sortSkillSets sorts each domain's skill names by proficiency descending,
// then alphabetically ascending within the same level.
func sortSkillSets(sets []SkillSet, list []SkillItem) {
	rank := map[string]int{"Advanced": 3, "Intermediate": 2, "Familiar": 1, "Beginner": 0}
	idx := make(map[string]SkillItem, len(list))
	for _, s := range list {
		idx[s.Name] = s
	}
	for i := range sets {
		sort.SliceStable(sets[i].Skills, func(a, b int) bool {
			ra := rank[idx[sets[i].Skills[a]].Level]
			rb := rank[idx[sets[i].Skills[b]].Level]
			if ra != rb {
				return ra > rb
			}
			return sets[i].Skills[a] < sets[i].Skills[b]
		})
	}
}

// sortTestimonials sorts each testimonial by length of the quote ascending.
func sortTestimonials(list []Testimonial) {
	slices.SortStableFunc(list, func(a, b Testimonial) int {
		if len(a.Quote) < len(b.Quote) {
			return -1
		}
		if len(a.Quote) > len(b.Quote) {
			return 1
		}
		return 0
	})
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

// fullDate returns the long-form date with day for use in title/hover text
// ("January 2, 2006"). Returns empty string for dates without a day component
// or for empty (Present) dates.
func fullDate(iso string) string {
	if iso == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", iso)
	if err != nil {
		return ""
	}
	return t.Format("January 2, 2006")
}

// fullDateRange returns "Start - End" using fullDate for both ends, for use
// as a title/hover on a single element showing a date range.
func fullDateRange(start, end string) string {
	s := fullDate(start)
	e := fullDate(end)
	if s == "" && e == "" {
		return ""
	}
	if s == "" {
		return e
	}
	if e == "" {
		return s
	}
	return s + " - " + e
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
func certTitle(c Certificate) string {
	if c.ID == "" {
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
	if c.ID == "" {
		return ""
	}
	if c.VerificationCode != "" {
		return c.ID + " / " + c.VerificationCode
	}
	return c.ID
}

// CertGroup is a set of certificates from the same issuer.
type CertGroup struct {
	Issuer                 string
	Certs                  []Certificate
	SharedID               string // non-empty when all certs share the same ID
	SharedVerificationCode string
}

// groupCerts groups a sorted certificate slice by issuer, preserving order of
// first appearance. Sets SharedID when every cert in the group has the same ID.
func groupCerts(certs []Certificate) []CertGroup {
	seen := map[string]int{}
	var groups []CertGroup
	for _, c := range certs {
		issuer := c.Issuer
		if issuer == "" {
			issuer = "Other"
		}
		if i, ok := seen[issuer]; ok {
			groups[i].Certs = append(groups[i].Certs, c)
		} else {
			seen[issuer] = len(groups)
			groups = append(groups, CertGroup{Issuer: issuer, Certs: []Certificate{c}})
		}
	}
	for i, g := range groups {
		if len(g.Certs) == 0 {
			continue
		}
		id, vc := g.Certs[0].ID, g.Certs[0].VerificationCode
		allSame := id != ""
		for _, c := range g.Certs[1:] {
			if c.ID != id || c.VerificationCode != vc {
				allSame = false
				break
			}
		}
		if allSame {
			groups[i].SharedID = id
			groups[i].SharedVerificationCode = vc
		}
	}
	return groups
}

// certGroupNames formats the cert names within a group:
// two certs joined with " & ", three or more joined with ", ".
func certGroupNames(g CertGroup) string {
	names := make([]string, len(g.Certs))
	for i, c := range g.Certs {
		names[i] = c.Name
	}
	if len(names) == 2 {
		return names[0] + " & " + names[1]
	}
	return strings.Join(names, ", ")
}

// certGroupID returns the shared ID parenthetical (e.g. "140-027-434") or "".
func certGroupID(g CertGroup) string {
	if g.SharedID == "" {
		return ""
	}
	if g.SharedVerificationCode != "" {
		return g.SharedID + " / " + g.SharedVerificationCode
	}
	return g.SharedID
}

// TemplateData is the top-level data passed to the HTML template.
type TemplateData struct {
	Basics          Basics
	EmployerGroups  []EmployerGroup
	Projects        []Project
	SkillSets       []SkillSet
	SkillList       []SkillItem
	Certificates    []Certificate
	CertGroups      []CertGroup
	Education       []Education
	Languages       []Language
	Interests       []Interest
	Testimonials    []Testimonial
	GoogleFontsLink template.HTML
	NameFontCSS     template.CSS
}
