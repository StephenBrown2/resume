# Rust implementation plan

**Read `plans/shared-context.md` first.** This document covers only Rust-specific decisions.

---

## Goal

A Rust binary in `rust/` that reads `resume.yaml`, groups work entries by employer, and writes `docs/index.html`. Should be idiomatic Rust 2024 edition with clear error propagation.

---

## Language version and edition

- **Rust 1.96.0** (latest stable). Declare `edition = "2024"` in `Cargo.toml`.
- Minimum supported Rust version (MSRV) is not a concern - target latest stable.

---

## Dependencies (`Cargo.toml`)

```toml
[dependencies]
serde       = { version = "1", features = ["derive"] }
serde_yaml  = "0.9"
serde_json  = "1"        # required by jsonschema crate for input type
minijinja   = "2"        # Jinja2-compatible templating with auto-escaping
clap        = { version = "4", features = ["derive"] }
thiserror   = "2"
jsonschema  = "0.46"     # JSON Schema Draft 2020-12 validation

[dev-dependencies]
# none required
```

**`minijinja`** is preferred over `askama` here because the template can be edited without a recompile, matches Jinja2 syntax closely, and has first-class HTML auto-escaping. Use `minijinja::Environment` with auto-escaping enabled.

**`jsonschema`** (crate by Stranger6667) supports Draft 2020-12, which matches the `$schema` declared in `schema.json`.

---

## File structure

```
rust/
  Cargo.toml
  Cargo.lock
  src/
    main.rs       # CLI entry, error handling, orchestration
    model.rs      # serde-annotated structs for the YAML schema
    grouping.rs   # employer grouping logic
    render.rs     # minijinja setup, template filters, HTML generation
    template.html # embedded via include_str!()
```

---

## Data model (`src/model.rs`)

Use `#[derive(Debug, Deserialize)]` on all structs. Optional fields use `Option<T>`.

```rust
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Resume {
    pub basics:       Basics,
    pub disposition:  Option<Disposition>,
    pub work:         Vec<WorkEntry>,
    pub projects:     Vec<Project>,
    pub skills:       Skills,
    pub certificates: Vec<Certificate>,
    pub education:    Vec<Education>,
    pub languages:    Vec<Language>,
    pub interests:    Vec<Interest>,
    pub testimonials: Vec<Testimonial>,
    pub references:   Vec<Reference>,
}

#[derive(Debug, Deserialize)]
pub struct WorkEntry {
    pub employer:       String,
    #[serde(rename = "employerGroup")]
    pub employer_group: Option<String>,
    pub position:       String,
    pub url:            Option<String>,
    #[serde(rename = "startDate")]
    pub start_date:     String,
    #[serde(rename = "endDate")]
    pub end_date:       Option<String>,
    pub summary:        Option<String>,
    pub location:       Option<String>,
    pub highlights:     Vec<String>,
    pub keywords:       Vec<String>,
}

#[derive(Debug, Deserialize)]
pub struct Skills {
    pub sets: Vec<SkillSet>,
    pub list: Vec<SkillItem>,
}

#[derive(Debug, Deserialize)]
pub struct SkillSet {
    pub name:   String,
    pub skills: Vec<String>,
}

#[derive(Debug, Deserialize)]
pub struct SkillItem {
    pub name:    String,
    pub level:   String,
    pub summary: Option<String>,
    pub years:   Option<u32>,
}

// ... Basics, Location, Profile, Disposition, Relocation, Project,
//     Certificate, Education, Language, Interest, Testimonial, Reference
// (define all; see schema.json)
```

---

## Error handling (`src/main.rs`)

Use `thiserror` to define a single top-level `Error` enum:

```rust
#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("I/O error: {0}")]
    Io(#[from] std::io::Error),
    #[error("YAML parse error: {0}")]
    Yaml(#[from] serde_yaml::Error),
    #[error("template error: {0}")]
    Template(#[from] minijinja::Error),
    #[error("schema validation failed")]
    Validation,
}

pub type Result<T> = std::result::Result<T, Error>;
```

Propagate with `?` throughout; the `main` function returns `Result<()>`.

---

## CLI (`src/main.rs`)

