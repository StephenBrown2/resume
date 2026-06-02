# Shared context for all resume renderer implementations

This document is referenced by the five language-specific implementation plans. Read it in full before proceeding with any individual plan.

---

## What you are building

A standalone CLI program that reads `resume.yaml` and writes a styled, single-file HTML resume to `docs/index.html`. The program lives in its own subdirectory of the repo (`go/`, `rust/`, `java/`, `python/`, or `elixir/`) and replaces the current `goresume`-based build step.

The repo root is `/home/stephen/Projects/resume/`. All paths in these plans are relative to that root.

---

## Input: `resume.yaml`

The file uses a hybrid schema combining JSON Resume and FRESH conventions. Read `schema.json` for the full field-level definitions. The key structural points:

### Top-level sections

```
$schema, basics, disposition, work, projects, skills, certificates,
education, languages, interests, testimonials, references
```

### `basics`
```yaml
name, label, email, phone, url, summary
location: { city, region, countryCode }
profiles: [{ network, username, url }]
```

### `disposition`
```yaml
travel: integer (0–100, percent willing to travel)
authorization: string
commitment: [string]
remote: boolean
relocation: { willing: boolean, destinations: [string] }
```

### `work[]`
Each entry is one position. Multiple consecutive entries with the same `employer` represent promotions. An optional `employerGroup` string field may appear to explicitly link entries across name changes.
```yaml
employer, position, url, startDate, endDate, summary, location
highlights: [string]
keywords: [string]
employerGroup: string  # optional; group key overriding employer name
```

### `projects[]`
```yaml
name, description, url, type
highlights: [string]
keywords: [string]
roles: [string]
startDate, endDate
```

### `skills`
```yaml
sets:
  - name: string
    skills: [string]   # references into list[].name
list:
  - name, level, summary, years
```

### `certificates[]`
```yaml
name, date, url, issuer
id: string              # optional; cert ID (e.g. "140-027-434", "LPI000223384")
verificationCode: string  # optional; secondary lookup code (LPI only)
```

When `id` is present, render it as a `title` attribute on the cert link/span (hover tooltip on screen), and in print mode append a `.print-only` span: ` (id)`. If `verificationCode` is also present, append both: ` (id / verificationCode)` and include both in the tooltip: `"ID: {id} · Verification Code: {verificationCode}"`.

Sort certificates by `date` descending (most recent first) before rendering.

### `education[]`
```yaml
institution, url, area, studyType, startDate, endDate, score, location
```

### `languages[]`
```yaml
language, fluency, years
```

### `interests[]`
```yaml
name, summary
```

### `testimonials[]`
```yaml
name, role, category, url, email, quote
```

Sort testimonials by `quote` ascending (shortest first) before rendering.

### `references[]`  (contact-only, no quote)
```yaml
name, role, category, email, url
```

---

## Schema validation

Before rendering, validate the parsed YAML data against the JSON Schema file using the language's JSON Schema validation library (see each language plan for the specific library). The schema declares `"$schema": "https://json-schema.org/draft/2020-12/schema"`.

Resolve the schema path in this order: (1) `--schema` flag if explicitly set; (2) the `$schema` field in the YAML, resolved relative to the input file's directory; (3) `schema.json` in the same directory as the input file.

Validation failures should be printed to stderr with a descriptive message and cause the program to exit with a non-zero status. Validation is a pre-flight check - a failed validation does not necessarily mean the data is unrenderable, so consider using `--skip-validation` as an escape hatch.

---

## Employer grouping algorithm

This is the most important non-trivial logic. The `work[]` array is already in reverse-chronological order (most recent first).

**Step 1 - Assign a group key to each entry:**
- If `employerGroup` is present, use it as the group key.
- Otherwise, use the `employer` name as the group key.

**Step 2 - Group consecutive entries sharing the same group key** into an `EmployerGroup`. Non-consecutive entries with the same key are treated as separate groups (two separate stints at the same company).

**Step 3 - Compute group-level metadata:**
- `displayName`: the `employer` of the first (most recent) entry in the group.
- `url`: the `url` of the first entry.
- `startDate`: the earliest `startDate` among all entries in the group.
- `endDate`: the latest `endDate` among all entries (empty string / absent = "Present").
- `formerNames`: any distinct `employer` values beyond the first, in display order - used to note name changes, e.g. "(formerly Rackspace Technology)".

---

## Output: `docs/index.html`

The output is a self-contained HTML file. Read the existing `docs/index.html` in full to understand all CSS custom properties, class names, and section structure. The new implementations must produce equivalent output with the additions below.

### Header: title label

Render `basics.label` as a `.title-label` element directly above the `.name` heading:

