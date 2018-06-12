package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"
)

var (
	write = flag.Bool("w", false, "write out to file (default stdout)")
)

func main() {

	flag.Parse()

	filepath.Walk(flag.Arg(0), func(path string, stat os.FileInfo, err error) error {
		if stat.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		// strip imports
		ast.FilterFile(f, func(_ string) bool { return true })

		var buf bytes.Buffer
		if err := format.Node(&buf, fset, f); err != nil {
			panic(err)
		}

		// re-set imports
		data, err := imports.Process(path, buf.Bytes(), nil)
		if err != nil {
			panic(err)
		}

		if *write {
			if err := ioutil.WriteFile(path, data, stat.Mode()); err != nil {
				panic(err)
			}
		} else {
			fmt.Println("# ", path)
			fmt.Println(string(data))
		}

		return nil

	})

}
