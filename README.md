# Stephen Brown II - Résumé Generator

This is the source code repository for my résumé. My résumé is maintained in a YAML data file which,
after being converted to JSON, is used as the data source for a variety of output formats and themes.

## Info

* [HackMyResume](https://github.com/hacksalot/HackMyResume)
* [FRESH resume schema](https://github.com/fresh-standard/fresh-resume-schema)
* [JSON resume](http://jsonresume.org/getting-started/)

## Setup

Install [wkhtmltopdf](http://wkhtmltopdf.org/downloads.html),
[hackmyresume](https://github.com/hacksalot/HackMyResume), and
[yq](https://yq.readthedocs.io/en/latest/) which are required for generation.

```shell
yay -S wkhtmltopdf nodejs-hackmyresume go-yq
# or
brew install wkhtmltopdf yq
npm install -g resume-cli hackmyresume
```

## Résumé Generation

Run the following command to output the résumé in all formats. The `out` directory is ignored by Git.

```shell
yq -Mj eval resume-fresh.yaml > resume-fresh.json
hackmyresume BUILD resume-fresh.json TO out/stephen-brown-ii.all -t modern
hackmyresume BUILD resume-fresh.json TO out/stephen-brown-ii.all -t ../some-folder/my-custom-theme/
hackmyresume BUILD resume-fresh.json TO out/stephen-brown-ii.all -t node_modules/jsonresume-theme-modern
```

Pre-defined FRESH themes are: `positive`, `modern`, `compact`, `basis` or `awesome` (Only supports LATEX, JSON, and YML formats)

## Convert FRESH résumé to JRS format

The FRESH format is the master file. The generated `resume.json` is ignored by Git.

```shell
hackmyresume CONVERT resume-fresh.json resume.json
```

## Résumé Analysis and Validation

Analyze and report on the résumé data:

```shell
hackmyresume ANALYZE resume-fresh.json
hackmyresume VALIDATE resume-fresh.json
```
