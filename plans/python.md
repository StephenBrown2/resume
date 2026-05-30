# Python implementation plan

**Read `plans/shared-context.md` first.** This document covers only Python-specific decisions.

---

## Goal

A Python CLI script in `python/` that reads `resume.yaml`, groups work entries by employer, and writes `docs/index.html`. Uses modern Python with dataclasses, type hints throughout, and Jinja2 for templating.

---

## Language version

- **Python 3.13** (latest stable as of mid-2025). No compatibility shims for older versions.
- Use `match` statements (structural pattern matching) where appropriate.
- Use `|` for union types in annotations (e.g. `str | None` instead of `Optional[str]`).

---

## Package management

Use **`uv`** (not pip/poetry/pipenv) — it is the current standard for Python project management. Initialize with `uv init --app resume-renderer` inside the `python/` directory.

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
requires-python = ">=3.13"
dependencies = [
    "pyyaml>=6.0",
    "jinja2>=3.1",
]

[project.scripts]
resume-renderer = "resume_renderer.main:main"
```

No other external packages. `pyyaml` for YAML; `jinja2` for HTML templating with auto-escaping.

---

## Data model (`model.py`)

Use `@dataclass` with `field(default_factory=list)` for list fields and `None` defaults for optional scalars. Define a `from_dict` classmethod on each class (or use a simple recursive loader function).

```python
from dataclasses import dataclass, field

@dataclass
class WorkEntry:
    employer: str
    position: str
    startDate: str
    endDate: str | None = None
    employerGroup: str | None = None
    url: str | None = None
    summary: str | None = None
    location: str | None = None
    highlights: list[str] = field(default_factory=list)
    keywords: list[str] = field(default_factory=list)

@dataclass
class SkillSet:
    name: str
    skills: list[str] = field(default_factory=list)

@dataclass
class SkillItem:
    name: str
    level: str
    summary: str | None = None
    years: int | None = None

@dataclass
class Skills:
    sets: list[SkillSet] = field(default_factory=list)
    list: list[SkillItem] = field(default_factory=list)

@dataclass
class Basics:
    name: str
    label: str = ""
    email: str = ""
    phone: str = ""
    url: str = ""
    summary: str = ""
    location: Location = field(default_factory=lambda: Location())
    profiles: list[Profile] = field(default_factory=list)

@dataclass
class Location:
    city: str = ""
    region: str = ""
    countryCode: str = ""

@dataclass
class Profile:
    network: str = ""
    username: str = ""
    url: str = ""

# ... Disposition, Relocation, Project, Certificate, Education,
#     Language, Interest, Testimonial, Reference
# (follow same pattern; all fields with sensible defaults)

@dataclass
class Resume:
    basics: Basics = field(default_factory=Basics)
    disposition: Disposition | None = None
    work: list[WorkEntry] = field(default_factory=list)
    projects: list[Project] = field(default_factory=list)
    skills: Skills = field(default_factory=Skills)
    certificates: list[Certificate] = field(default_factory=list)
    education: list[Education] = field(default_factory=list)
    languages: list[Language] = field(default_factory=list)
    interests: list[Interest] = field(default_factory=list)
    testimonials: list[Testimonial] = field(default_factory=list)
    references: list[Reference] = field(default_factory=list)
```

### Loading from dict

Write a `load_resume(data: dict) -> Resume` function that recursively constructs the dataclass tree from the raw `yaml.safe_load()` output. Handle missing keys with `.get()` and defaults. Avoid `dacite` or other third-party loaders.

```python
def load_resume(data: dict) -> Resume:
    basics_raw = data.get("basics", {})
    basics = Basics(
        name=basics_raw.get("name", ""),
        # ... etc
        location=Location(**basics_raw.get("location", {})),
        profiles=[Profile(**p) for p in basics_raw.get("profiles", [])],
    )
    work = [
        WorkEntry(**{k: v for k, v in entry.items()})
        for entry in data.get("work", [])
    ]
    # ... etc
    return Resume(basics=basics, work=work, ...)
```

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
    {# single position — render as bare job div with employer in meta #}
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

## `main.py`

```python
import argparse
import sys
from pathlib import Path
import yaml
from .model import load_resume
from .grouping import group_work
from .render import build_env

def main() -> None:
    parser = argparse.ArgumentParser(description="Render resume.yaml to HTML")
    parser.add_argument("-i", "--input",  default="../resume.yaml")
    parser.add_argument("-o", "--output", default="../docs/index.html")
    args = parser.parse_args()

    raw = yaml.safe_load(Path(args.input).read_text(encoding="utf-8"))
    resume = load_resume(raw)
    groups = group_work(resume.work)
    skill_map = {item.name: item for item in resume.skills.list}

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
    )

    out = Path(args.output)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(html, encoding="utf-8")
    print(f"wrote {out}", file=sys.stderr)

if __name__ == "__main__":
    main()
```

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
- `markupsafe.Markup` is part of Jinja2's dependency tree — no extra install needed.
- Use `__future__.annotations` (`from __future__ import annotations`) at the top of each file to enable PEP 563 deferred evaluation if you use forward references in type hints, though Python 3.13 handles most cases correctly without it.
- Do not use `pydantic` — it adds a significant dependency for no benefit over plain dataclasses here. The YAML is trusted internal data, not user input requiring validation.
- Jinja2's `autoescape` with `["html"]` will escape values interpolated with `{{ }}`. Pre-built `Markup` objects (from `nbsp_words`) are treated as already-safe and not double-escaped. Template literal text (`&middot;`, `&amp;`) is always passed through verbatim.
