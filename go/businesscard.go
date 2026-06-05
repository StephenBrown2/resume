package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	qrcode "github.com/skip2/go-qrcode"
)

// BusinessCardData is passed to the business card and sheet HTML templates.
type BusinessCardData struct {
	Basics          Basics
	QRCodeDataURI   template.URL
	GoogleFontsLink template.HTML
	NameFontCSS     template.CSS
}

// BusinessCardSheetData wraps BusinessCardData for the 10-up Letter sheet.
// CardSlots is length 10; the template ranges over it and passes $ (root) to
// the named "card" sub-template so the card fields remain accessible.
type BusinessCardSheetData struct {
	BusinessCardData
	CardSlots []struct{}
	NoGrid    bool
}

// buildCardData generates the QR code PNG and assembles BusinessCardData.
func buildCardData(basics Basics, nameFontCSS template.CSS, googleFontsLink template.HTML) (BusinessCardData, error) {
	png, err := qrcode.Encode(basics.URL, qrcode.Medium, 300)
	if err != nil {
		return BusinessCardData{}, fmt.Errorf("generate qr code: %w", err)
	}
	return BusinessCardData{
		Basics:          basics,
		QRCodeDataURI:   template.URL("data:image/png;base64," + base64.StdEncoding.EncodeToString(png)),
		GoogleFontsLink: googleFontsLink,
		NameFontCSS:     nameFontCSS,
	}, nil
}

// renderCardToFile parses tmplSrc, executes it with data, writes the HTML to a
// temp file in the same directory as pdfPath (so snap-sandboxed Chromium can
// reach it), and exports to pdfPath via Chromium headless.
func renderCardToFile(tmplName, tmplSrc string, data any, pdfPath string) error {
	funcMap := template.FuncMap{"stripScheme": stripScheme}
	tmpl, err := template.New(tmplName).Funcs(funcMap).Parse(tmplSrc)
	if err != nil {
		return fmt.Errorf("parse %s template: %w", tmplName, err)
	}
	absPDF, err := filepath.Abs(pdfPath)
	if err != nil {
		return fmt.Errorf("resolve pdf path: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(absPDF), "resume-card-*.html")
	if err != nil {
		return fmt.Errorf("create temp html: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmpl.Execute(tmp, data); err != nil {
		tmp.Close()
		return fmt.Errorf("render %s: %w", tmplName, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp html: %w", err)
	}
	return exportPDF(tmpPath, pdfPath)
}

// generateBusinessCardSheet renders a 10-up (2×5) Letter-size sheet PDF
// for printing on card stock and hand-cutting. Set noGrid to omit the dashed
// cut guides (use when printing on perforated stock such as Avery 5371).
func generateBusinessCardSheet(basics Basics, pdfPath string, nameFontCSS template.CSS, googleFontsLink template.HTML, noGrid bool) error {
	card, err := buildCardData(basics, nameFontCSS, googleFontsLink)
	if err != nil {
		return err
	}
	data := BusinessCardSheetData{
		BusinessCardData: card,
		CardSlots:        make([]struct{}, 10),
		NoGrid:           noGrid,
	}
	return renderCardToFile("business-card-sheet", businessCardSheetTemplate, data, pdfPath)
}

// cardCSS is shared card layout CSS, scoped under .card.
const cardCSS = `.card {
  width: 3.5in;
  height: 2in;
  padding: 0.19in;
  display: grid;
  grid-template-columns: 1fr 0.95in;
  gap: 0.1in;
  background: var(--page);
  color: var(--ink);
  overflow: hidden;
}
.card .left { display: flex; flex-direction: column; }
.card .name {
  font-family: var(--name-font);
  font-size: 20pt;
  line-height: 1.05;
  color: var(--black);
}
.card .label {
  font-size: 8pt;
  font-weight: 600;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--accent);
  margin-top: 5pt;
}
.card .contact {
  margin-top: auto;
  font-size: 6.5pt;
  color: var(--muted);
  line-height: 1.75;
}
.card .contact a { color: inherit; text-decoration: none; }
.card .right { display: flex; align-items: flex-end; justify-content: flex-end; }
.card .qr img { width: 0.95in; height: 0.95in; display: block; }`

// cardSubTemplate defines a named "card" sub-template used by both the single
// card and sheet templates. The sheet calls {{template "card" $}} inside a
// range so that $ (root data) is passed instead of the loop variable.
const cardSubTemplate = `{{define "card"}}<div class="card">
  <div class="left">
    <div>
      <div class="name">{{.Basics.Name}}</div>
      <div class="label">{{.Basics.Label}}</div>
    </div>
    <div class="contact">
      <a href="mailto:{{.Basics.Email}}">{{.Basics.Email}}</a><br>
      {{- if .Basics.Phone}}{{.Basics.Phone}}{{end -}}
      {{- if .Basics.Location.City}} &middot; {{.Basics.Location.City}}, {{.Basics.Location.Region}}{{end}}<br>
      <a href="{{.Basics.URL}}">{{stripScheme .Basics.URL}}</a>
      {{- range .Basics.Profiles}}<br>{{.Network}}: <a href="{{.URL}}">{{stripScheme .URL}}</a>{{end}}
    </div>
  </div>
  <div class="right">
    <div class="qr"><img src="{{.QRCodeDataURI}}" alt="QR code for {{.Basics.URL}}"></div>
  </div>
</div>{{end}}`

const businessCardTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
` + fontsHead + `
<style>
` + cssVars + `
@page { size: 3.5in 2in; margin: 0; }
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
html, body { width: 3.5in; height: 2in; }
body {
  font-family: var(--sans);
` + cssBodySmoothing + `
}
` + cardCSS + `
</style>
</head>
<body>
{{template "card" .}}
</body>
</html>
` + cardSubTemplate

const businessCardSheetTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
` + fontsHead + `
<style>
` + cssVars + `
@page { size: letter portrait; margin: 0; }
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
html, body { width: 8.5in; height: 11in; background: white; }
body {
  font-family: var(--sans);
` + cssBodySmoothing + `
}
.sheet {
  width: 8.5in;
  height: 11in;
  padding: 0.5in 0.75in;
  display: grid;
  grid-template-columns: 3.5in 3.5in;
  grid-template-rows: repeat(5, 2in);
  gap: 0;
}
.slot { width: 3.5in; height: 2in; }
{{if not .NoGrid}}.slot { border: 0.4pt dashed #bbb; }{{end}}
` + cardCSS + `
</style>
</head>
<body>
<div class="sheet">
  {{range .CardSlots}}<div class="slot">{{template "card" $}}</div>{{end}}
</div>
</body>
</html>
` + cardSubTemplate
