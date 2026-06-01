package main

const resumeTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Basics.Name}} — Resume</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
{{.GoogleFontsLink}}
<link rel="stylesheet" href="https://rsms.me/inter/inter.css">
<style>
  :root {
    --black:  #0a0a0a;
    --ink:    #1c1c1c;
    --muted:  #6a6a6a;
    --rule:   #e0e0e0;
    --page:   #fafaf8;
    --accent: #c0561a;
    --serif:  'Instrument Serif', Georgia, serif;
    --sans:   'Inter', 'Inter Variable', system-ui, sans-serif;
    --name-font: {{.NameFontCSS}};
  }

  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  html { font-size: 15px; }

  body {
    font-family: var(--sans);
    background: var(--page);
    color: var(--ink);
    line-height: 1.5;
    -webkit-font-smoothing: antialiased;
    font-feature-settings: 'dlig' 1, 'calt' 1, 'ss01' 1, 'ss04' 1, 'ss07' 1;
  }

  .page {
    max-width: 805px;
    margin: 0 auto;
    padding: 48px 48px 64px;
  }

  header {
    display: grid;
    grid-template-columns: 1fr auto;
    align-items: end;
    gap: 24px;
    padding-bottom: 20px;
    border-bottom: 1.5px solid var(--black);
    margin-bottom: 26px;
  }

  .title-label {
    font-size: 0.68rem;
    font-weight: 600;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: var(--accent);
    margin-bottom: 5px;
  }

  .name {
    font-family: var(--name-font);
    font-size: 2.5rem;
    font-weight: 400;
    letter-spacing: -0.01em;
    line-height: 1;
    color: var(--black);
  }

  a { color: var(--muted); text-decoration: none; }
  a:hover { color: var(--accent); }

  .contact {
    text-align: right;
    font-size: 0.76rem;
    color: var(--muted);
    line-height: 1.85;
  }

  section { margin-bottom: 22px; }

  .section-label {
    font-size: 0.64rem;
    font-weight: 600;
    letter-spacing: 0.15em;
    text-transform: uppercase;
    color: var(--accent);
    margin-bottom: 11px;
  }

  .summary p {
    font-size: 0.87rem;
    font-weight: 300;
    line-height: 1.7;
    color: var(--ink);
  }

  .job { margin-bottom: 16px; }
  .job:last-child { margin-bottom: 0; }

  .job-header {
    display: grid;
    grid-template-columns: 1fr auto;
    align-items: baseline;
    margin-bottom: 1px;
  }

  .job-title { font-size: 0.87rem; font-weight: 600; color: var(--black); }
  .job-dates  { font-size: 0.74rem; color: var(--muted); font-variant-numeric: tabular-nums; }

  .job-meta {
    font-size: 0.76rem;
    color: var(--muted);
    margin-bottom: 7px;
  }
  .job-meta a { font-weight: 500; }

  .job-summary {
    font-size: 0.78rem;
    color: var(--muted);
    font-style: italic;
    margin-bottom: 5px;
  }

  .highlights { list-style: none; }
  .highlights li {
    font-size: 0.81rem;
    font-weight: 300;
    line-height: 1.55;
    padding-left: 13px;
    position: relative;
    margin-bottom: 3px;
    color: var(--ink);
  }
  .highlights li::before {
    content: '–';
    position: absolute;
    left: 0;
    color: var(--muted);
  }

  .tags { display: flex; flex-wrap: wrap; gap: 4px; margin-top: 7px; }
  .tag {
    font-size: 0.65rem;
    font-weight: 500;
    color: var(--muted);
    border: 1px solid var(--rule);
    border-radius: 3px;
    padding: 1px 6px;
  }

  .job-divider {
    border: none;
    border-top: 1px solid var(--rule);
    margin: 14px 0;
  }

  .employer-group { margin-bottom: 0; }

  .employer-header {
    display: grid;
    grid-template-columns: 1fr auto;
    align-items: baseline;
    gap: 8px;
    padding-bottom: 6px;
    margin-bottom: 8px;
    border-bottom: 1px solid var(--rule);
  }

  .employer-name {
    font-size: 0.87rem;
    font-weight: 700;
    color: var(--black);
  }
  .employer-name a, .project-name a { color: inherit; }

  .employer-former {
    font-size: 0.73rem;
    color: var(--muted);
    margin-left: 6px;
  }

  .employer-group .job {
    padding-left: 10px;
    border-left: 2px solid var(--rule);
    margin-bottom: 0;
  }

  .position-divider {
    border: none;
    border-top: 1px dashed var(--rule);
    margin: 10px 0 10px 10px;
  }

  .job-dates,
  .skill-item,
  .edu-detail {
    font-variant-numeric: tabular-nums;
  }

  .projects-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
  }

  .project {
    border: 1px solid var(--rule);
    border-radius: 4px;
    padding: 10px 12px;
  }

  .project-name { font-size: 0.81rem; font-weight: 600; color: var(--black); margin-bottom: 3px; }

  .project-desc { font-size: 0.75rem; color: var(--muted); font-weight: 300; line-height: 1.5; }

  .project-tags { display: flex; flex-wrap: wrap; gap: 3px; margin-top: 5px; }

  .skills-domains {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 18px 40px;
  }

  .skill-group-label {
    font-size: 0.64rem;
    font-weight: 600;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--muted);
    padding-bottom: 5px;
    border-bottom: 1px solid var(--rule);
    margin-bottom: 6px;
  }

  .skill-item {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    font-size: 0.79rem;
    padding-bottom: 4px;
    border-bottom: 1px solid var(--rule);
  }

  .skill-name { font-weight: 500; color: var(--black); }
  .skill-level {
    font-size: 0.64rem;
    font-weight: 600;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    color: var(--muted);
  }
  .skill-level.adv { color: var(--black); }
  .skill-level.mid { color: var(--ink); }

  .footer-grid {
    display: grid;
    grid-template-columns: 1fr 1.7fr;
    gap: 0 40px;
    align-items: start;
  }

  .edu-degree { font-size: 0.81rem; font-weight: 600; color: var(--black); }
  .edu-detail { font-size: 0.73rem; color: var(--muted); margin-top: 1px; }
  .cert-list { font-size: 0.76rem; color: var(--muted); line-height: 1.8; }

  .testimonial {
    border-left: 2px solid var(--accent);
    padding: 10px 16px;
    margin-top: 4px;
    page-break-inside: avoid;
    break-inside: avoid;
  }

  .testimonial blockquote {
    font-size: 0.82rem;
    font-weight: 300;
    font-style: italic;
    line-height: 1.65;
    color: var(--ink);
    margin-bottom: 6px;
  }

  .testimonial cite {
    font-size: 0.72rem;
    font-weight: 500;
    color: var(--muted);
    font-style: normal;
    letter-spacing: 0.02em;
  }

  .print-only { display: none; }

  @page {
    size: letter;
    margin: 0.3in;
  }

  @media print {
    html { font-size: 12pt; }
    body { background: white; }
    .page { max-width: 100%; padding: 0; }

    .print-only { display: inline; }
    a { color: inherit !important; text-decoration: none !important; }

    .job            { page-break-inside: avoid; }
    .project        { page-break-inside: avoid; }
    .section-intro  { break-inside: avoid; page-break-inside: avoid; }
    header          { page-break-after: avoid; }
  }

  @media (max-width: 600px) {
    .page { padding: 28px 20px 48px; }
    header { grid-template-columns: 1fr; }
    .contact { text-align: left; }
    .skills-domains { grid-template-columns: 1fr; }
    .projects-grid { grid-template-columns: 1fr; }
    .footer-grid { grid-template-columns: 1fr; gap: 14px 0; }
    .name { font-size: 1.9rem; }
  }
