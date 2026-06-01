export filename := "resume.yaml"

# ── Go ───────────────────────────────────────────────────────────────────────

[working-directory: 'go']
go-build:
    go build -o resume-renderer .

[working-directory: 'go']
go-render: go-build
    ./resume-renderer --input ../{{filename}} --output ../docs/index.html

[working-directory: 'go']
go-validate: go-build
    ./resume-renderer --input ../{{filename}} --output /dev/null

# ── Generic (calls all implemented language renderers) ───────────────────────

build: go-render

validate: go-validate

# ── Dev helpers ──────────────────────────────────────────────────────────────

serve directory="docs":
    python3 -m http.server --directory {{directory}}

dev: build serve

watch:
  #!/usr/bin/env sh
  inotifywait -m -r . \
    --exclude "(.*\\.pdf$)|public|justfile|\\.git" \
    -e close_write,move,create,delete \
    | while read -r directory events filename; do
      just build
    done