```html
<p class="title-label">Senior Software Engineer</p>
<h1 class="name">Stephen Brown II</h1>
```

The `.title-label` class is already defined in `docs/index.html` (small caps, accent color, letter-spaced).

### Font customization

The name at the top of the resume uses a display/serif font specified at render time via the `--name-font` CLI flag. The rest of the page uses Inter (body, labels, dates, tags, etc.).

**`--name-font` flag** accepts a Google Fonts family name exactly as it appears in the Google Fonts catalog (e.g. `"Instrument Serif"`, `"Playfair Display"`, `"EB Garamond"`). Default: `"Instrument Serif"`.

The program must:
1. Convert the font name to a URL-safe form by replacing spaces with `+` (e.g. `"Playfair Display"` → `"Playfair+Display"`).
2. Generate a Google Fonts `<link>` tag for that family, requesting both upright and italic styles. Use the format:
   ```
   https://fonts.googleapis.com/css2?family={URL_NAME}:ital@0;1&display=swap
   ```
   This works for most serif families; if the font requires weight axes, that is acceptable to leave for manual adjustment.
3. Inject the font name into the CSS `:root` as `--name-font`:
   ```css
   :root {
     --name-font: 'Instrument Serif', Georgia, serif;
     /* other vars unchanged */
   }
   ```
   The fallback stack after the chosen font name should always be `Georgia, serif`.
4. Apply `--name-font` to the `.name` selector instead of `--serif`:
   ```css
   .name {
     font-family: var(--name-font);
     /* all other .name properties unchanged */
   }
   ```

The `--serif` CSS variable (if kept) may remain for any other uses of a serif font in the template, but the name heading must use `--name-font`.

### Tabular numbers for years

Inter supports the `tnum` OpenType feature. Apply it to elements that display year ranges so columns align visually in the skills section and employment dates. Add to the `<style>` block:

```css
.job-dates,
.skill-item,
.edu-detail {
  font-variant-numeric: tabular-nums;
}
```

### Character disambiguation for Inter

Inter includes OpenType stylistic sets that improve legibility by disambiguating similar characters. Enable them on `body` alongside antialiasing:

```css
body {
  -webkit-font-smoothing: antialiased;
  font-feature-settings: 'dlig' 1, 'calt' 1, 'ss01' 1, 'ss04' 1, 'ss07' 1;
}
```

- `dlig` - discretionary ligatures
- `calt` - contextual alternates
- `ss01` - alternate digit one (serif base, distinct from lowercase L)
- `ss04` - disambiguation, no slashed zero (open zero; conflicts with `ss02`)
- `ss07` - open digit seven

### New employer-group HTML structure

When a group has more than one position, wrap the positions in an `employer-group` div instead of rendering them as bare `job` divs:

```html
<div class="employer-group">
  <div class="employer-header">
    <div>
      <span class="employer-name"><a href="https://jumpcloud.com">JumpCloud, Inc.</a></span>
      <!-- only rendered if former names exist: -->
      <span class="employer-former">(formerly OldName)</span>
    </div>
    <span class="job-dates">May 2021 – May 2026</span>
  </div>

  <div class="job">
    <div class="job-header">
      <span class="job-title">Senior Software Engineer</span>
      <span class="job-dates">Oct 2022 – May 2026</span>
    </div>
    <div class="job-meta">Remote</div>
    <ul class="highlights">
      <li>...</li>
    </ul>
    <div class="tags">
      <span class="tag">Go</span>
    </div>
  </div>

  <div class="position-divider"></div>

  <div class="job">
    <div class="job-header">
      <span class="job-title">Software Engineer 3</span>
      <span class="job-dates">May 2021 – Oct 2022</span>
    </div>
    <div class="job-meta">Longmont, CO</div>
    <ul class="highlights">
      <li>...</li>
    </ul>
    <div class="tags">
      <span class="tag">Go</span>
    </div>
  </div>
</div>
```

When a group has exactly one position, render it as the existing bare `job` div (employer name shown in `.job-meta`, no wrapping group element):

```html
<div class="job">
  <div class="job-header">
    <span class="job-title">Software Developer</span>
    <span class="job-dates">Oct 2019 – May 2021</span>
  </div>
  <div class="job-meta"><a href="https://objectrocket.com">ObjectRocket</a> &middot; Remote</div>
  <ul class="highlights">
    <li>...</li>
  </ul>
  <div class="tags">
    <span class="tag">Go</span>
  </div>
</div>
```

Between top-level items (employer groups and lone jobs alike) insert `<hr class="job-divider">` - but not after the last item.

### Section structure with `section-intro`

