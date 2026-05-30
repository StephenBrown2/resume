# Python implementation plan

**Read `plans/shared-context.md` first.** This document covers only Python-specific decisions.

---

## Goal

A Python CLI script in `python/` that reads `resume.yaml`, groups work entries by employer, and writes `docs/index.html`. Uses modern Python with dataclasses, type hints throughout, and Jinja2 for templating.

---

## Language version

- **Python 3.14.5** (latest stable). No compatibility shims for older versions.
- Use `match` statements (structural pattern matching) where appropriate.
- Use `|` for union types in annotations (e.g. `str | None` instead of `Optional[str]`).

---

## Package management

Use **`uv`** (not pip/poetry/pipenv) - it is the current standard for Python project management. Initialize with `uv init --app resume-renderer` inside the `python/` directory.

```
python/
  pyproject.toml
  uv.lock
  README.md          (optional)
  src/
    resume_renderer/
      __init__.py
      main.py          # CLI entry point
      model.py         # dataclasses for the YAML schema
      grouping.py      # employer grouping logic
      render.py        # Jinja2 setup and rendering
      filters.py       # custom Jinja2 filters
      template.html    # Jinja2 template
```

---

## Dependencies (`pyproject.toml`)

```toml
[project]
name = "resume-renderer"
version = "0.1.0"
requires-python = ">=3.14"
dependencies = [
    "pyyaml>=6.0",
    "jinja2>=3.1",
    "pydantic>=2.0",
    "jsonschema>=4.26",
]

[project.scripts]
resume-renderer = "resume_renderer.main:main"
```

- **`pyyaml`** - YAML parsing.
- **`jinja2`** - HTML templating with auto-escaping.
- **`pydantic`** - data model validation and coercion; replaces hand-written `from_dict` loaders. Use `model_validate` to construct models from the raw YAML dict.
- **`jsonschema`** 4.26.0 - JSON Schema Draft 2020-12 validation against `schema.json` before Pydantic model construction.

---

## Data model (`model.py`)

Use **Pydantic v2** `BaseModel` for all types. Pydantic handles optional fields, default values, and type coercion from the raw YAML dict automatically via `model_validate`. Use `model_config = ConfigDict(extra="allow")` on all models to match the `additionalProperties: true` policy of the schema.

```python
from pydantic import BaseModel, ConfigDict

class WorkEntry(BaseModel):
    model_config = ConfigDict(extra="allow")
    employer: str
    position: str
    startDate: str
    endDate: str | None = None
    employerGroup: str | None = None
    url: str | None = None
    summary: str | None = None
    location: str | None = None
    highlights: list[str] = []
    keywords: list[str] = []

class SkillSet(BaseModel):
    model_config = ConfigDict(extra="allow")
    name: str
    skills: list[str] = []

class SkillItem(BaseModel):
    model_config = ConfigDict(extra="allow")
    name: str
    level: str
    summary: str | None = None
    years: int | None = None

class Skills(BaseModel):
    model_config = ConfigDict(extra="allow")
    sets: list[SkillSet] = []
    list: list[SkillItem] = []

class Location(BaseModel):
    model_config = ConfigDict(extra="allow")
    city: str = ""
    region: str = ""
    countryCode: str = ""

class Profile(BaseModel):
    model_config = ConfigDict(extra="allow")
    network: str = ""
    username: str = ""
    url: str = ""

class Basics(BaseModel):
    model_config = ConfigDict(extra="allow")
    name: str
    label: str = ""
    email: str = ""
    phone: str = ""
    url: str = ""
    summary: str = ""
    location: Location = Location()
    profiles: list[Profile] = []

# ... Disposition, Relocation, Project, Certificate, Education,
#     Language, Interest, Testimonial, Reference
# (follow same pattern)

class Resume(BaseModel):
    model_config = ConfigDict(extra="allow")
    basics: Basics
    disposition: Disposition | None = None
    work: list[WorkEntry] = []
    projects: list[Project] = []
    skills: Skills = Skills()
    certificates: list[Certificate] = []
    education: list[Education] = []
    languages: list[Language] = []
    interests: list[Interest] = []
    testimonials: list[Testimonial] = []
    references: list[Reference] = []
```

### Loading

```python
raw = yaml.safe_load(Path(args.input).read_text(encoding="utf-8"))
resume = Resume.model_validate(raw)
```

No manual `from_dict` function needed - Pydantic's `model_validate` recurses into nested models automatically.

---

## Employer grouping (`grouping.py`)