```rust
use clap::Parser;

#[derive(Parser)]
#[command(about = "Renders resume.yaml to HTML")]
struct Args {
    #[arg(short, long, default_value = "../resume.yaml")]
    input: std::path::PathBuf,

    #[arg(short, long, default_value = "../docs/index.html")]
    output: std::path::PathBuf,

    #[arg(short = 'f', long, default_value = "Instrument Serif",
          help = "Google Fonts family name for the name heading")]
    name_font: String,

    #[arg(long, help = "Skip JSON Schema validation")]
    skip_validation: bool,
}
```

## Schema validation (`src/main.rs`)

After parsing the YAML file content to a `serde_json::Value` (convert via `serde_yaml` → `serde_json`), validate against `schema.json`:

```rust
use jsonschema::validator_for;

fn validate(schema_path: &str, data: &serde_json::Value) -> Result<()> {
    let schema_str = std::fs::read_to_string(schema_path)?;
    let schema: serde_json::Value = serde_json::from_str(&schema_str)?;
    let validator = validator_for(&schema)?;
    let result = validator.validate(data);
    if let Err(errors) = result {
        for e in errors {
            eprintln!("validation error: {e}");
        }
        return Err(Error::Validation);
    }
    Ok(())
}
```

Add `serde_json = "1"` to `Cargo.toml` dependencies (needed for the `jsonschema` crate's input type). The pipeline is:

1. Read YAML bytes.
2. Deserialize to `serde_yaml::Value`, then convert to `serde_json::Value` via `serde_json::to_value(serde_yaml::from_slice(bytes)?)`.
3. Validate with jsonschema.
4. Deserialize to typed `Resume` struct via `serde_yaml::from_slice(bytes)`.

## `--name-font` handling (`src/render.rs`)

```rust
pub fn font_link(name: &str) -> String {
    let url_name = name.replace(' ', "+");
    format!(
        r#"<link href="https://fonts.googleapis.com/css2?family={url_name}:ital@0;1&display=swap" rel="stylesheet">"#
    )
}

pub fn name_font_css(name: &str) -> String {
    format!("'{name}', Georgia, serif")
}
```

Add `font_link` and `name_font_css` outputs to `TemplateData`:

```rust
#[derive(serde::Serialize)]
pub struct TemplateData {
    pub google_fonts_link: String,  // raw HTML <link> tag
    pub name_font_css:     String,  // CSS value for --name-font
    // ... rest of fields
}
```

In the template, insert them unescaped:

```jinja2
{{ google_fonts_link | safe }}
```

and in the CSS:

```jinja2
--name-font: {{ name_font_css | safe }};
```

---

## Employer grouping (`src/grouping.rs`)

```rust
#[derive(Debug)]
pub struct EmployerGroup {
    pub display_name:  String,
    pub former_names:  Vec<String>,  // distinct names beyond display_name, in order seen
    pub url:           Option<String>,
    pub start_date:    String,       // earliest in group
    pub end_date:      Option<String>, // latest (None = "Present")
    pub positions:     Vec<WorkEntry>,
}

pub fn group_work(entries: Vec<WorkEntry>) -> Vec<EmployerGroup> {
    // Algorithm from shared-context.md:
    // - key = entry.employer_group.as_deref().unwrap_or(&entry.employer)
    // - group consecutive same-key entries
    // - track former_names as Vec of distinct employer strings beyond the first seen
}
```

The function takes ownership of `entries` and returns groups in the same (reverse-chronological) order.

---

## Date formatting (`src/render.rs`)

```rust
pub fn format_date(iso: &str) -> String {
    // parse "YYYY-MM-DD", "YYYY-MM", "YYYY"
    // return "Mon YYYY" (e.g. "Oct 2022") or "YYYY" for year-only
    // empty string → "Present"
}
```

Use the standard library only - parse the date string manually by splitting on `-`. Map month number to `["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"]`.

---

## Non-breaking spaces (`src/render.rs`)

```rust
pub fn nbsp_short_words(s: &str) -> String {
    // For each word of ≤4 chars followed by a longer word,
    // replace the separating space with &nbsp;
    // Returns the modified string with HTML entity inserted
}
```

Register as a custom Jinja filter: `env.add_filter("nbsp_words", ...)`.

---

## Template and rendering (`src/render.rs`)

```rust
use minijinja::{Environment, context};

pub fn render(data: TemplateData) -> Result<String> {
    let mut env = Environment::new();
    env.set_auto_escape_callback(|_| minijinja::AutoEscape::Html);
    env.add_template("resume", include_str!("template.html"))?;
    env.add_filter("format_date", format_date_filter);
    env.add_filter("nbsp_words", nbsp_words_filter);
    env.add_filter("level_class", level_class_filter);

    let tmpl = env.get_template("resume")?;
    Ok(tmpl.render(context!(data => data))?)
}
```

Pass a `TemplateData` struct that is `serde::Serialize` (derive it), so minijinja can serialize it into its value system:

```rust
#[derive(serde::Serialize)]
pub struct TemplateData {
    pub basics:          Basics,
    pub employer_groups: Vec<EmployerGroup>,
    pub projects:        Vec<Project>,
    pub skill_sets:      Vec<SkillSet>,
    pub skill_list:      Vec<SkillItem>,
    pub certificates:    Vec<Certificate>,
    pub education:       Vec<Education>,
    pub languages:       Vec<Language>,
    pub interests:       Vec<Interest>,
    pub testimonials:    Vec<Testimonial>,
}
```

Note: `minijinja` serializes via `serde::Serialize`, so all template data types need `#[derive(Serialize)]` in addition to `Deserialize`.

---

## Template (`src/template.html`)

The Jinja2 template mirrors the structure in `docs/index.html`. Key iteration patterns:

```jinja2
{% for group in employer_groups %}
  {% if group.positions | length > 1 %}
    <div class="employer-group">
      <div class="employer-header">
        ...{{ group.start_date | format_date }} – {{ group.end_date | format_date }}...
      </div>
      {% for pos in group.positions %}
        <div class="job">...</div>
        {% if not loop.last %}<div class="position-divider"></div>{% endif %}
      {% endfor %}
    </div>
  {% else %}
    <div class="job">
      ...
      <div class="job-meta">
        <a href="{{ group.url }}">{{ group.display_name }}</a> &middot; {{ group.positions[0].location }}
      </div>
      ...
    </div>
  {% endif %}
  {% if not loop.last %}<hr class="job-divider">{% endif %}
{% endfor %}
```

For the skills section, look up level per skill name:

```jinja2
{% for set in skill_sets %}
  <div>
    <div class="skill-group-label">{{ set.name }}</div>
    {% for skill_name in set.skills %}
      {% set item = skill_list | selectattr("name", "equalto", skill_name) | first %}
      <div class="skill-item">
        <span class="skill-name">{{ skill_name }}</span>
        <span class="skill-level {{ item.level | level_class }}">{{ item.level }}</span>
      </div>
    {% endfor %}
  </div>
{% endfor %}
```

---

## `rust/README.md`

Create `rust/README.md` documenting this implementation. It should cover:

- **Prerequisites:** Rust 1.96.0+ (edition 2024). Install via [rustup](https://rustup.rs/).
- **Build:** `cargo build --release` (or `just rust-build` from the repo root).
- **Run:** `./target/release/resume-renderer [flags]` - list all flags with `--help`.
- **Flags:** table matching the CLI interface in `shared-context.md` (`--input`, `--output`, `--name-font`, `--skip-validation`).
- **Output:** writes `docs/index.html` (relative to the repo root when using the default path).
- **Template:** note that `src/template.html` is embedded at compile time via `include_str!()` - the binary is fully self-contained.

---

## Build and run

```sh
cd rust
cargo build --release
./target/release/resume-renderer --input ../resume.yaml --output ../docs/index.html
```

Add to the repo `justfile`:

```just
rust-build:
    cd rust && cargo build --release

rust-render: rust-build
    rust/target/release/resume-renderer --input resume.yaml --output docs/index.html
```

---

## Notes

- Use `serde_yaml 0.9` which wraps `libyaml`; it handles the YAML multi-line strings in the resume correctly.
- `minijinja` auto-escape in Html mode will escape `&`, `<`, `>`, `"` in template variables. HTML entities written literally in the template (`&middot;`, `&nbsp;`) are passed through unescaped because they are template text, not values - this is correct.
- The `include_str!()` macro embeds the template at compile time; the binary is fully self-contained.
- For the `selectattr` filter in the skills template: minijinja 2.x supports it natively. If it is unavailable, pre-build a `skill_map: HashMap<String, SkillItem>` in `TemplateData` and use `skill_map[skill_name]` in the template instead.
