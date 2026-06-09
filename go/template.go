package main

const resumeTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Basics.Name}} - Resume</title>
` + fontsHead + `
<style>
` + cssVars + `

  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  html { font-size: 15px; }

  body {
    font-family: var(--sans);
    background: var(--page);
    color: var(--ink);
    line-height: 1.5;
` + cssBodySmoothing + `
  }

  .page {
    max-width: 805px;
    margin: 0 auto;
    padding: 48px 48px 64px;
  }

  header {
    display: grid;
    grid-template-columns: 1fr 1fr;
    align-items: start;
    gap: 24px;
    padding-bottom: 10px;
    border-bottom: 1.5px solid var(--black);
    margin-bottom: 13px;
  }

  .title-label {
    font-size: 0.68rem;
    font-weight: 600;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: var(--accent);
    margin-bottom: 5px;
    width: fit-content;
    white-space: nowrap;
  }

  .name {
    font-family: var(--name-font);
    font-size: 2.5rem;
    font-weight: 400;
    line-height: 1;
    color: var(--black);
  }

  a { color: var(--muted); text-decoration: none; }
  a:hover { color: var(--accent); }

  .header-right {
    text-align: right;
  }

  .contact {
    font-size: 0.76rem;
    color: var(--muted);
    line-height: 1.85;
    margin-bottom: 0.5rem;
  }

  .profiles {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    justify-content: flex-end;
  }

  .profile-link {
    font-size: 0.72rem;
    color: var(--muted);
    border: 1px solid var(--rule);
    border-radius: 3px;
    padding: 2px 8px;
  }

  .profile-url {
    display: none;
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
    content: '-';
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

  time { font: inherit; }
  time[title] { cursor: help; }
  .print-only { display: none; }
  .print-employer { display: none; }
  .screen-only { display: revert; }

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

    header {
      grid-template-columns: 1fr 1fr 1fr;
    }

    .header-right {
      display: contents;
    }

    .contact {
      grid-column: 2;
      text-align: center;
      margin-bottom: 0;
    }

    .profiles {
      grid-column: 3;
      display: block;
      text-align: right;
      justify-content: flex-end;
    }

    .profile-link {
      display: block;
      border: none;
      padding: 0;
      font-size: 0.72rem;
      margin-bottom: 2px;
    }

    .profile-url {
      display: inline;
    }

    .profiles .print-only {
      display: block;
      font-size: 0.72rem;
      margin-bottom: 2px;
    }

    .job            { page-break-inside: avoid; }
    .project        { page-break-inside: avoid; }
    .section-intro  { break-inside: avoid; page-break-inside: avoid; }
    .section-label  { break-after: avoid; page-break-after: avoid; }
    header          { page-break-after: avoid; }

    body { font-feature-settings: normal; }
    .screen-only { display: none; }
    .footer-grid { grid-template-columns: 1fr 2fr; gap: 17px; }

    .employer-header { display: none; }
    .employer-group .job { padding-left: 0; border-left: none; }
    .position-divider { border-top: 1px solid var(--rule); margin: 10px 0; }
    .print-employer { display: inline; }
  }

  @media (max-width: 600px) {
    .page { padding: 28px 20px 48px; }
    header { grid-template-columns: 1fr; }
    .header-right { text-align: left; }
    .profiles { justify-content: flex-start; }
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
    <div class="header-right">
      <div class="contact">
        <a href="mailto:{{.Basics.Email}}">{{.Basics.Email}}</a><br>
        {{.Basics.Phone}} &middot; {{.Basics.Location.City}}, {{.Basics.Location.Region}}<br>
      </div>
      <div class="profiles">
        {{- range .Basics.Profiles}}
        <a href="{{.URL}}" class="profile-link" title="{{.Network}}">
          <span class="profile-name">{{.Network}}</span>
          <span class="profile-url">: {{stripScheme .URL}}</span>
        </a>
        {{- end}}
        <span class="print-only"><a href="{{.Basics.URL}}">{{stripScheme .Basics.URL}}</a></span>
      </div>
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
      <div class="section-label">Skills</div>
    </div>
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
      <div class="section-label">Education &amp; Certifications</div>
      <div class="footer-grid">
      <div>
        {{- range .Education}}
        <div class="edu-degree">{{.StudyType}} {{.Area}}</div>
        <div class="edu-detail">
          {{- if .URL}}<a href="{{.URL}}">{{.Institution}}</a>{{else}}{{.Institution}}{{end}}
          {{- if .EndDate}} &middot; <time{{if fullDateRange .StartDate .EndDate}} title="{{fullDateRange .StartDate .EndDate}}"{{end}}>{{formatDate .EndDate}}</time>{{end}}
        </div>
        {{- end}}
      </div>
      <div class="cert-list screen-only">
        {{- range $i, $cert := .Certificates}}
        {{- if gt $i 0}} &nbsp;&middot;&nbsp; {{end}}
        {{- if $cert.URL}}<a href="{{$cert.URL}}"{{with certTitle $cert}} title="{{.}}"{{end}}>{{$cert.Name}}</a>{{else}}<span{{with certTitle $cert}} title="{{.}}"{{end}}>{{$cert.Name}}</span>{{end}}
        {{- with certPrintID $cert}}<span class="print-only"> ({{.}})</span>{{end}}
        {{- end}}
      </div>
      <div class="cert-list print-only">
        {{- range $i, $g := .CertGroups}}
        {{- if gt $i 0}} &nbsp;&middot;&nbsp; {{end}}
        {{- $g.Issuer}}: {{certGroupNames $g}}{{with certGroupID $g}} ({{.}}){{end}}
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
<script>
(function() {
  var nameEl  = document.querySelector('h1.name');
  var labelEl = document.querySelector('p.title-label');
  var colEl   = document.querySelector('header > div:first-child');

  function fitHeader() {
    if (!nameEl || !labelEl || !colEl) return;
    nameEl.style.fontSize   = '';
    labelEl.style.fontSize  = '';
    nameEl.style.whiteSpace = 'nowrap';

    var colW   = colEl.getBoundingClientRect().width;
    var htmlFs = parseFloat(getComputedStyle(document.documentElement).fontSize);
    if (colW <= 0 || htmlFs <= 0) { nameEl.style.whiteSpace = ''; return; }

    var nr = document.createRange();
    nr.selectNodeContents(nameEl);
    var nw = nr.getBoundingClientRect().width;
    if (nw > 0 && isFinite(colW / nw)) {
      nameEl.style.fontSize = ((parseFloat(getComputedStyle(nameEl).fontSize) * colW / nw) / htmlFs).toFixed(4) + 'rem';
    }
    // clear nowrap only after font-size is applied so the scaled text never re-wraps
    nameEl.style.whiteSpace = '';

    var lr = document.createRange();
    lr.selectNodeContents(labelEl);
    var lw = lr.getBoundingClientRect().width;
    if (lw > 0 && isFinite(colW / lw)) {
      labelEl.style.fontSize = ((parseFloat(getComputedStyle(labelEl).fontSize) * colW / lw) / htmlFs).toFixed(4) + 'rem';
    }
  }

  function resetFit() {
    if (nameEl)  { nameEl.style.fontSize = ''; nameEl.style.whiteSpace = ''; }
    if (labelEl) labelEl.style.fontSize = '';
  }

  document.fonts.ready.then(fitHeader);
  window.addEventListener('resize', fitHeader);
  window.addEventListener('beforeprint', resetFit);
  window.addEventListener('afterprint',  fitHeader);
})();
</script>
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
        <span class="job-dates"><time{{if fullDate .StartDate}} title="{{fullDate .StartDate}}"{{end}}>{{formatDate .StartDate}}</time> &#8211; <time{{if fullDate .EndDate}} title="{{fullDate .EndDate}}"{{end}}>{{formatDate .EndDate}}</time></span>
      </div>
      {{- range $j, $pos := .Positions}}
      {{- if gt $j 0}}
      <div class="position-divider"></div>
      {{- end}}
      <div class="job">
        <div class="job-header">
          <span class="job-title">{{$pos.Position}}</span>
          <span class="job-dates"><time{{if fullDate $pos.StartDate}} title="{{fullDate $pos.StartDate}}"{{end}}>{{formatDate $pos.StartDate}}</time> &#8211; <time{{if fullDate $pos.EndDate}} title="{{fullDate $pos.EndDate}}"{{end}}>{{formatDate $pos.EndDate}}</time></span>
        </div>
        {{- if $pos.Summary}}
        <div class="job-summary">{{$pos.Summary}}</div>
        {{- end}}
        <div class="job-meta"><span class="print-employer">{{$.DisplayName}}{{if or $pos.Location (ne $pos.Employer $.DisplayName)}} &middot; {{end}}</span>{{- if ne $pos.Employer $.DisplayName}}{{$pos.Employer}}{{if $pos.Location}} &middot; {{end}}{{end}}{{- $pos.Location}}</div>
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
        <span class="job-dates"><time{{if fullDate .StartDate}} title="{{fullDate .StartDate}}"{{end}}>{{formatDate .StartDate}}</time> &#8211; <time{{if fullDate .EndDate}} title="{{fullDate .EndDate}}"{{end}}>{{formatDate .EndDate}}</time></span>
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
