package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

var (
	camelCase = flag.Bool("camelcase", false, "convert enum names into const names in camelcase")
)

var usage = func() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of goen:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgoen [flags] [file]\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetPrefix("goen: ")

	flag.Usage = usage
	flag.Parse()

	// FIXME: use it
	_ = camelCase

	args := flag.Args()

	log.Printf("args=%v", args)

	if len(args) == 0 {
		args = []string{"."}
	}

	filename := args[0]
	if isDir(filename) {
		log.Printf("dir is nor allowed: filename=%s", filename)
		return
	}

	log.Println("--->", filename)

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("parse file: err=%v", err)
		return
	}

	var groups []*ast.CommentGroup
	ast.Inspect(f, func(n ast.Node) bool {
		v, ok := n.(*ast.CommentGroup)
		if !ok {
			return true
		}
		groups = append(groups, v)
		return false
	})

	b := strings.Builder{}
	for i := range groups {
		b.WriteString(groups[i].Text())
		b.WriteRune('\n')
	}
	log.Println(b.String())
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Printf("get file stat: path=%s, error=%v", path, err)
		os.Exit(1)
	}
	return fileInfo.IsDir()
}
