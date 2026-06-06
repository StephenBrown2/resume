package main

import (
	"encoding/base64"
	"fmt"
	htmlpkg "html"
	htmltmpl "html/template"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	txttmpl "text/template"

	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const goldSentinel = "#ffffff"
const cardBackground = "#F8F8F8"

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
	NameFontSize float64
	NameY        float64
	Name         string
	Label        string // uppercased
	LabelFill    string
	LabelSize    float64
	LabelY       float64
	TextWidth    float64 // textLength value for SVG fallback
	QRBase64     string
	QRX          float64
	QRY          float64
	QRSize       float64
	ContactLines []svgContactLine
	DebugBox     bool // draw safe-area guide rect
}

const (
	svgContentX     = 32.0  // 18pt bleed + 14pt padding
	svgContentRight = 256.0 // 18pt bleed + 224pt usable - 14pt padding
	svgContentBotY  = 148.0 // 18 bleed + 144 live − 14 padding
	svgContactFS    = 6.5
	svgContactLineH = 11.0
	svgColGap       = 7.0 // gap between left column and QR
)

func buildSVGCardData(basics Basics, nameFont string, nameFontSize, labelFontSize, nameY, labelY, textWidth float64, debugBox bool) (svgCardData, error) {
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

	qrSize := float64(len(raw)) * svgContactLineH
	qrX := svgContentRight - qrSize
	qrY := svgContentBotY - qrSize

	q, err := qrcode.New(basics.URL, qrcode.Medium)
	if err != nil {
		return svgCardData{}, fmt.Errorf("generate qr code: %w", err)
	}
	q.BackgroundColor = color.Transparent
	q.DisableBorder = true
	png, err := q.PNG(300)
	if err != nil {
		return svgCardData{}, fmt.Errorf("encode qr png: %w", err)
	}

	startY := svgContentBotY - float64(len(raw))*svgContactLineH + svgContactFS
	lines := make([]svgContactLine, len(raw))
	for i, r := range raw {
		lines[i] = svgContactLine{Y: startY + float64(i)*svgContactLineH, Text: r.text, URL: r.url}
	}

	return svgCardData{
		NameFont:     nameFont,
		NameFontSize: nameFontSize,
		NameY:        nameY,
		Name:         basics.Name,
		Label:        strings.ToUpper(basics.Label),
		LabelFill:    goldSentinel,
		LabelSize:    labelFontSize,
		LabelY:       labelY,
		TextWidth:    textWidth,
		QRBase64:     base64.StdEncoding.EncodeToString(png),
		QRX:          qrX,
		QRY:          qrY,
		QRSize:       qrSize,
		ContactLines: lines,
		DebugBox:     debugBox,
	}, nil
}

// countContactLines returns how many lines the contact block will have for
// the given basics, matching the logic in buildSVGCardData.
func countContactLines(basics Basics) int {
	n := 1 // email always present
	if basics.Phone != "" || basics.Location.City != "" {
		n++
	}
	n++ // URL
	n += len(basics.Profiles)
	return n
}

// findFontFile returns the font file path for an exactly-matched fontconfig family.
// Returns an error if fontconfig resolves to a different family (i.e. the requested
// family is not installed).
func findFontFile(family, style string) (string, error) {
	out, err := exec.Command("fc-match", family+":style="+style, "--format=%{family}\t%{file}").Output()
	if err != nil {
		return "", fmt.Errorf("fc-match: %w", err)
	}
	parts := strings.SplitN(strings.TrimSpace(string(out)), "\t", 2)
	if len(parts) != 2 || parts[1] == "" {
		return "", fmt.Errorf("fc-match: unexpected output %q", out)
	}
	if !strings.EqualFold(strings.TrimSpace(parts[0]), family) {
		return "", fmt.Errorf("fc-match returned %q instead of %q", parts[0], family)
	}
	return parts[1], nil
}