Every `<section>` wraps its `.section-label` and first content block in `<div class="section-intro">` to keep them together across page breaks. For sections with a single content container (Profile, Projects, Skills, Education), the entire content is inside `section-intro`. For sections with multiple content items (Experience, References), only the label and first item go inside `section-intro`; remaining items follow as siblings:

```html
<section>
  <div class="section-intro">
    <div class="section-label">Experience</div>
    <!-- first employer-group or job -->
  </div>
  <hr class="job-divider">
  <!-- remaining groups/jobs -->
</section>
```

### New and modified CSS

Add the following inside the `<style>` block (in addition to all existing CSS from `docs/index.html`):

```css
/* ── Name font (set dynamically by --name-font flag) ── */
:root {
  --name-font: 'Instrument Serif', Georgia, serif; /* overridden by renderer */
}

.name { font-family: var(--name-font); }

/* ── Tabular numbers for date alignment ── */
.job-dates,
.skill-item,
.edu-detail { font-variant-numeric: tabular-nums; }

/* ── Employer group (multiple positions) ──── */
.employer-group { margin-bottom: 0; }

.employer-header {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: baseline;
  gap: 8px;
  padding-bottom: 6px;
  margin-bottom: 8px;
  border-bottom: 1px solid var(--rule);
}

.employer-name {
  font-size: 0.87rem;
  font-weight: 700;
  color: var(--black);
}
.employer-name a, .project-name a { color: inherit; }

.employer-former {
  font-size: 0.73rem;
  color: var(--muted);
  margin-left: 6px;
}

.employer-group .job {
  padding-left: 10px;
  border-left: 2px solid var(--rule);
  margin-bottom: 0;
}

.position-divider {
  border: none;
  border-top: 1px dashed var(--rule);
  margin: 10px 0 10px 10px;
}
```

Link styles cascade from the global `a { color: var(--muted); text-decoration: none; }` and `a:hover { color: var(--accent); }` rules. Do not add redundant per-component link overrides. The one exception is `.employer-name a { color: inherit; }` and `.project-name a { color: inherit; }` which intentionally override the muted default to use the parent element's black color.

The `<link>` block in `<head>` becomes dynamic - the Instrument Serif link is replaced by the font-specific link generated from `--name-font`, while the Inter link remains fixed:

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<!-- Generated from --name-font: -->
<link href="https://fonts.googleapis.com/css2?family={URL_NAME}:ital@0;1&display=swap" rel="stylesheet">
<!-- Fixed: -->
<link rel="stylesheet" href="https://rsms.me/inter/inter.css">
```

### Date formatting

Dates in the YAML are ISO 8601 strings: `"YYYY-MM-DD"`, `"YYYY-MM"`, or `"YYYY"`. Format them as `"Mon YYYY"` for display (e.g. `"2022-10-03"` → `"Oct 2022"`). An absent or empty `endDate` should display as `"Present"`.

### Skills level CSS class

When rendering `skills.list`, apply a modifier class on `.skill-level` based on the level string:
- `"Advanced"` → `class="skill-level adv"`
- `"Intermediate"` → `class="skill-level mid"`
- Anything else (Familiar, Beginner) → `class="skill-level"`

The skills section iterates `skills.sets` for the group structure and looks up each skill name in `skills.list` to get the level for display.

### Skill sort order

Within each domain, sort skills by proficiency descending, then name ascending:
`Advanced`(3) > `Intermediate`(2) > `Familiar`(1) > `Beginner`(0). Apply before rendering; do not alter the source data order.

### Keyword shuffle

Shuffle `keywords` arrays for work entries and open source projects on each render. Skills are **not** shuffled — they are sorted as above.

### HTML escaping

All user-supplied string values must be HTML-escaped before insertion (`&`, `<`, `>`, `"`, `'`).

### Non-breaking spaces in summary

The summary paragraph uses `&nbsp;` to prevent short connector words from starting a line alone. The rule: **if the next word is 4 characters or shorter, replace the preceding space with `&nbsp;`** (i.e. the nbsp *precedes* the short word, binding it to its predecessor).

Example: `"I own problems"` → `"I&nbsp;own problems"` (nbsp before `"own"`, not after `"I"`).
Example: `"building and maintaining"` → `"building&nbsp;and maintaining"`.

This applies only to `basics.summary`. Implement as a template function (`nbspSummary`) that returns `template.HTML` to bypass auto-escaping of the `&nbsp;` entity.

---

## CLI interface

All implementations must accept:

