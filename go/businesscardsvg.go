package main

import (
	"encoding/base64"
	"fmt"
	htmlpkg "html"
	htmltmpl "html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	txttmpl "text/template"

	qrcode "github.com/skip2/go-qrcode"
)

const goldSentinel = "#c8a832"

// extractFontName returns the first CSS font-family name from a value like
// "'Instrument Serif', Georgia, serif".
func extractFontName(css string) string {
	s := strings.TrimSpace(css)
	if s == "" {
		return "Instrument Serif"
	}
	if s[0] == '\'' || s[0] == '"' {
		if idx := strings.IndexByte(s[1:], s[0]); idx >= 0 {
			return s[1 : idx+1]
		}
	}
	if idx := strings.IndexByte(s, ','); idx >= 0 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}

type svgContactLine struct {
	Y    float64
	Text string
	URL  string // non-empty → wrap in <a>
}

type svgCardData struct {
	NameFont     string
	Name         string
	Label        string // uppercased
	LabelFill    string
	QRBase64     string
	ContactLines []svgContactLine
}

const (
	svgContentX      = 32.0  // 18pt bleed + 14pt padding
	svgLeftColW      = 149.0 // 224pt usable - 68pt QR - 7pt gap
	svgQRX           = 188.0 // svgContentX + svgLeftColW + 7
	svgQRY           = 80.0  // content bottom − QR size
	svgQRSize        = 68.0
	svgContentBotY   = 148.0 // 18 bleed + 144 live − 14 padding
	svgContactFS     = 6.5
	svgContactLineH  = 11.0
)

func buildSVGCardData(basics Basics, nameFont string) (svgCardData, error) {
	png, err := qrcode.Encode(basics.URL, qrcode.Medium, 300)
	if err != nil {
		return svgCardData{}, fmt.Errorf("generate qr code: %w", err)
	}

	type rawLine struct{ text, url string }
	var raw []rawLine
	raw = append(raw, rawLine{basics.Email, "mailto:" + basics.Email})

	phoneLine := basics.Phone
	if basics.Location.City != "" {
		sep := ""
		if phoneLine != "" {
			sep = " · "
		}
		phoneLine += sep + basics.Location.City + ", " + basics.Location.Region
	}
	if phoneLine != "" {
		raw = append(raw, rawLine{phoneLine, ""})
	}
	raw = append(raw, rawLine{stripScheme(basics.URL), basics.URL})
	for _, p := range basics.Profiles {
		raw = append(raw, rawLine{p.Network + ": " + stripScheme(p.URL), p.URL})
	}

	startY := svgContentBotY - float64(len(raw))*svgContactLineH + svgContactFS
	lines := make([]svgContactLine, len(raw))
	for i, r := range raw {
		lines[i] = svgContactLine{Y: startY + float64(i)*svgContactLineH, Text: r.text, URL: r.url}
	}

	return svgCardData{
		NameFont:     nameFont,
		Name:         basics.Name,
		Label:        strings.ToUpper(basics.Label),
		LabelFill:    goldSentinel,
		QRBase64:     base64.StdEncoding.EncodeToString(png),
		ContactLines: lines,
	}, nil
}

func xmlesc(s string) string { return htmlpkg.EscapeString(s) }

// svgCardTemplate produces a 4.0×2.5in SVG with bleed (viewBox 0 0 288 180).
// Gold elements use fill="#c8a832" as a sentinel; the export pipeline converts
// this to a proper spot color channel (Scribus) or solidColor ref (Inkscape).
const svgCardTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg"
  xmlns:xlink="http://www.w3.org/1999/xlink"
  width="4.0in" height="2.5in" viewBox="0 0 288 180">
  <defs><!-- spot colors injected by export pipeline --></defs>
  <rect width="288" height="180" fill="#fafaf8"/>
  <text x="32" y="49"
    font-family="{{xmlesc .NameFont}}, Liberation Serif, Georgia, serif"
    font-size="20" fill="#0a0a0a"
    textLength="149" lengthAdjust="spacingAndGlyphs">{{xmlesc .Name}}</text>
  <text x="32" y="63"
    font-family="Inter, Liberation Sans, Arial, sans-serif"
    font-size="8" font-weight="600"
    fill="{{.LabelFill}}"
    textLength="149" lengthAdjust="spacingAndGlyphs">{{xmlesc .Label}}</text>
{{- range .ContactLines}}
  {{- if .URL}}<a xlink:href="{{xmlesc .URL}}">{{end}}
  <text x="32" y="{{printf "%.2f" .Y}}"
    font-family="Inter, Liberation Sans, Arial, sans-serif"
    font-size="6.5" fill="#6a6a6a">{{xmlesc .Text}}</text>
  {{- if .URL}}</a>{{end}}
{{- end}}
  <image x="188" y="80" width="68" height="68"
    href="data:image/png;base64,{{.QRBase64}}"
    xlink:href="data:image/png;base64,{{.QRBase64}}"
    preserveAspectRatio="xMidYMid meet"/>