```python
from dataclasses import dataclass, field

@dataclass
class EmployerGroup:
    display_name: str
    former_names: list[str]
    url: str | None
    start_date: str        # earliest in group
    end_date: str | None   # None means "Present"
    positions: list[WorkEntry]

def group_work(entries: list[WorkEntry]) -> list[EmployerGroup]:
    groups: list[EmployerGroup] = []
    for entry in entries:
        key = entry.employerGroup or entry.employer
        if groups and _group_key(groups[-1]) == key:
            _extend_group(groups[-1], entry)
        else:
            groups.append(_new_group(entry))
    return groups

def _group_key(g: EmployerGroup) -> str:
    # Return the original key used to start this group
    # Store as a private field or infer from positions[0]
    return (g.positions[0].employerGroup or g.positions[0].employer)

def _extend_group(g: EmployerGroup, entry: WorkEntry) -> None:
    g.positions.append(entry)
    # Update start_date if earlier
    if entry.startDate < g.start_date:
        g.start_date = entry.startDate
    # Update former_names
    if entry.employer != g.display_name and entry.employer not in g.former_names:
        g.former_names.append(entry.employer)

def _new_group(entry: WorkEntry) -> EmployerGroup:
    return EmployerGroup(
        display_name=entry.employer,
        former_names=[],
        url=entry.url,
        start_date=entry.startDate,
        end_date=entry.endDate,
        positions=[entry],
    )
```

---

## Jinja2 setup and custom filters (`render.py`, `filters.py`)

```python
from jinja2 import Environment, FileSystemLoader, select_autoescape

def build_env(template_dir: Path) -> Environment:
    env = Environment(
        loader=FileSystemLoader(template_dir),
        autoescape=select_autoescape(["html"]),
    )
    env.filters["format_date"] = format_date
    env.filters["nbsp_words"]  = nbsp_words
    env.filters["level_class"] = level_class
    return env
```

### `format_date` filter (`filters.py`)

```python
MONTHS = ["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"]

def format_date(iso: str | None) -> str:
    if not iso:
        return "Present"
    parts = iso.split("-")
    if len(parts) >= 2:
        month_idx = int(parts[1]) - 1
        return f"{MONTHS[month_idx]} {parts[0]}"
    return parts[0]
```

### `nbsp_words` filter (`filters.py`)

```python
from markupsafe import Markup  # bundled with Jinja2

def nbsp_words(s: str) -> Markup:
    """Replace spaces after short words (≤4 chars) with &nbsp; to prevent bad wraps."""
    words = s.split(" ")
    result = []
    for i, word in enumerate(words):
        result.append(word)
        if i < len(words) - 1:
            next_word = words[i + 1]
            sep = "&nbsp;" if len(word) <= 4 and len(next_word) > 4 else " "
            result.append(sep)
    return Markup("".join(result))
```

`Markup` marks the string as safe so Jinja2 does not double-escape the `&nbsp;` entity.

### `level_class` filter

```python
def level_class(level: str) -> str:
    match level.lower():
        case "advanced":    return "adv"
        case "intermediate": return "mid"
        case _:             return ""
```

---

## Template (`template.html`)

Standard Jinja2 with auto-escaping. Key patterns:

```jinja2
{% for group in employer_groups %}
  {% if group.positions | length > 1 %}
    <div class="employer-group">
      <div class="employer-header">
        <div>
          <span class="employer-name">
            {% if group.url %}<a href="{{ group.url }}">{% endif %}
            {{ group.display_name }}
            {% if group.url %}</a>{% endif %}
          </span>
          {% if group.former_names %}
          <span class="employer-former">
            (formerly {{ group.former_names | join(", ") }})
          </span>
          {% endif %}
        </div>
        <span class="job-dates">
          {{ group.start_date | format_date }} – {{ group.end_date | format_date }}
        </span>
      </div>
      {% for pos in group.positions %}
        {% include "_job_position.html" %}
        {% if not loop.last %}<div class="position-divider"></div>{% endif %}
      {% endfor %}
    </div>
  {% else %}
    {# single position - render as bare job div with employer in meta #}
    {% set pos = group.positions[0] %}
    {% set show_employer = true %}
    {% include "_job_position.html" %}
  {% endif %}
  {% if not loop.last %}<hr class="job-divider">{% endif %}
{% endfor %}
```

Split `_job_position.html` out as a partial template in the same directory to avoid repetition between the grouped and single-employer cases.

For skills, build `skill_map` in Python before rendering:

```python
skill_map = {item.name: item for item in resume.skills.list}
```

Pass it to the template context. In template: `{{ skill_map[skill_name].level | level_class }}`.

---

## Schema validation (`main.py`)

Before constructing Pydantic models, validate the raw dict against `schema.json` using `jsonschema`:

```python
import jsonschema
import json

def validate_schema(raw: dict, schema_path: str) -> None:
    schema = json.loads(Path(schema_path).read_text())
    validator = jsonschema.Draft202012Validator(schema)
    errors = list(validator.iter_errors(raw))
    if errors:
        for e in errors:
            print(f"validation error: {e.json_path}: {e.message}", file=sys.stderr)
        sys.exit(1)
```