</style>
</head>
<body>
<div class="page">

  <header>
    <div>
      <p class="title-label">{{.Basics.Label}}</p>
      <h1 class="name">{{.Basics.Name}}</h1>
    </div>
    <div class="contact">
      <a href="mailto:{{.Basics.Email}}">{{.Basics.Email}}</a><br>
      {{.Basics.Phone}} &middot; {{.Basics.Location.City}}, {{.Basics.Location.Region}}<br>
      {{range .Basics.Profiles -}}
      <a href="{{.URL}}">{{stripScheme .URL}}</a><br>
      {{- end}}
      <span class="print-only"><a href="{{.Basics.URL}}">{{stripScheme .Basics.URL}}</a></span>
    </div>
  </header>

  <section class="summary">
    <div class="section-intro">
      <div class="section-label">Profile</div>
      <p>{{nbspSummary .Basics.Summary}}</p>
    </div>
  </section>

  <section>
    <div class="section-intro">
      <div class="section-label">Experience</div>
      {{- with index .EmployerGroups 0}}{{template "employer-group" .}}{{end}}
    </div>
    {{- range $i, $group := .EmployerGroups}}
    {{- if gt $i 0}}
    <hr class="job-divider">
    {{template "employer-group" $group}}
    {{- end}}
    {{- end}}
  </section>

  <section>
    <div class="section-intro">
      <div class="section-label">Open Source &amp; Projects</div>
      <div class="projects-grid">
        {{- range .Projects}}
        <div class="project">
          <div class="project-name">
            {{- if .URL}}<a href="{{.URL}}">{{.Name}}</a>{{else}}{{.Name}}{{end}}
          </div>
          <div class="project-desc">{{.Description}}</div>
          {{- if .Keywords}}
          <div class="project-tags">
            {{- range .Keywords}}<span class="tag">{{.}}</span>{{end}}
          </div>
          {{- end}}
        </div>
        {{- end}}
      </div>
    </div>
  </section>

  <section>
    <div class="section-intro">
      <div class="section-label">Skills</div>
      <div class="skills-domains">
        {{- range .SkillSets}}
        <div>
          <div class="skill-group-label">{{.Name}}</div>
          {{- $list := $.SkillList}}
          {{- range .Skills}}
          {{- $skill := skillByName $list .}}
          <div class="skill-item">
            <span class="skill-name">{{$skill.Name}}</span>
            {{- $cls := levelClass $skill.Level}}
            {{- if $cls}}
            <span class="skill-level {{$cls}}">{{$skill.Level}}</span>
            {{- else}}
            <span class="skill-level">{{$skill.Level}}</span>
            {{- end}}
          </div>
          {{- end}}
        </div>
        {{- end}}
      </div>
    </div>
  </section>

  <section>
    <div class="section-intro">
      <div class="section-label">Education &amp; Certifications</div>
      <div class="footer-grid">
      <div>
        {{- range .Education}}
        <div class="edu-degree">{{.StudyType}} {{.Area}}</div>
        <div class="edu-detail">
          {{- if .URL}}<a href="{{.URL}}">{{.Institution}}</a>{{else}}{{.Institution}}{{end}}
          {{- if .EndDate}} &middot; {{formatDate .EndDate}}{{end}}
        </div>
        {{- end}}
      </div>
      <div class="cert-list">
        {{- range $i, $cert := .Certificates}}
        {{- if gt $i 0}} &nbsp;&middot;&nbsp; {{end}}
        {{- if $cert.URL}}<a href="{{$cert.URL}}"{{with certTitle $cert}} title="{{.}}"{{end}}>{{$cert.Name}}</a>{{else}}<span{{with certTitle $cert}} title="{{.}}"{{end}}>{{$cert.Name}}</span>{{end}}
        {{- with certPrintID $cert}}<span class="print-only"> ({{.}})</span>{{end}}
        {{- end}}
      </div>
    </div>
    </div>
  </section>

  <section class="references">
    <div class="section-intro">
      <div class="section-label">References</div>
      {{- with index .Testimonials 0}}{{template "testimonial" .}}{{end}}
    </div>
    {{- range $i, $t := .Testimonials}}
    {{- if gt $i 0}}
    {{template "testimonial" $t}}
    {{- end}}
    {{- end}}
  </section>