</svg>`

func renderSVGCard(data svgCardData, svgPath string) error {
	tmpl, err := txttmpl.New("card-svg").Funcs(txttmpl.FuncMap{"xmlesc": xmlesc}).Parse(svgCardTemplate)
	if err != nil {
		return fmt.Errorf("parse svg template: %w", err)
	}
	f, err := os.Create(svgPath)
	if err != nil {
		return fmt.Errorf("create svg: %w", err)
	}
	if err := tmpl.Execute(f, data); err != nil {
		f.Close()
		return fmt.Errorf("render svg: %w", err)
	}
	return f.Close()
}

// generateBusinessCard renders a 4.0×2.5in SVG (with bleed) and exports it to
// a PDF with a Gold spot color channel via Scribus or Inkscape.
// The SVG is kept alongside the PDF (same path, .svg extension).
func generateBusinessCard(basics Basics, pdfPath string, nameFontCSS htmltmpl.CSS, _ htmltmpl.HTML) error {
	data, err := buildSVGCardData(basics, extractFontName(string(nameFontCSS)))
	if err != nil {
		return err
	}
	svgPath := strings.TrimSuffix(pdfPath, ".pdf") + ".svg"
	if err := renderSVGCard(data, svgPath); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", svgPath)
	return exportSVGtoPDF(svgPath, pdfPath)
}

// exportSVGtoPDF detects Scribus (preferred) or Inkscape and exports the SVG
// to a PDF with a proper Gold spot color separation.
func exportSVGtoPDF(svgPath, pdfPath string) error {
	absSVG, err := filepath.Abs(svgPath)
	if err != nil {
		return fmt.Errorf("resolve svg path: %w", err)
	}
	absPDF, err := filepath.Abs(pdfPath)
	if err != nil {
		return fmt.Errorf("resolve pdf path: %w", err)
	}
	if scribus, err := exec.LookPath("scribus"); err == nil {
		return exportViaScribus(scribus, absSVG, absPDF)
	}
	if inkscape, err := exec.LookPath("inkscape"); err == nil {
		return exportViaInkscape(inkscape, absSVG, absPDF)
	}
	return fmt.Errorf("neither scribus nor inkscape found in PATH")
}

// scribyPyFmt is a format string for the Scribus Python script. SVG and PDF
// paths are baked in (%q) so no args need to be passed on the command line,
// avoiding Scribus treating them as documents to open.
const scribyPyFmt = `import scribus
scribus.openDoc(%q)
scribus.defineColorCMYKFloat("Gold", 0.0, 20.0, 80.0, 10.0)
scribus.setSpotColor("Gold", True)
scribus.replaceColor("FromSVG#c8a832", "Gold")
pdf = scribus.PDFfile()
pdf.file = %q
pdf.version = 15
pdf.outdst = 1
pdf.fontEmbedding = 1
pdf.save()
scribus.closeDoc()
`

func exportViaScribus(scribus, svgPath, pdfPath string) error {
	inputSVG := svgPath
	// Pre-outline text with Inkscape so textLength is baked into path geometry.
	if inkscape, err := exec.LookPath("inkscape"); err == nil {
		outlined, err := outlineTextWithInkscape(inkscape, svgPath, pdfPath)
		if err == nil {
			defer os.Remove(outlined)
			inputSVG = outlined
		}
	}

	script, err := os.CreateTemp("", "resume-scribus-*.py")
	if err != nil {
		return fmt.Errorf("create scribus script: %w", err)
	}
	defer os.Remove(script.Name())
	if _, err := fmt.Fprintf(script, scribyPyFmt, inputSVG, pdfPath); err != nil {
		script.Close()
		return err
	}
	if err := script.Close(); err != nil {
		return err
	}
	cmd := exec.Command(scribus, "--no-gui", "--python-script", script.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("scribus export: %w\n%s", err, out)
	}
	return nil
}

// outlineTextWithInkscape converts all text elements to paths so that SVG
// textLength constraints are baked into the path geometry before Scribus
// imports the file (Scribus ignores textLength on import).
func outlineTextWithInkscape(inkscape, svgPath, nearPath string) (string, error) {
	tmp, err := os.CreateTemp(filepath.Dir(nearPath), "resume-outlined-*.svg")
	if err != nil {
		return "", fmt.Errorf("create outlined svg: %w", err)
	}
	outPath := tmp.Name()
	tmp.Close()

	actions := "select-all;object-to-path;export-filename:" + outPath + ";export-type:svg;export-do"
	cmd := exec.Command(inkscape, "--actions", actions, svgPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(outPath)
		return "", fmt.Errorf("inkscape outline: %w\n%s", err, out)
	}
	return outPath, nil
}

func exportViaInkscape(inkscape, svgPath, pdfPath string) error {
	svgBytes, err := os.ReadFile(svgPath)
	if err != nil {
		return fmt.Errorf("read svg: %w", err)
	}
	processed := preprocessForInkscape(string(svgBytes))

	tmp, err := os.CreateTemp(filepath.Dir(pdfPath), "resume-inkscape-*.svg")
	if err != nil {
		return fmt.Errorf("create temp svg: %w", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(processed); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	actions := "select-all;object-to-path;export-filename:" + pdfPath + ";export-pdf-version:1.5;export-do"
	cmd := exec.Command(inkscape, "--actions", actions, tmp.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("inkscape export: %w\n%s", err, out)
	}
	return nil
}

// preprocessForInkscape adds the Inkscape-specific solidColor spot color
// definition and rewrites sentinel fills to reference it.
func preprocessForInkscape(svg string) string {
	svg = strings.Replace(svg,
		`<svg xmlns="http://www.w3.org/2000/svg"`,
		`<svg xmlns="http://www.w3.org/2000/svg"`+
			` xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"`+
			` xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.0.dtd"`,
		1)
	solidColor := `<solidColor id="gold-solid" style="solid-color:` + goldSentinel + `;solid-opacity:1" inkscape:label="Gold"/>`
	svg = strings.Replace(svg,
		`<defs><!-- spot colors injected by export pipeline --></defs>`,
		`<defs>`+solidColor+`</defs>`,
		1)
	svg = strings.ReplaceAll(svg, `fill="`+goldSentinel+`"`, `fill="url(#gold-solid)"`)
	return svg
}
