package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Enums   []EnumParam
	Package string
}

type EnumParam struct {
	Name   string
	Type   string
	Values []string
}

type EnumFile struct {
	Package string
	Enums   []Enum
}

type Enum struct {
	Name   string
	Type   string
	Values []Value
}

type Value struct {
	Orig string
	Name string
}

const (
	formatIota        = "%s%s %s = iota"
	formatIotaPlusOne = "%s%s %s = iota + 1"
	format            = "%s%s"
)

var (
	plusOne = flag.Bool("one", false, "iota starts from one")
)

func main() {
	log.SetPrefix("goen: ")

	flag.Parse()

	srcPath := flag.Arg(0)
	if srcPath == "" {
		log.Fatal("[ERROR] source file path must not be empty")
	}

	dstPath := flag.Arg(1)
	if dstPath == "" {
		log.Fatal("[ERROR] destination file path must not be empty")
	}

	log.Printf("[INFO] src=%s, dst=%s", srcPath, dstPath)

	tmpl, err := template.New("enum").Parse(EnumTemplate)
	assertErr(err)

	cfg, err := readConfig(srcPath)
	assertErr(err)

	pkg, err := extractPackage(dstPath)
	assertErr(err)

	enumFile := prepare(pkg, cfg)

	err = saveTemplate(tmpl, enumFile, dstPath)
	assertErr(err)
}

func prepare(pkg string, cfg *Config) *EnumFile {
	eFile := &EnumFile{
		Package: pkg,
		Enums:   make([]Enum, len(cfg.Enums)),
	}
	for i := range cfg.Enums {
		t := cfg.Enums[i].Type
		eFile.Enums[i].Name = cfg.Enums[i].Name
		eFile.Enums[i].Type = t
		eFile.Enums[i].Values = make([]Value, len(cfg.Enums[i].Values))

		// TODO: to use different values depending on the specified type:
		//	- strings
		//	- numerics (iouta)
		for j := range cfg.Enums[i].Values {
			n := toCamelCase(cfg.Enums[i].Values[j])
			if j == 0 {
				n = fmt.Sprintf(getIotaFormat(), eFile.Enums[i].Name, n, t)
			} else {
				n = fmt.Sprintf(format, eFile.Enums[i].Name, n)
			}
			eFile.Enums[i].Values[j] = Value{
				Orig: cfg.Enums[i].Values[j],
				Name: n,
			}
		}
	}
	return eFile
}

func toCamelCase(name string) string {
	b := strings.Builder{}
	ss := strings.Split(name, "_")
	for i := range ss {
		runes := []rune(ss[i])
		for j := 0; j < len(runes); j++ {
			if j == 0 {
				b.WriteRune(runes[j])
			} else {
				b.WriteRune(unicode.ToLower(runes[j]))
			}
		}
	}
	return b.String()
}

func extractPackage(dstPath string) (dir string, err error) {
	// FIXME: need check.
	dir, _ = path.Split(dstPath)
	if dir == "" {
		dir = "main"
	} else {
		dir = strings.TrimSuffix(dir, "/")
	}
	return
}

func saveTemplate(tmpl *template.Template, e *EnumFile, dstPath string) (err error) {
	_, err = os.Stat(dstPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		osErr := os.Remove(dstPath)
		if osErr != nil {
			return osErr
		}
	}

	var f *os.File
	f, err = os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer func() {
		fErr := f.Close()
		if fErr != nil {
			err = fErr
		}
	}()

	return tmpl.Execute(f, e)
}

func readConfig(path string) (*Config, error) {
	fBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var c *Config
	err = yaml.Unmarshal(fBytes, &c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml config: %w", err)
	}
	return c, nil
}

func assertErr(err error) {
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
}

func getIotaFormat() string {
	if *plusOne {
		return formatIotaPlusOne
	}
	return formatIota
}
