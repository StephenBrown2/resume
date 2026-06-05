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

[working-directory: 'go']
go-pdf: go-render
    ./resume-renderer --input ../{{filename}} --output ../docs/index.html --pdf ../output/resume.pdf

[working-directory: 'go']
go-card: go-build
    ./resume-renderer --input ../{{filename}} --output /dev/null --business-card ../output/business-card.pdf

[working-directory: 'go']
go-sheet: go-build
    ./resume-renderer --input ../{{filename}} --output /dev/null --sheet ../output/card-sheet.pdf

[working-directory: 'go']
go-setup:
    curl -sSfL https://golangci-lint.run/install.sh | sh -s

[working-directory: 'go']
go-fmt:
    ./bin/golangci-lint fmt

[working-directory: 'go']
go-lint:
    ./bin/golangci-lint run

# ── Generic (calls all implemented language renderers) ───────────────────────

build: go-render

validate: go-validate

setup: go-setup

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