Call before `Resume.model_validate(raw)`. Pydantic validation catches structural issues too; jsonschema catches schema-declared constraints (e.g. URI format, integer ranges) that Pydantic won't.

## `--name-font` flag

```python
parser.add_argument("-f", "--name-font", default="Instrument Serif",
                    help="Google Fonts family name for the name heading")
```

In `main()`, compute the font link and CSS value:

```python
font_url = args.name_font.replace(" ", "+")
google_fonts_link = (
    f'<link href="https://fonts.googleapis.com/css2?family={font_url}:ital@0;1'
    f'&amp;display=swap" rel="stylesheet">'
)
name_font_css = f"'{args.name_font}', Georgia, serif"
```

Pass both to the template context:

```python
html = tmpl.render(
    google_fonts_link=Markup(google_fonts_link),  # Markup = pre-escaped, safe
    name_font_css=name_font_css,
    # ... rest of context
)
```

In the template:

```jinja2
{{ google_fonts_link }}   {# already Markup, not double-escaped #}
...
--name-font: {{ name_font_css }};
```

## `main.py`

```python
import argparse
import sys
from pathlib import Path
import yaml
from .model import Resume
from .grouping import group_work
from .render import build_env

def main() -> None:
    parser = argparse.ArgumentParser(description="Render resume.yaml to HTML")
    parser.add_argument("-i", "--input",  default="../resume.yaml")
    parser.add_argument("-o", "--output", default="../docs/index.html")
    parser.add_argument("-f", "--name-font", default="Instrument Serif",
                        help="Google Fonts family name for the name heading")
    parser.add_argument("--skip-validation", action="store_true")
    args = parser.parse_args()

    raw = yaml.safe_load(Path(args.input).read_text(encoding="utf-8"))
    if not args.skip_validation:
        validate_schema(raw, "schema.json")
    resume = Resume.model_validate(raw)
    groups = group_work(resume.work)
    skill_map = {item.name: item for item in resume.skills.list}

    font_url = args.name_font.replace(" ", "+")
    google_fonts_link = Markup(
        f'<link href="https://fonts.googleapis.com/css2?family={font_url}:ital@0;1'
        f'&amp;display=swap" rel="stylesheet">'
    )
    name_font_css = f"'{args.name_font}', Georgia, serif"

    template_dir = Path(__file__).parent
    env = build_env(template_dir)
    tmpl = env.get_template("template.html")

    html = tmpl.render(
        basics=resume.basics,
        employer_groups=groups,
        projects=resume.projects,
        skill_sets=resume.skills.sets,
        skill_map=skill_map,
        certificates=resume.certificates,
        education=resume.education,
        languages=resume.languages,
        interests=resume.interests,
        testimonials=resume.testimonials,
        google_fonts_link=Markup(google_fonts_link),
        name_font_css=name_font_css,
    )

    out = Path(args.output)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(html, encoding="utf-8")
    print(f"wrote {out}", file=sys.stderr)

if __name__ == "__main__":
    main()
```

---

## `python/README.md`

Create `python/README.md` documenting this implementation. It should cover:

- **Prerequisites:** Python 3.14+ and [uv](https://docs.astral.sh/uv/). Install uv via `curl -LsSf https://astral.sh/uv/install.sh | sh` or `brew install uv`.
- **Setup:** `uv sync` (installs dependencies into a local `.venv`).
- **Run:** `uv run resume-renderer [flags]` (or `just python-render` from the repo root).
- **Flags:** table matching the CLI interface in `shared-context.md` (`--input`, `--output`, `--name-font`, `--skip-validation`).
- **Output:** writes `docs/index.html` (relative to the repo root when using the default path).
- **Package layout:** note the `src/resume_renderer/` structure and the `resume-renderer` entry point defined in `pyproject.toml`.

---

## Build and run

```sh
cd python
uv sync
uv run resume-renderer --input ../resume.yaml --output ../docs/index.html
```

Or after installing the project:
```sh
uv pip install -e python/
resume-renderer --input resume.yaml --output docs/index.html
```

Add to repo `justfile`:

```just
python-render:
    cd python && uv run resume-renderer --input ../resume.yaml --output ../docs/index.html
```

---

## Notes

- `yaml.safe_load` is sufficient; no need for `ruamel.yaml` unless round-trip preservation is required (it is not here).
- `markupsafe.Markup` is part of Jinja2's dependency tree - no extra install needed.
- Use `__future__.annotations` (`from __future__ import annotations`) at the top of each file to enable PEP 563 deferred evaluation if you use forward references in type hints, though Python 3.13 handles most cases correctly without it.
- Jinja2's `autoescape` with `["html"]` will escape values interpolated with `{{ }}`. Pre-built `Markup` objects (from `nbsp_words`) are treated as already-safe and not double-escaped. Template literal text (`&middot;`, `&amp;`) is always passed through verbatim.
