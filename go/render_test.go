package main

import (
	"html/template"
	"testing"
)

// ── groupWork ────────────────────────────────────────────────────────────────

func TestGroupWork_SingleEntry(t *testing.T) {
	entries := []WorkEntry{
		{Employer: "Acme", Position: "Engineer", StartDate: "2020-01", EndDate: "2022-01"},
	}
	groups := groupWork(entries)
	if len(groups) != 1 {
		t.Fatalf("want 1 group, got %d", len(groups))
	}
	if groups[0].DisplayName != "Acme" {
		t.Errorf("DisplayName: want Acme, got %s", groups[0].DisplayName)
	}
	if len(groups[0].Positions) != 1 {
		t.Errorf("want 1 position, got %d", len(groups[0].Positions))
	}
}

func TestGroupWork_ConsecutiveSameEmployer(t *testing.T) {
	entries := []WorkEntry{
		{Employer: "JumpCloud", Position: "Senior SWE", StartDate: "2022-10", EndDate: "2026-05"},
		{Employer: "JumpCloud", Position: "SWE 3", StartDate: "2021-05", EndDate: "2022-10"},
	}
	groups := groupWork(entries)
	if len(groups) != 1 {
		t.Fatalf("want 1 group, got %d", len(groups))
	}
	g := groups[0]
	if g.DisplayName != "JumpCloud" {
		t.Errorf("DisplayName: want JumpCloud, got %s", g.DisplayName)
	}
	if len(g.Positions) != 2 {
		t.Errorf("want 2 positions, got %d", len(g.Positions))
	}
	if g.StartDate != "2021-05" {
		t.Errorf("StartDate: want 2021-05, got %s", g.StartDate)
	}
	if g.EndDate != "2026-05" {
		t.Errorf("EndDate: want 2026-05, got %s", g.EndDate)
	}
}

func TestGroupWork_NonConsecutiveSameEmployer(t *testing.T) {
	entries := []WorkEntry{
		{Employer: "Acme", Position: "Senior", StartDate: "2023-01", EndDate: ""},
		{Employer: "Other", Position: "Mid", StartDate: "2020-01", EndDate: "2023-01"},
		{Employer: "Acme", Position: "Junior", StartDate: "2018-01", EndDate: "2020-01"},
	}
	groups := groupWork(entries)
	if len(groups) != 3 {
		t.Fatalf("want 3 groups (non-consecutive same employer = separate), got %d", len(groups))
	}
}

func TestGroupWork_EmployerGroupField(t *testing.T) {
	// Explicit employerGroup: displayName = group key; subsidiaries don't
	// appear in FormerNames (they show up in per-position job-meta instead).
	entries := []WorkEntry{
		{Employer: "SubCo", EmployerGroup: "ParentCorp", Position: "Senior", StartDate: "2022-01", EndDate: ""},
		{Employer: "ParentCorp", EmployerGroup: "ParentCorp", Position: "Junior", StartDate: "2020-01", EndDate: "2022-01"},
	}
	groups := groupWork(entries)
	if len(groups) != 1 {
		t.Fatalf("want 1 group via EmployerGroup key, got %d", len(groups))
	}
	g := groups[0]
	if g.DisplayName != "ParentCorp" {
		t.Errorf("DisplayName: want ParentCorp (group key), got %s", g.DisplayName)
	}
	if len(g.FormerNames) != 0 {
		t.Errorf("FormerNames: want [] (subsidiaries not formerNames), got %v", g.FormerNames)
	}
	if g.EndDate != "" {
		t.Errorf("EndDate: want empty (Present), got %q", g.EndDate)
	}
}

func TestGroupWork_NameChange(t *testing.T) {
	// No explicit employerGroup; consecutive same key via employer → name change
	// tracked in FormerNames.
	entries := []WorkEntry{
		{Employer: "NewName", Position: "Senior", StartDate: "2022-01", EndDate: ""},
		{Employer: "OldName", Position: "Junior", StartDate: "2020-01", EndDate: "2022-01"},
	}
	groups := groupWork(entries)
	if len(groups) != 2 {
		// Different employer names with no shared employerGroup → separate groups
		t.Fatalf("want 2 groups (no shared key), got %d", len(groups))
	}
}

