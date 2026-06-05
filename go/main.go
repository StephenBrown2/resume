package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	goyaml "github.com/goccy/go-yaml"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

func main() {
	input := flag.String("input", "../resume.yaml", "path to resume YAML")
	output := flag.String("output", "../docs/index.html", "path to write HTML")
	pdfOutput := flag.String("pdf", "", "path to write PDF (requires chromium; uses print-layout CSS)")
	cardOutput := flag.String("business-card", "", "path to write business card PDF (requires chromium)")
	sheetOutput := flag.String("sheet", "", "path to write 10-up Letter card sheet PDF for hand-cutting (requires chromium)")
	noGrid := flag.Bool("no-grid", false, "omit cut guides on --sheet output (use with perforated stock such as Avery 5371)")
	nameFont := flag.String("name-font", "Instrument Serif", "Google Fonts family for name heading")
	schema := flag.String("schema", "", "path to JSON Schema file (default: schema.json next to --input)")
	since := flag.String("since", "", "exclude jobs whose end date is before this date (YYYY, YYYY-MM, or YYYY-MM-DD)")
	skipVal := flag.Bool("skip-validation", false, "skip JSON Schema validation")
	flag.Parse()

	data, err := os.ReadFile(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read input: %v\n", err)
		os.Exit(1)
	}

	if !*skipVal {
		var raw map[string]any
		if err := goyaml.Unmarshal(data, &raw); err != nil {
			fmt.Fprintf(os.Stderr, "parse yaml for validation: %v\n", err)
			os.Exit(1)
		}
		if *schema == "" {
			*schema = resolveSchemaPath(*input, raw)
		}
		if err := validateSchema(*schema, raw); err != nil {
			fmt.Fprintln(os.Stderr, "validation error:", err)
			os.Exit(1)
		}
	}

	var resume Resume
	if err := goyaml.Unmarshal(data, &resume); err != nil {
		fmt.Fprintf(os.Stderr, "parse yaml: %v\n", err)
		os.Exit(1)
	}

	fontURLName := strings.ReplaceAll(*nameFont, " ", "+")
	googleFontsLink := template.HTML(fmt.Sprintf(
		`<link href="https://fonts.googleapis.com/css2?family=%s:ital@0;1&display=swap" rel="stylesheet">`,
		fontURLName,
	))
	nameFontCSS := template.CSS(fmt.Sprintf("'%s', Georgia, serif", *nameFont))

	sort.Slice(resume.Certificates, func(i, j int) bool {
		return resume.Certificates[i].Date > resume.Certificates[j].Date
	})
	shuffleKeywords(&resume)
	sortSkillSets(resume.Skills.Sets, resume.Skills.List)
	sortTestimonials(resume.Testimonials)

	work := resume.Work
	if *since != "" {
		work = filterWorkSince(work, *since)
	}
	groups := groupWork(work)

	tmplData := TemplateData{
		Basics:          resume.Basics,
		EmployerGroups:  groups,
		Projects:        resume.Projects,
		SkillSets:       resume.Skills.Sets,
		SkillList:       resume.Skills.List,
		Certificates:    resume.Certificates,
		Education:       resume.Education,
		Languages:       resume.Languages,
		Interests:       resume.Interests,
		Testimonials:    resume.Testimonials,
		GoogleFontsLink: googleFontsLink,
		NameFontCSS:     nameFontCSS,
	}

	funcMap := template.FuncMap{
		"formatDate":  formatDate,
		"nbspSummary": nbspShortWords,
		"levelClass":  levelClass,
		"skillByName": skillByName,
		"stripScheme": stripScheme,
		"certTitle":   certTitle,
		"certPrintID": certPrintID,
	}

	tmpl, err := template.New("resume").Funcs(funcMap).Parse(resumeTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse template: %v\n", err)
		os.Exit(1)
	}

	out, err := os.Create(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create output: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := out.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "close: %v\n", err)
			os.Exit(1)
		}
	}()

	if err := tmpl.Execute(out, tmplData); err != nil {
		fmt.Fprintf(os.Stderr, "render template: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "wrote %s\n", *output)

	if *pdfOutput != "" {
		if *output == "/dev/null" || *output == "-" {
			fmt.Fprintln(os.Stderr, "error: --pdf requires a real --output path (not /dev/null or -)")
			os.Exit(1)
		}
		if err := exportPDF(*output, *pdfOutput); err != nil {
			fmt.Fprintf(os.Stderr, "export pdf: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", *pdfOutput)
	}

	if *cardOutput != "" {
		if err := generateBusinessCard(resume.Basics, *cardOutput, nameFontCSS, googleFontsLink); err != nil {
			fmt.Fprintf(os.Stderr, "export business card: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", *cardOutput)
	}

	if *sheetOutput != "" {
		if err := generateBusinessCardSheet(resume.Basics, *sheetOutput, nameFontCSS, googleFontsLink, *noGrid); err != nil {
			fmt.Fprintf(os.Stderr, "export card sheet: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", *sheetOutput)
	}
}

// exportPDF renders htmlPath to a PDF at pdfPath using Chromium headless.
// The print-layout CSS (@media print) is applied automatically.
func exportPDF(htmlPath, pdfPath string) error {
	chrome, err := findChrome()
	if err != nil {
		return err
	}

	absHTML, err := filepath.Abs(htmlPath)
	if err != nil {
		return fmt.Errorf("resolve html path: %w", err)
	}
	absPDF, err := filepath.Abs(pdfPath)
	if err != nil {
		return fmt.Errorf("resolve pdf path: %w", err)
	}

	cmd := exec.Command(chrome,
		"--headless=new",
		"--disable-gpu",
		"--no-sandbox",
		"--print-to-pdf="+absPDF,
		"--print-to-pdf-no-header",
		"file://"+absHTML,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w\n%s", err, out)
	}
	return nil
}

// findChrome returns the path to the first Chromium/Chrome binary found in PATH.
func findChrome() (string, error) {
	for _, name := range []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable"} {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("no chromium/chrome binary found in PATH; install chromium to use --pdf")
}

// resolveSchemaPath derives the schema path from the YAML $schema field if
// present, falling back to schema.json in the same directory as the input file.
// The $schema value is resolved relative to the input file's directory.
func resolveSchemaPath(inputPath string, raw map[string]any) string {
	dir := filepath.Dir(inputPath)
	if s, ok := raw["$schema"].(string); ok && s != "" {
		if filepath.IsAbs(s) {
			return s
		}
		return filepath.Join(dir, s)
	}
	return filepath.Join(dir, "schema.json")
}

func validateSchema(schemaPath string, data map[string]any) error {
	c := jsonschema.NewCompiler()
	sch, err := c.Compile(schemaPath)
	if err != nil {
		return fmt.Errorf("compile schema: %w", err)
	}
	if err := sch.Validate(data); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
