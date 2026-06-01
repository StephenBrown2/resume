# Go resume renderer

Reads `resume.yaml`, validates it against `schema.json`, and writes a styled single-file HTML resume to `docs/index.html`.

## Prerequisites

Go 1.26+. Install via [mise](https://mise.jdx.dev/), [asdf](https://asdf-vm.com/), or [go.dev/dl](https://go.dev/dl/).

## Build

```sh
go build -o resume-renderer .
# or from repo root:
just go-build
```

## Run

```sh
./resume-renderer [flags]
./resume-renderer --help
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `--input` | `../resume.yaml` | Path to input YAML file |
| `--output` | `../docs/index.html` | Path to write HTML output |
| `--name-font` | `Instrument Serif` | Google Fonts family name for the name heading |
| `--schema` | `../schema.json` | Path to JSON Schema file |
| `--skip-validation` | false | Skip JSON Schema validation of the YAML |

## Output

Writes a self-contained HTML file. Default path `../docs/index.html` is relative to the `go/` directory, so it lands at `docs/index.html` in the repo root.

## Module

Module path: `github.com/StephenBrown2/resume/go`. No CGo; pure Go only.
