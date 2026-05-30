# Post-implementation cleanup plan

Execute this plan **after** at least one language implementation is working and producing correct output from `resume.yaml` → `docs/index.html`.

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

### `themes/`

```
themes/actual.html
themes/block.html
themes/positive.html
themes/simple-compact.html
themes/simple.html
```

Template files for the old `goresume` tool. The new renderers embed their own template; these are unused. Delete the entire `themes/` directory.

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

After cleanup, the default recipe (run by bare `just`) should invoke whichever language implementation has been chosen as the primary build.

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
