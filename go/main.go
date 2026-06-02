package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	goyaml "github.com/goccy/go-yaml"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

func main() {
	input := flag.String("input", "../resume.yaml", "path to resume YAML")
	output := flag.String("output", "../docs/index.html", "path to write HTML")
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
	defer out.Close()

	if err := tmpl.Execute(out, tmplData); err != nil {
		fmt.Fprintf(os.Stderr, "render template: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "wrote %s\n", *output)
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