// firstInstalledFont returns the path to the first family in families that
// fontconfig can resolve exactly at the given style.
func firstInstalledFont(families []string, style string) (string, error) {
	for _, f := range families {
		if path, err := findFontFile(f, style); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("none of %v installed (style=%q)", families, style)
}

// textAdvancePt returns the advance width in points for text at ptSize, using
// the font file at fontPath. Loaded at 72 DPI so 1 fixed-point pixel = 1 pt.
func textAdvancePt(text, fontPath string, ptSize float64) (float64, error) {
	raw, err := os.ReadFile(fontPath)
	if err != nil {
		return 0, fmt.Errorf("read font: %w", err)
	}
	f, err := opentype.Parse(raw)
	if err != nil {
		return 0, fmt.Errorf("parse font: %w", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{Size: ptSize, DPI: 72})
	if err != nil {
		return 0, fmt.Errorf("create face: %w", err)
	}
	defer face.Close()
	var total fixed.Int26_6
	for _, r := range text {
		if adv, ok := face.GlyphAdvance(r); ok {
			total += adv
		}
	}
	return float64(total) / 64.0, nil
}

// ptSizeToFill returns the font size (pt) that makes text's advance width equal
// targetWidth. Falls back to nominalSize if measurement fails.
func ptSizeToFill(text, fontPath string, nominalSize, targetWidth float64) float64 {
	w, err := textAdvancePt(text, fontPath, nominalSize)
	if err != nil || w <= 0 {
		fmt.Fprintf(os.Stderr, "warn: font measurement failed (%v); using nominal size\n", err)
		return nominalSize
	}
	return nominalSize * targetWidth / w
}


func xmlesc(s string) string { return htmlpkg.EscapeString(s) }

// svgCardTemplate produces a 4.0×2.5in SVG with bleed (viewBox 0 0 288 180).
// Gold elements use fill="#ffffff" (white) as a sentinel; the export pipeline
// converts this to a proper spot color channel (Scribus) or solidColor ref
// (Inkscape). White is required: gold toner applies over the base layer, so a
// non-white base tints the gold.
const svgCardTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg"
  xmlns:xlink="http://www.w3.org/1999/xlink"
  width="4.0in" height="2.5in" viewBox="0 0 288 180">
  <defs><!-- spot colors injected by export pipeline --></defs>
  <rect width="288" height="180" fill="` + cardBackground + `"/>
  <text x="32" y="{{printf "%.2f" .NameY}}"
    font-family="{{xmlesc .NameFont}}, Liberation Serif, Georgia, serif"
    font-size="{{printf "%.2f" .NameFontSize}}" fill="#0a0a0a"
    textLength="{{printf "%.2f" .TextWidth}}" lengthAdjust="spacingAndGlyphs">{{xmlesc .Name}}</text>
  <text x="32" y="{{printf "%.2f" .LabelY}}"
    font-family="Inter, Liberation Sans, Arial, sans-serif"
    font-size="{{printf "%.2f" .LabelSize}}" font-weight="600"
    fill="{{.LabelFill}}"
    textLength="{{printf "%.2f" .TextWidth}}" lengthAdjust="spacingAndGlyphs">{{xmlesc .Label}}</text>
{{- range .ContactLines}}
  {{- if .URL}}<a xlink:href="{{xmlesc .URL}}">{{end}}
  <text x="32" y="{{printf "%.2f" .Y}}"
    font-family="Inter, Liberation Sans, Arial, sans-serif"
    font-size="6.5" fill="#6a6a6a">{{xmlesc .Text}}</text>
  {{- if .URL}}</a>{{end}}
{{- end}}
  <image x="{{printf "%.2f" .QRX}}" y="{{printf "%.2f" .QRY}}" width="{{printf "%.2f" .QRSize}}" height="{{printf "%.2f" .QRSize}}"
    href="data:image/png;base64,{{.QRBase64}}"
    xlink:href="data:image/png;base64,{{.QRBase64}}"
    preserveAspectRatio="xMidYMid meet"/>
{{- if .DebugBox}}
  <rect x="32" y="32" width="224" height="116" fill="none" stroke="#ff0000" stroke-width="0.5"/>
{{- end}}
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

const (
	defaultNameFontSize  = 20.0
	defaultLabelFontSize = 9.25
	// font metric ratios (Instrument Serif name, Inter 600 label)
	serifCapHeightRatio    = 0.65
	serifDescenderRatio    = 0.30
	sansCapHeightRatio     = 0.73
	textTopPad             = 4.0 // pt above cap-height from content top
	interTextGap           = 3.0 // pt between name descender bottom and label cap top
)

// goMetricsAdjust corrects for GPOS kerning and optical overhang that
// golang.org/x/image/font/opentype does not account for. Scribus renders
// the same font/text ~5-7% wider than the simple advance-width sum, so we
// target a slightly narrower effective width. textLength="224" in the SVG
// then stretches the result back to the full safe-area width.
const goMetricsAdjust = 0.93

// computeCardLayout computes font sizes and y positions for the name and label
// using Go font metrics. The measured width target is scaled by goMetricsAdjust
// so Scribus's larger rendered width lands within the 224pt safe area.
func computeCardLayout(name, label, nameFont string, targetWidth float64) (namePt, labelPt, nameY, labelY float64) {
	namePt = defaultNameFontSize
	labelPt = defaultLabelFontSize

	effective := targetWidth * goMetricsAdjust
	nameFamilies := []string{nameFont, "Liberation Serif", "Georgia", "DejaVu Serif"}
	labelFamilies := []string{"Inter", "Liberation Sans", "Arial", "DejaVu Sans"}

	if namePath, err := firstInstalledFont(nameFamilies, "Regular"); err == nil {
		namePt = ptSizeToFill(name, namePath, defaultNameFontSize, effective)
	} else {
		fmt.Fprintf(os.Stderr, "warn: name font not found (%v); using %.2fpt\n", err, namePt)
	}
	if labelPath, err := firstInstalledFont(labelFamilies, "SemiBold"); err == nil {
		labelPt = ptSizeToFill(strings.ToUpper(label), labelPath, defaultLabelFontSize, effective)
	} else if labelPath, err := firstInstalledFont(labelFamilies, "Bold"); err == nil {
		labelPt = ptSizeToFill(strings.ToUpper(label), labelPath, defaultLabelFontSize, effective)
	} else {
		fmt.Fprintf(os.Stderr, "warn: label font not found (%v); using %.2fpt\n", err, labelPt)
	}

	nameY = svgContentX + textTopPad + namePt*serifCapHeightRatio
	labelY = nameY + namePt*serifDescenderRatio + interTextGap + labelPt*sansCapHeightRatio
	return
}

// generateBusinessCard renders a 4.0×2.5in SVG (with bleed) and exports it to
// a PDF with a spot color channel via Scribus.
// spotColorName selects the channel name (e.g. "Gold", "RDG_Gold", "PANTONE 871 C").
// debugCard adds a visible safe-area guide rect to the SVG for layout verification.
// The SVG is kept alongside the PDF (same path, .svg extension).
func generateBusinessCard(basics Basics, pdfPath string, nameFontCSS htmltmpl.CSS, _ htmltmpl.HTML, spotColorName string, debugCard bool) error {
	nameFont := extractFontName(string(nameFontCSS))
	targetWidth := svgContentRight - svgContentX // full safe-area width = 224pt

	nameFontSize, labelFontSize, nameY, labelY := computeCardLayout(
		basics.Name, strings.ToUpper(basics.Label), nameFont, targetWidth)
	fmt.Fprintf(os.Stderr, "font sizes: name=%.2fpt label=%.2fpt  y: name=%.2f label=%.2f\n",
		nameFontSize, labelFontSize, nameY, labelY)

	data, err := buildSVGCardData(basics, nameFont, nameFontSize, labelFontSize, nameY, labelY, targetWidth, debugCard)
	if err != nil {
		return err
	}
	svgPath := strings.TrimSuffix(pdfPath, ".pdf") + ".svg"
	if err := renderSVGCard(data, svgPath); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", svgPath)
	return exportSVGtoPDF(svgPath, pdfPath, spotColorName)
}

// exportSVGtoPDF exports the SVG to a PDF with a spot color channel via Scribus.
func exportSVGtoPDF(svgPath, pdfPath, spotColorName string) error {
	absSVG, err := filepath.Abs(svgPath)
	if err != nil {
		return fmt.Errorf("resolve svg path: %w", err)
	}
	absPDF, err := filepath.Abs(pdfPath)
	if err != nil {
		return fmt.Errorf("resolve pdf path: %w", err)
	}
	scribus, err := exec.LookPath("scribus")
	if err != nil {
		return fmt.Errorf("scribus not found in PATH; install scribus to use --business-card")
	}
	return exportViaScribus(scribus, absSVG, absPDF, spotColorName)
}

// scribyPyFmt is a format string for the Scribus Python script. SVG path,
// PDF path, and spot color name (×3) are baked in (%q / %s) so no args need
// to be passed on the command line, avoiding Scribus treating them as
// documents to open. Args: svgPath, spotColorName (×3), pdfPath.
const scribyPyFmt = `import scribus
scribus.openDoc(%q)
scribus.defineColorCMYKFloat(%q, 0.0, 20.0, 80.0, 10.0)
scribus.setSpotColor(%q, True)
scribus.replaceColor("FromSVG#ffffff", %q)
pdf = scribus.PDFfile()
pdf.file = %q
pdf.version = 15
pdf.outdst = 1
pdf.fontEmbedding = 1
pdf.save()
scribus.closeDoc()
`

func exportViaScribus(scribus, svgPath, pdfPath, spotColorName string) error {
	script, err := os.CreateTemp("", "resume-scribus-*.py")
	if err != nil {
		return fmt.Errorf("create scribus script: %w", err)
	}
	defer os.Remove(script.Name())
	if _, err := fmt.Fprintf(script, scribyPyFmt, svgPath, spotColorName, spotColorName, spotColorName, pdfPath); err != nil {
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

