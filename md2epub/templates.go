package main

import (
	"html/template"
)

var (
	tpage = template.Must(template.New("").Parse(pageTemplateText))
	tnav  = template.Must(template.New("").Parse(navTemplateText))
)

const (
	pageTemplateText = `<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="{{ if .lang }}{{ .lang }}{{ else }}en{{ end }}">
<head>
<meta charset="UTF-8" />
<title>{{ .title }}</title>
</head>
<body>
{{ if .title }}<h1>{{ .title }}</h1>

{{ end }}{{ .content }}
</body>
</html>`

	navTemplateText = `<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="{{ if .lang }}{{ .lang }}{{ else }}en{{ end }}">
<head>
<meta charset="UTF-8" />
<title>{{ .title }}</title>
</head>
<body>
<nav epub:type="toc">
<h1>{{ .title }}</h1>
<ol>{{ range .toc }}{{ if and .Title .Spine .Filename}}
<li><a href="{{ .Filename }}">{{ .Title }}</a></li>{{ end }}{{ end }}
</ol>
</nav>
</body>
</html>`
)
