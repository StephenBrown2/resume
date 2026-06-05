package main

// cssVars is the shared :root CSS variable block. Both the resume template and
// the business card template embed this via string concatenation so that the
// palette and font stacks stay in one place.
const cssVars = `  :root {
    --black:  #0a0a0a;
    --ink:    #1c1c1c;
    --muted:  #6a6a6a;
    --rule:   #e0e0e0;
    --page:   #fafaf8;
    --accent: #c0561a;
    --serif:  'Instrument Serif', Georgia, serif;
    --sans:   'Inter', 'Inter Variable', system-ui, sans-serif;
    --name-font: {{.NameFontCSS}};
  }`

// fontsHead is the shared <link> block for Google Fonts and Inter. Embed it
// inside <head> in each template. Requires .GoogleFontsLink in template data.
const fontsHead = `<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
{{.GoogleFontsLink}}
<link rel="stylesheet" href="https://rsms.me/inter/inter.css">`

// cssBodySmoothing is the shared body antialiasing and OpenType feature settings.
// Indented to 4 spaces to sit inside a body { } block.
const cssBodySmoothing = `    -webkit-font-smoothing: antialiased;
    font-feature-settings: 'dlig' 1, 'calt' 1, 'ss01' 1, 'ss04' 1, 'ss07' 1;`