func TestGroupWork_PresentEndDate(t *testing.T) {
	entries := []WorkEntry{
		{Employer: "Acme", Position: "Senior", StartDate: "2022-01", EndDate: ""},
		{Employer: "Acme", Position: "Junior", StartDate: "2020-01", EndDate: "2022-01"},
	}
	groups := groupWork(entries)
	if len(groups) != 1 {
		t.Fatalf("want 1 group, got %d", len(groups))
	}
	if groups[0].EndDate != "" {
		t.Errorf("EndDate: want empty (Present), got %q", groups[0].EndDate)
	}
}

func TestGroupWork_MultipleGroups(t *testing.T) {
	entries := []WorkEntry{
		{Employer: "JumpCloud", Position: "Senior SWE", StartDate: "2022-10", EndDate: "2026-05"},
		{Employer: "JumpCloud", Position: "SWE 3", StartDate: "2021-05", EndDate: "2022-10"},
		{Employer: "ObjectRocket", Position: "SWE", StartDate: "2019-10", EndDate: "2021-05"},
		{Employer: "Rackspace", Position: "Admin", StartDate: "2014-01", EndDate: "2019-10"},
	}
	groups := groupWork(entries)
	if len(groups) != 3 {
		t.Fatalf("want 3 groups, got %d", len(groups))
	}
	if groups[0].DisplayName != "JumpCloud" {
		t.Errorf("group 0: want JumpCloud, got %s", groups[0].DisplayName)
	}
	if len(groups[0].Positions) != 2 {
		t.Errorf("group 0 positions: want 2, got %d", len(groups[0].Positions))
	}
	if groups[1].DisplayName != "ObjectRocket" {
		t.Errorf("group 1: want ObjectRocket, got %s", groups[1].DisplayName)
	}
	if groups[2].DisplayName != "Rackspace" {
		t.Errorf("group 2: want Rackspace, got %s", groups[2].DisplayName)
	}
}

// ── formatDate ───────────────────────────────────────────────────────────────

func TestFormatDate(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", "Present"},
		{"2006", "2006"},
		{"2022-10", "Oct 2022"},
		{"2022-10-03", "Oct 2022"},
		{"2021-05-16", "May 2021"},
		{"2009-12-15", "Dec 2009"},
		{"2014-01-06", "Jan 2014"},
	}
	for _, c := range cases {
		got := formatDate(c.in)
		if got != c.want {
			t.Errorf("formatDate(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// ── nbspShortWords ───────────────────────────────────────────────────────────

func TestNbspShortWords(t *testing.T) {
	cases := []struct {
		in   string
		want template.HTML
	}{
		// nbsp precedes ≤4-char words (binds short word to preceding word)
		{"I own problems", "I&nbsp;own problems"},
		{"I own problems end-to-end", "I&nbsp;own problems end-to-end"},
		{"hello world", "hello world"},
		{"12+ years building systems", "12+ years building systems"},
		{"a b c longer", "a&nbsp;b&nbsp;c longer"},
		{"Focused on reliability", "Focused&nbsp;on reliability"},
		{"building and maintaining", "building&nbsp;and maintaining"},
		{"single", "single"},
	}
	for _, c := range cases {
		got := nbspShortWords(c.in)
		if got != c.want {
			t.Errorf("nbspShortWords(%q)\n  got  %q\n  want %q", c.in, got, c.want)
		}
	}
}

// ── levelClass ───────────────────────────────────────────────────────────────

func TestLevelClass(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Advanced", "adv"},
		{"Intermediate", "mid"},
		{"Familiar", ""},
		{"Beginner", ""},
		{"", ""},
	}
	for _, c := range cases {
		got := levelClass(c.in)
		if got != c.want {
			t.Errorf("levelClass(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// ── skillByName ──────────────────────────────────────────────────────────────

func TestSkillByName(t *testing.T) {
	list := []SkillItem{
		{Name: "Go", Level: "Advanced", Years: 6},
		{Name: "Python", Level: "Advanced", Years: 8},
		{Name: "Perl", Level: "Familiar", Years: 3},
	}
	got := skillByName(list, "Python")
	if got.Name != "Python" || got.Level != "Advanced" {
		t.Errorf("skillByName Python: got %+v", got)
	}
	missing := skillByName(list, "Rust")
	if missing.Name != "Rust" {
		t.Errorf("missing skill: want Name=Rust, got %+v", missing)
	}
}
