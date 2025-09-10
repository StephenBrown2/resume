yaml:
    yq -M -o=json eval resume.yaml > resume.json

serve:
    python3 -m http.server

go: yaml
    goresume export --resume resume.json --html-theme positive --html-output docs/index.html
    python3 -m http.server --directory docs