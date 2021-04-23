package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Enums   []Enum
	Package string
}

type Enum struct {
	Name   string
	Type   string
	Values []string
}

const enumTemplate = `package {{.Package}}
{{range .Enums}}
type {{.Name}} {{.Type}}

const (
	{{- range .Values}}
	{{ . -}}
	{{- end}}
)
{{end}}
`

var (
	marshal = flag.Bool("marshal", false, "marshal")
)

func main() {
	log.SetPrefix("goen: ")
	flag.Parse()

	fromFilePath := flag.Args()[0]
	toFilePath := flag.Args()[1]
	if toFilePath == "" {
		log.Fatal("[ERROR] to file is empty")
	}

	fb, err := os.ReadFile(fromFilePath)
	assertErr(err)

	var cfg *Configuration
	err = yaml.Unmarshal(fb, &cfg)
	assertErr(err)

	cfg.Package, _ = path.Split(toFilePath)
	cfg.Package = strings.ReplaceAll(cfg.Package, "/", "")

	for i := range cfg.Enums {
		for j := range cfg.Enums[i].Values {
			name := cfg.Enums[i].Name
			if j == 0 {
				cfg.Enums[i].Values[j] = fmt.Sprintf("%s%s %s = iota", name, cfg.Enums[i].Values[j], name)
			} else {
				cfg.Enums[i].Values[j] = fmt.Sprintf("%s%s", name, cfg.Enums[i].Values[j])
			}
		}
	}

	tmpl, err := template.New("enums").Parse(enumTemplate)
	assertErr(err)

	err = os.Remove(toFilePath)
	assertErr(err)

	toF, err := os.OpenFile(toFilePath, os.O_RDWR|os.O_CREATE, 0755)
	assertErr(err)

	err = tmpl.Execute(toF, cfg)
	assertErr(err)
}

func assertErr(err error) {
	if err != nil {
		log.Fatalf("[ERROR] %w+", err)
	}
}
