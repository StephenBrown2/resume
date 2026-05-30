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
```

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

### `references[]`  (contact-only, no quote)
```yaml
name, role, category, email, url
```

---

## Schema validation

Before rendering, validate the parsed YAML data against `schema.json` using the language's JSON Schema validation library (see each language plan for the specific library). The schema declares `"$schema": "https://json-schema.org/draft/2020-12/schema"`.

Validation failures should be printed to stderr with a descriptive message and cause the program to exit with a non-zero status. Validation is a pre-flight check — a failed validation does not necessarily mean the data is unrenderable, so consider using `--skip-validation` as an escape hatch.

---

## Employer grouping algorithm

This is the most important non-trivial logic. The `work[]` array is already in reverse-chronological order (most recent first).

**Step 1 — Assign a group key to each entry:**
- If `employerGroup` is present, use it as the group key.
- Otherwise, use the `employer` name as the group key.

**Step 2 — Group consecutive entries sharing the same group key** into an `EmployerGroup`. Non-consecutive entries with the same key are treated as separate groups (two separate stints at the same company).

**Step 3 — Compute group-level metadata:**
- `displayName`: the `employer` of the first (most recent) entry in the group.
- `url`: the `url` of the first entry.
- `startDate`: the earliest `startDate` among all entries in the group.
- `endDate`: the latest `endDate` among all entries (empty string / absent = "Present").
- `formerNames`: any distinct `employer` values beyond the first, in display order — used to note name changes, e.g. "(formerly Rackspace Technology)".

---

## Output: `docs/index.html`

The output is a self-contained HTML file. Read the existing `docs/index.html` in full to understand all CSS custom properties, class names, and section structure. The new implementations must produce equivalent output with the additions below.

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

### New employer-group HTML structure

When a group has **more than one position**, wrap the positions in an `employer-group` div instead of rendering them as bare `job` divs:

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

When a group has **exactly one position**, render it as the existing bare `job` div (employer name shown in `.job-meta`, no wrapping group element):

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

Between top-level items (employer groups and lone jobs alike) insert `<hr class="job-divider">` — but not after the last item.

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
.edu-detail {
  font-variant-numeric: tabular-nums;
}

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
.employer-name a { color: inherit; text-decoration: none; }
.employer-name a:hover { color: var(--accent); }

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

The `<link>` block in `<head>` becomes dynamic — the Instrument Serif link is replaced by the font-specific link generated from `--name-font`, while the Inter link remains fixed:

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

### HTML escaping

All user-supplied string values must be HTML-escaped before insertion (`&`, `<`, `>`, `"`, `'`).

### Non-breaking spaces in summary

The summary paragraph in the header section uses non-breaking spaces (`&nbsp;`) between short words and the following word to prevent awkward line breaks. The pattern: for any word 4 characters or shorter followed by a longer word, replace the space between them with `&nbsp;`. Example: `"I own problems"` → `"I&nbsp;own problems"`. This applies only to the `basics.summary` paragraph.

---

## CLI interface

All implementations must accept:

| Flag / Arg | Default | Description |
|---|---|---|
| `--input` / `-i` | `../resume.yaml` | Path to input YAML file |
| `--output` / `-o` | `../docs/index.html` | Path to write HTML output |
| `--name-font` / `-f` | `Instrument Serif` | Google Fonts family name for the name heading |
| `--skip-validation` | false | Skip JSON Schema validation of the YAML |
| `--help` / `-h` | — | Print usage |

---

## What to preserve from `docs/index.html`

- All CSS custom properties and their values in the `:root` block (modified as described above)
- All `@media print` and `@media (max-width: 600px)` rules
- Section order: header, summary, experience, open source & projects, skills, education & certifications, references (testimonials)
- The `print-only` span in the contact block containing the full resume URL
- The `&middot;` separator used in `.job-meta` and `.contact`
- The `.footer-grid` layout for education + certifications on one row

---

## Testing

After generating `docs/index.html`, open it in a browser and verify:
- All sections are present and in correct order
- JumpCloud entries are grouped under a single employer header with the overall date range
- Dates format correctly and align with tabular numbers
- No raw YAML keys or unescaped HTML appear
- The name heading uses the specified font, loaded from Google Fonts
- `@media print` layout looks reasonable

The existing `docs/index.html` is the structural reference. Diff the two if you need to verify equivalence.
