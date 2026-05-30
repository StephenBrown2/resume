# Stephen Brown II - Resume Generator

Source repository for my resume. The canonical data lives in `resume.yaml`
(validated against `schema.json`); a renderer reads it and writes a
self-contained `docs/index.html`, which GitHub Pages serves at
<https://stephenbrown2.github.io/resume/>.

---

## Repository layout

```
resume.yaml          # resume data (hybrid JSON Resume / FRESH schema)
schema.json          # JSON Schema Draft 2020-12 for resume.yaml
docs/index.html      # generated output - committed and served by GitHub Pages
themes/              # alternate HTML layouts (not yet wired into any renderer)
plans/               # implementation plans for renderers in Go, Rust, Python, Elixir, and Java
```

---

## Current build

Install [goresume](https://github.com/nikaro/goresume) and
[just](https://just.systems/):

```shell
# with Go
go install github.com/nikaro/goresume@latest

# with Homebrew
brew install nikaro/tap/goresume just

# on ArchLinux
yay -S goresume-bin just
```

Generate `docs/index.html` and serve it locally:

```shell
just build    # export resume.yaml -> docs/index.html using the default theme
just serve    # serve docs/ on http://localhost:8000
just go       # build + serve in one step
```

Validate the resume data:

```shell
just validate
```

Watch for changes and rebuild automatically (requires `inotifywait`):

```shell
just watch
```

---

## Themes

The `themes/` directory contains alternate HTML layouts ported from
[FRESH Resume themes](https://github.com/fresh-themes). Pass a theme name
to `just build` to use one:

```shell
just build block           # default
just build positive
just build simple
just build simple-compact
just build actual
```

---

## GitHub Pages

`docs/index.html` is committed to the repository. GitHub Pages is configured
to serve from the `docs/` folder on the `master` branch - no CI build step
is required. Regenerate and commit `docs/index.html` whenever `resume.yaml`
changes.

---

## Planned renderers

Self-contained renderers in Go, Rust, Python, Elixir, and Java are in
development. Each will read `resume.yaml` directly (no `goresume` dependency)
and write `docs/index.html`. Implementation plans are in `plans/`.
