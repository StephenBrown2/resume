export filename := "resume.yaml"

build theme="block":
    goresume export --resume {{filename}} --html-theme {{theme}} --html-output docs/index.html

serve directory="docs":
    python3 -m http.server --directory {{directory}}

go: build serve

validate:
    goresume validate --resume {{filename}}

watch:
  #!/usr/bin/env sh
  inotifywait -m -r . \
    --exclude "(.*\\.pdf$)|public|justfile|\\.git" \
    -e close_write,move,create,delete \
    | while read -r directory events filename; do
      just go
    done
