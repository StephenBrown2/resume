# Go implementation plan

**Read `plans/shared-context.md` first.** This document covers only Go-specific decisions.

---

## Goal

Replace the two existing schema-type files (`json-resume.go`, `fresh-resume.go`) and the `goresume` CLI dependency with a single self-contained Go program in `go/`. The program reads `resume.yaml`, groups work entries by employer, and writes `docs/index.html`.

---

## Language version and module

- Go 1.26.3 (latest stable). Use `go 1.26.3` in `go.mod`.
- Module path: `github.com/StephenBrown2/resume/go` (or `resume-renderer` - use whatever fits `go.mod` conventions for a local tool).
- Use Go's standard library wherever possible; external dependencies should be minimal.

---

## Dependencies

| Package | Version | Purpose |
|---|---|---|
| `github.com/goccy/go-yaml` | latest | YAML parsing (faster than `gopkg.in/yaml.v3`; better error messages; handles anchors and complex tags correctly) |
| `github.com/santhosh-tekuri/jsonschema/v6` | v6.0.2 | JSON Schema Draft 2020-12 validation |
| `html/template` | stdlib | HTML rendering with auto-escaping |

No other external dependencies. `html/template` handles HTML escaping automatically; use it instead of `text/template`.

---

## File structure

```
go/
  go.mod
  go.sum
  main.go         # CLI entry point, flag parsing
  resume.go       # data model structs
  render.go       # HTML rendering logic, template funcs
  template.go     # embedded HTML template string (using go:embed or a raw string literal)
```

Delete `json-resume.go` and `fresh-resume.go` from the repo root - they were schema type definitions for a different tool and are now superseded.

---

## Data model (`resume.go`)

Define structs matching the hybrid YAML schema. Use `yaml:"field"` tags throughout. All optional fields should use pointer types or `omitempty`.

```go
type Resume struct {
    Basics       Basics        `yaml:"basics"`
    Disposition  Disposition   `yaml:"disposition"`
    Work         []WorkEntry   `yaml:"work"`
    Projects     []Project     `yaml:"projects"`
    Skills       Skills        `yaml:"skills"`
    Certificates []Certificate `yaml:"certificates"`
    Education    []Education   `yaml:"education"`
    Languages    []Language    `yaml:"languages"`
    Interests    []Interest    `yaml:"interests"`
    Testimonials []Testimonial `yaml:"testimonials"`
    References   []Reference   `yaml:"references"`
}

type Basics struct {
    Name     string   `yaml:"name"`
    Label    string   `yaml:"label"`
    Email    string   `yaml:"email"`
    Phone    string   `yaml:"phone"`
    URL      string   `yaml:"url"`
    Summary  string   `yaml:"summary"`
    Location Location `yaml:"location"`
    Profiles []Profile `yaml:"profiles"`
}

type Location struct {
    City        string `yaml:"city"`
    Region      string `yaml:"region"`
    CountryCode string `yaml:"countryCode"`
}

type Profile struct {
    Network  string `yaml:"network"`
    Username string `yaml:"username"`
    URL      string `yaml:"url"`
}

type Disposition struct {
    Travel        int        `yaml:"travel"`
    Authorization string     `yaml:"authorization"`
    Commitment    []string   `yaml:"commitment"`
    Remote        bool       `yaml:"remote"`
    Relocation    Relocation `yaml:"relocation"`
}

type Relocation struct {
    Willing      bool     `yaml:"willing"`
    Destinations []string `yaml:"destinations"`
}

type WorkEntry struct {
    Employer      string   `yaml:"employer"`
    EmployerGroup string   `yaml:"employerGroup"` // optional; overrides grouping key
    Position      string   `yaml:"position"`
    URL           string   `yaml:"url"`
    StartDate     string   `yaml:"startDate"`
    EndDate       string   `yaml:"endDate"`
    Summary       string   `yaml:"summary"`
    Location      string   `yaml:"location"`
    Highlights    []string `yaml:"highlights"`
    Keywords      []string `yaml:"keywords"`
}

type Skills struct {
    Sets []SkillSet  `yaml:"sets"`
    List []SkillItem `yaml:"list"`
}

type SkillSet struct {
    Name   string   `yaml:"name"`
    Skills []string `yaml:"skills"`
}

type SkillItem struct {
    Name    string `yaml:"name"`
    Level   string `yaml:"level"`
    Summary string `yaml:"summary"`
    Years   int    `yaml:"years"`
}

// ... Project, Certificate, Education, Language, Interest, Testimonial, Reference
// (follow the same pattern; see schema.json for all fields)
```

