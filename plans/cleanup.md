# Post-implementation cleanup plan

Execute this plan **after all five language implementations** (Go, Rust, Python, Elixir, Java) are complete and each produces correct output from `resume.yaml` → `docs/index.html`.

---

## Files to delete

### Root-level legacy Go schema files

```
fresh-resume.go
json-resume.go
```

These are schema type definitions (`FRESHResume`, `JSONResume`) that were written for the old `goresume` CLI tool. They have no `main` package or entry point and are not referenced by anything in the repo. The Go implementation plan explicitly calls for their removal. Confirm with `grep -r "FRESHResume\|JSONResume" .` that nothing references them before deleting.

### `resume.html`

```
resume.html
```

A stale HTML snapshot at the repo root, distinct from `docs/index.html`. It appears to predate the current build pipeline. Once `docs/index.html` is being generated correctly by a renderer, this file has no purpose.

### `Variable Web Fonts/` (root-level copy)

```
Variable Web Fonts/
```

There are two copies of this directory: one at the repo root and one inside `docs/`. The one at the repo root is not served and not referenced. Confirm the `docs/` copy is the one in use (it's co-located with `docs/index.html`), then delete the root-level copy.

---

## Files to update

### `justfile`

The existing `justfile` drives the old `goresume`-based build. Replace its contents with the new recipes added by the language implementations. The old recipes to remove:

```just
build theme="block":
    goresume export --resume {{filename}} --html-theme {{theme}} --html-output docs/index.html

validate:
    goresume validate --resume {{filename}}
```

The `serve` and `watch` recipes can stay as-is (they are tool-agnostic). The `go` alias recipe (`go: build serve`) should be renamed or removed to avoid shadowing the new `go-render` recipe.

**Default recipe:** The bare `just` invocation should print available render recipes rather than silently run one language's build. Use Just's built-in `--list` for this:

```just
default:
    @just --list
```

**Working-directory syntax:** Any recipe that previously used `cd <dir> && ...` must be converted to use Just's `[working-directory]` attribute instead:

```just
[working-directory: 'go']
go-build:
    go build -o resume-renderer .

[working-directory: 'go']
go-render: go-build
    ./resume-renderer --input ../resume.yaml --output ../docs/index.html

[working-directory: 'rust']
rust-build:
    cargo build --release

[working-directory: 'rust']
rust-render: rust-build
    ./target/release/resume-renderer --input ../resume.yaml --output ../docs/index.html

[working-directory: 'python']
python-render:
    uv run resume-renderer --input ../resume.yaml --output ../docs/index.html

[working-directory: 'elixir']
elixir-build:
    mix escript.build

[working-directory: 'elixir']
elixir-render: elixir-build
    ./resume_renderer --input ../resume.yaml --output ../docs/index.html

[working-directory: 'java']
java-build:
    mvn -q package -DskipTests

java-render: java-build
    java -jar java/target/resume-renderer-1.0-SNAPSHOT.jar \
         --input resume.yaml --output docs/index.html
```

The `serve` and `watch` recipes operate from the repo root and need no working-directory attribute.

### `.gitignore`

Add ignore patterns for build artifacts produced by the new implementations:

```gitignore
# Go
go/resume-renderer

# Rust
rust/target/

# Python
python/.venv/
python/__pycache__/

# Elixir
elixir/_build/
elixir/deps/
elixir/resume_renderer

# Java
java/target/
```

### `README.md`

Update to describe the new build process. At minimum, remove any references to `goresume` and point to `just <lang>-render` (or the chosen default) as the way to regenerate `docs/index.html`.

### `.pre-commit-config.yaml`

The current hooks (`trailing-whitespace`, `end-of-file-fixer`, `check-json`, `check-yaml`, `check-added-large-files`) are all still appropriate. No changes needed unless new language-specific linters are desired.

---

## Verification checklist

Before committing the cleanup:

- [ ] `grep -r "goresume" .` returns no results (outside of git history)
- [ ] `grep -r "FRESHResume\|JSONResume" .` returns no results
- [ ] `just <lang>-render` (for the chosen primary language) produces `docs/index.html` successfully
- [ ] `docs/index.html` opens correctly in a browser
- [ ] `git status` shows only intentional deletions and modifications
- [ ] No language build artifacts are staged (`go/resume-renderer`, `rust/target/`, etc.)

---

## Future step: theme extensibility

**Do not touch `themes/` during cleanup.** The five HTML files there are kept as the starting point for a multi-theme system. They predate the new schema and are not yet wired into any renderer, but they represent reusable layout alternatives worth preserving.

### What needs to happen before themes are usable

1. **Update each theme to the new schema field names.** The existing themes were written against the old `goresume` / JSON Resume field layout. They reference field names and section structures that may not match the hybrid schema (e.g. `name` vs `employer`, `company` vs `employer`, `reference` vs `testimonial`/`references` split, `skills[].keywords` vs `skills.sets`/`skills.list`, etc.). Each theme file must be audited and updated to use the fields defined in `schema.json`.

2. **Port each theme to each language's template format.** Each renderer uses a different template engine:

   | Language | Template engine | Theme file extension |
   |---|---|---|
   | Go | `html/template` | `.gohtml` or `.html` |
   | Rust | Minijinja (Jinja2-compatible) | `.html.jinja` or `.html` |
   | Python | Jinja2 | `.html` |
   | Elixir | EEx | `.html.eex` |
   | Java | Pebble (Jinja2/Twig-style) | `.html` |

   Jinja2-syntax themes (Rust, Python, Java/Pebble) will be the most portable — a single theme file may work across all three with minor adjustments. Go and Elixir will require separate ports.

3. **Add a `--theme` flag to each renderer.** The flag should accept a theme name (e.g. `block`, `simple`, `actual`) and resolve it to the corresponding template file, replacing the single embedded template. The default theme should produce output equivalent to the current `docs/index.html`.

4. **Move themes into the renderer subdirectories** (or a shared `themes/` directory that each renderer reads from) once the template format question is resolved. The current `themes/` directory at the repo root can serve as the canonical source of truth until then.

5. **Add `just <lang>-render --theme <name>` recipes** (or a `theme` parameter) to the justfile once the flag exists in all renderers.