| Flag / Arg | Default | Description |
|---|---|---|
| `--input` / `-i` | `../resume.yaml` | Path to input YAML file |
| `--output` / `-o` | `../docs/index.html` | Path to write HTML output (use `../docs/{lang}-index.html` when testing) |
| `--name-font` / `-f` | `Instrument Serif` | Google Fonts family name for the name heading |
| `--schema` | _(derived)_ | Path to JSON Schema file. Resolution order: (1) this flag if set, (2) `$schema` field in the YAML resolved relative to the input file's directory, (3) `schema.json` in the same directory as the input file |
| `--since` | _(none)_ | Exclude work entries whose `endDate` falls before this date. Accepts `YYYY`, `YYYY-MM`, or `YYYY-MM-DD`. Entries with no `endDate` (current role) are always included. |
| `--skip-validation` | false | Skip JSON Schema validation of the YAML |
| `--help` / `-h` | - | Print usage |

---

## Formatting and linting

Each language implementation must be formatted and lint-clean before merging. Every language adds two justfile recipes: `{lang}-fmt` (apply formatting in-place) and `{lang}-lint` (report issues, non-zero exit on failure). Both also run as pre-commit hooks.

| Language | Setup | Formatter | Linter |
|---|---|---|---|
| Go | `just go-setup` | `golangci-lint fmt` | `golangci-lint run` |
| Python | `uv sync` (ruff is a dev dep) | `ruff format` | `ruff check` |
| Rust | `rustup component add rustfmt clippy` | `cargo fmt` | `cargo clippy -- -D warnings` |
| Java | Maven plugins; no extra install | `mvn spotless:apply` | `mvn checkstyle:check` |
| Elixir | `mix deps.get` (credo is a dep) | `mix format` | `mix credo --strict` |

Each language adds a `{lang}-setup` justfile recipe that installs any tools not managed by the language's own package manager. The generic `setup` recipe calls all implemented language setups.

See individual language plans for tool configuration details.

---

## Justfile integration

Each language lives in its own subdirectory and adds three recipes to the repo-root `justfile`. Use the `[working-directory: '{lang}']` attribute instead of `cd` in the recipe body.

```just
[working-directory: '{lang}']
{lang}-build:
    <build command>

[working-directory: '{lang}']
{lang}-render: {lang}-build
    ./{binary} --input ../resume.yaml --output ../docs/index.html

[working-directory: '{lang}']
{lang}-validate: {lang}-build
    ./{binary} --input ../resume.yaml --output /dev/null

[working-directory: '{lang}']
{lang}-setup:
    <install tools not managed by the lang's package manager>

[working-directory: '{lang}']
{lang}-fmt:
    <format command>

[working-directory: '{lang}']
{lang}-lint:
    <lint command>
```

The generic recipes delegate to all implemented languages:

```just
build: {lang}-render          # extend as more langs are added

validate: {lang}-validate
```

When testing, pass `--output ../docs/{lang}-index.html` to avoid overwriting the canonical `docs/index.html`. The `dev` recipe (`build` + `serve`) and `watch` recipe both call `build`.

---

## What to preserve from `docs/index.html`

- All CSS custom properties and their values in the `:root` block (modified as described above)
- All `@media print` and `@media (max-width: 600px)` rules, including section-title orphan prevention:
  ```css
  .section-intro  { break-inside: avoid; page-break-inside: avoid; }
  .section-label { page-break-after: avoid; break-after: avoid; }
  ```
  `break-before/after: avoid` are unreliable hints in both Firefox and Chromium. The reliable fix is wrapping `.section-label` and its first content sibling in `<div class="section-intro">` and using `break-inside: avoid` on that container. Each section uses this wrapper so the label is always atomically bound to its content.
- Section order: header, summary, experience, skills, open source & projects, education & certifications, references (testimonials)
- The `print-only` span in the contact block containing the full resume URL
- The `&middot;` separator used in `.job-meta` and `.contact`
- The `.footer-grid` layout for education + certifications on one row
- Link styles cascade from global `a` / `a:hover` rules — no per-component overrides needed except `.employer-name a` and `.project-name a` which use `color: inherit` to show black instead of muted

---

## Testing

**During development, write to `docs/{lang}-index.html`** (e.g. `docs/go-index.html`) using `--output docs/go-index.html` so the canonical `docs/index.html` is not overwritten. Open the lang-prefixed file in a browser to review, then diff against the canonical to verify equivalence before promoting.

After generating output, verify:
- All sections are present and in correct order
- JumpCloud entries are grouped under a single employer header with the overall date range
- Dates format correctly and align with tabular numbers
- No raw YAML keys or unescaped HTML appear
- The name heading uses the specified font, loaded from Google Fonts
- `@media print` layout looks reasonable

The canonical `docs/index.html` is the structural reference. Diff the lang-prefixed output against it to verify equivalence.