---

## Employer grouping (`render.go`)

Define an intermediate type:

```go
type EmployerGroup struct {
    DisplayName string
    FormerNames []string // distinct employer names beyond DisplayName, in order
    URL         string
    StartDate   string   // earliest across all positions
    EndDate     string   // latest (empty = "Present")
    Positions   []WorkEntry
}
```

Implement `groupWork(entries []WorkEntry) []EmployerGroup`:

1. Walk entries in order.
2. For each entry, compute its key: `entry.EmployerGroup` if non-empty, else `entry.Employer`.
3. If the key matches the current open group's key, append the entry to the current group and update `StartDate` / `FormerNames`.
4. Otherwise, close the current group and start a new one.
5. Return the completed groups.

---

## Date formatting (`render.go`)

```go
func formatDate(iso string) string {
    if iso == "" {
        return "Present"
    }
    // Parse YYYY-MM-DD, YYYY-MM, or YYYY
    // Return "Jan 2006" format for YYYY-MM-DD and YYYY-MM
    // Return "2006" for YYYY-only
}
```

Use `time.Parse` with multiple layout attempts. Map month number to 3-letter abbreviation.

---

## Non-breaking space insertion (`render.go`)

The rule is: if the *next* word is ≤4 characters, replace the preceding space with `&nbsp;` (nbsp *precedes* short words, not trails them). This binds short connector words to their predecessor, preventing them from stranding at the start of a line.

```go
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
```

This is used only in the summary paragraph via the `nbspSummary` template function.

---

## HTML template (`template.go`)

Embed the full HTML template as a raw string constant or use `//go:embed template.html`. Use `html/template` syntax.

Key template functions to register:
- `formatDate` - date string → display string
- `nbspSummary` - applies `nbspShortWords` to summary text, returns `template.HTML`
- `levelClass` - level string → CSS class suffix (`"adv"`, `"mid"`, or `""`)
- `skillByName` - looks up a `SkillItem` in the list by name for a given skill set entry

Pass a single data struct to the template:

```go
type TemplateData struct {
    Basics        Basics
    EmployerGroups []EmployerGroup
    Projects      []Project
    SkillSets     []SkillSet
    SkillList     []SkillItem
    Certificates  []Certificate
    Education     []Education
    Languages     []Language
    Interests     []Interest
    Testimonials  []Testimonial
}
```

---

## Schema validation (`main.go`)

After parsing YAML into `map[string]any` (before unmarshalling into the struct), validate against `schema.json` using `santhosh-tekuri/jsonschema/v6`:

```go
import (
    "github.com/santhosh-tekuri/jsonschema/v6"
    goyaml "github.com/goccy/go-yaml"
)

func validateSchema(schemaPath string, data map[string]any) error {
    c := jsonschema.NewCompiler()
    sch, err := c.Compile(schemaPath)
    if err != nil {
        return fmt.Errorf("compile schema: %w", err)
    }
    if err := sch.Validate(data); err != nil {
        return fmt.Errorf("schema validation: %w", err)
    }
    return nil
}
```

Pass `--skip-validation` to bypass this step. The `goccy/go-yaml` library's `Unmarshal` with a `map[string]any` target produces the raw data needed for jsonschema validation, then unmarshal again (or reuse) into the typed struct.

## `--name-font` flag and Google Fonts URL

```go
nameFont := flag.String("name-font", "Instrument Serif", "Google Fonts family for the name heading")

// Convert to URL form
fontURL := strings.ReplaceAll(*nameFont, " ", "+")
googleFontsLink := fmt.Sprintf(
    `<link href="https://fonts.googleapis.com/css2?family=%s:ital@0;1&display=swap" rel="stylesheet">`,
    fontURL,
)
// CSS var
nameFontCSS := fmt.Sprintf("'%s', Georgia, serif", *nameFont)
```

