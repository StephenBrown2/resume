# Stephen Brown II - Résumé Generator

This is the source code repository for my résumé. My résumé is maintained in a YAML data file which,
after being converted to JSON, is used as the data source for a variety of output formats and themes.

## Info

* [JSON resume](http://jsonresume.org/getting-started/)

## Setup

Install [goresume](https://github.com/nikaro/goresume) and
[yq](https://yq.readthedocs.io/en/latest/) which are required for generation.

A Justfile is also included, if you have the [just](https://just.systems/) command runner installed.

```shell
# with Go
go install github.com/nikaro/goresume@latest
go install github.com/mikefarah/yq/v4@latest

# with Homebrew
brew install nikaro/tap/goresume yq just

# on ArchLinux
yay -S goresume-bin go-yq just
```

## Résumé Generation

Run the following command to export the resume.yaml to `docs/index.html` and serve it on port 8000.

```shell
just go
```

A custom theme (ported from FRESH Resume themes) is in `themes/positive.html`.

### Github Pages

The generated page is already in the correct folder for Github Pages to pick it up and serve it at
<https://stephenbrown2.github.io/resume/>.


## Résumé Analysis and Validation

Validate the résumé data:

```shell
just validate
```
