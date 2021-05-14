package main

const EnumTemplate = `package {{.Package}}
{{range .Enums}}
type {{.Name}} {{.Type}}

const (
	{{- range .Values}}
	{{ .Name -}}
	{{- end}}
)
{{end}}
`