Pass both `googleFontsLink` and `nameFontCSS` into the template data struct so the template can inject them in the right places.

## `main.go`

```go
func main() {
    input    := flag.String("input",    "../resume.yaml",     "path to resume YAML")
    output   := flag.String("output",   "../docs/index.html", "path to write HTML")
    nameFont := flag.String("name-font","Instrument Serif",   "Google Fonts family for name heading")
    skipVal  := flag.Bool("skip-validation", false,           "skip JSON Schema validation")
    flag.Parse()

    data, err := os.ReadFile(*input)
    // handle err

    // Unmarshal to raw map for schema validation
    var raw map[string]any
    err = goyaml.Unmarshal(data, &raw)
    // handle err

    if !*skipVal {
        if err := validateSchema("schema.json", raw); err != nil {
            fmt.Fprintln(os.Stderr, "validation error:", err)
            os.Exit(1)
        }
    }

    // Unmarshal to typed struct
    var resume Resume
    err = goyaml.Unmarshal(data, &resume)
    // handle err

    groups := groupWork(resume.Work)

    tmplData := TemplateData{
        Basics:         resume.Basics,
        EmployerGroups: groups,
        // ...
    }

    tmpl := template.Must(template.New("resume").Funcs(funcMap).Parse(resumeTemplate))
    
    out, err := os.Create(*output)
    // handle err
    defer out.Close()
    
    err = tmpl.Execute(out, tmplData)
    // handle err
    
    fmt.Fprintf(os.Stderr, "wrote %s\n", *output)
}
```

---

## `go/README.md`

Create `go/README.md` documenting this implementation. It should cover:

- **Prerequisites:** Go 1.26.3+. Install via `go install` or a version manager such as `mise` or `asdf`.
- **Build:** `go build -o resume-renderer .` (or `just go-build` from the repo root).
- **Run:** `./resume-renderer [flags]` - list all flags with `./resume-renderer --help`.
- **Flags:** table matching the CLI interface in `shared-context.md` (`--input`, `--output`, `--name-font`, `--skip-validation`).
- **Output:** writes `docs/index.html` (relative to the repo root when using the default path).
- **Module:** note the module path and that no CGo is used.

---

## Justfile integration

Add language-specific recipes to the `justfile` using `[working-directory]` rather than `cd`:

```just
[working-directory: 'go']
go-build:
    go build -o resume-renderer .

[working-directory: 'go']
go-render: go-build
    ./resume-renderer --input ../resume.yaml --output ../docs/index.html

[working-directory: 'go']
go-validate: go-build
    ./resume-renderer --input ../resume.yaml --output /dev/null
```

The generic `build` and `validate` recipes call their language-specific counterparts. As other languages are implemented, their render recipes are added to `build`.

When testing, pass `--output ../docs/go-index.html` to avoid overwriting the canonical `docs/index.html`.

---

## Notes

- The two existing root-level `.go` files (`json-resume.go`, `fresh-resume.go`) were part of a now-deleted tool and should be removed.
- Do not use `text/template`; always use `html/template` to ensure proper escaping.
- The `template.HTML` type in `html/template` is an escape hatch for trusted pre-escaped content - use it only for `nbspSummary` output and `&middot;` / `&amp;` / `&nbsp;` literals in the template itself.
- CSS comments are stripped by `html/template` when it processes `<style>` blocks. Do not use `/* */` comments in the template; they will be replaced with whitespace.
- Multi-line CSS selectors: `html/template`'s CSS sanitizer drops lines that contain only a selector fragment (e.g. `.foo,` with no `{`). Keep comma-separated selectors on a single line: `.foo, .bar { ... }`.
- Bullet character: use the literal `–` en-dash character in the CSS `content` property (`content: '–';`). Do not use the CSS escape `'\2013'` — `html/template`'s CSS context sanitizer may alter backslash sequences, producing a visually different glyph.
- CSS values injected via template variables must use the appropriate `template.CSS` type (not `string`) to avoid `ZgotmplZ` substitution. In particular, `NameFontCSS` in `TemplateData` must be `template.CSS`.
- Build with `go build ./...` from the `go/` directory. No CGo; pure Go only.
