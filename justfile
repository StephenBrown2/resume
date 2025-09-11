serve:
    python3 -m http.server

go:
    goresume export --resume resume.yaml --html-theme positive --html-output docs/index.html
    python3 -m http.server --directory docs

validate:
    goresume validate --resume resume.yaml