</div>
</body>
</html>

{{define "employer-group"}}
{{- if gt (len .Positions) 1}}
    <div class="employer-group">
      <div class="employer-header">
        <div>
          {{- if .URL}}
          <span class="employer-name"><a href="{{.URL}}">{{.DisplayName}}</a></span>
          {{- else}}
          <span class="employer-name">{{.DisplayName}}</span>
          {{- end}}
          {{- range .FormerNames}}
          <span class="employer-former">(formerly {{.}})</span>
          {{- end}}
        </div>
        <span class="job-dates">{{formatDate .StartDate}} &#8211; {{formatDate .EndDate}}</span>
      </div>
      {{- range $j, $pos := .Positions}}
      {{- if gt $j 0}}
      <div class="position-divider"></div>
      {{- end}}
      <div class="job">
        <div class="job-header">
          <span class="job-title">{{$pos.Position}}</span>
          <span class="job-dates">{{formatDate $pos.StartDate}} &#8211; {{formatDate $pos.EndDate}}</span>
        </div>
        {{- if $pos.Summary}}
        <div class="job-summary">{{$pos.Summary}}</div>
        {{- end}}
        {{- if or $pos.Location (ne $pos.Employer $.DisplayName)}}
        <div class="job-meta">
          {{- if ne $pos.Employer $.DisplayName}}{{$pos.Employer}}{{if $pos.Location}} &middot; {{end}}{{end}}
          {{- $pos.Location}}</div>
        {{- end}}
        {{- if $pos.Highlights}}
        <ul class="highlights">
          {{- range $pos.Highlights}}
          <li>{{.}}</li>
          {{- end}}
        </ul>
        {{- end}}
        {{- if $pos.Keywords}}
        <div class="tags">
          {{- range $pos.Keywords}}<span class="tag">{{.}}</span>{{end}}
        </div>
        {{- end}}
      </div>
      {{- end}}
    </div>
{{- else}}
    {{- with index .Positions 0}}
    <div class="job">
      <div class="job-header">
        <span class="job-title">{{.Position}}</span>
        <span class="job-dates">{{formatDate .StartDate}} &#8211; {{formatDate .EndDate}}</span>
      </div>
      {{- if .Summary}}
      <div class="job-summary">{{.Summary}}</div>
      {{- end}}
      <div class="job-meta">
        {{- if .URL}}<a href="{{.URL}}">{{$.DisplayName}}</a>{{else}}{{$.DisplayName}}{{end}}
        {{- if .Location}} &middot; {{.Location}}{{end}}
      </div>
      {{- if .Highlights}}
      <ul class="highlights">
        {{- range .Highlights}}
        <li>{{.}}</li>
        {{- end}}
      </ul>
      {{- end}}
      {{- if .Keywords}}
      <div class="tags">
        {{- range .Keywords}}<span class="tag">{{.}}</span>{{end}}
      </div>
      {{- end}}
    </div>
    {{- end}}
{{- end}}
{{end}}

{{define "testimonial"}}
    <div class="testimonial">
      <blockquote>&#8220;{{.Quote}}&#8221;</blockquote>
      <cite>
        {{- if .URL}}<a href="{{.URL}}">{{.Name}} &middot; {{.Role}}</a>
        {{- else}}{{.Name}} &middot; {{.Role}}
        {{- end}}
      </cite>
    </div>
{{end}}
`
