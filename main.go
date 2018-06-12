package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var pkgInfo *build.Package

type Err struct {
	Name    string
	Code    string
	Comment string
}

type config struct {
	typeNames  string
	exportPath string
	apiVersion string
}

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func realMain() error {
	var (
		typeNames  = flag.String("type", "", "Required, multiple types with ','")
		exportPath = flag.String("export", "", "Options, export path")
		apiVersion = flag.String("version", "1.0", "Options, api version")
	)

	flag.Parse()
	if len(*typeNames) == 0 {
		log.Fatal("-type Required")
	}

	cfg := &config{
		typeNames:  *typeNames,
		exportPath: *exportPath,
		apiVersion: *apiVersion,
	}

	consts := cfg.ParsePackage()

	for _, v := range consts {
		src := cfg.genString(v)

		outputName := filepath.Join(".", fmt.Sprintf("%s_custom_string.go", pkgInfo.Name))
		if len(*exportPath) != 0 {
			p := strings.ToLower(*exportPath)
			os.MkdirAll(p, os.ModePerm)
			outputName = filepath.Join(p, fmt.Sprintf("%s_custom_string.go", pkgInfo.Name))
		}

		err := ioutil.WriteFile(outputName, src, 0644)
		if err != nil {
			log.Fatalf("writing output: %s", err)
		}
	}

	return nil
}

func (c *config) ParsePackage() map[string][]Err {
	mtypes := strings.Split(c.typeNames, ",")
	typesMap := make(map[string][]Err, len(mtypes))
	for _, v := range mtypes {
		typesMap[strings.TrimSpace(v)] = []Err{}
	}

	var err error
	fset := token.NewFileSet()
	pkgInfo, err = build.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}
	aps, err := parser.ParseDir(fset, pkgInfo.Dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_string.go")
	}, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range mtypes {

		var fs []*ast.File
		for _, f := range aps[pkgInfo.Name].Files {
			fs = append(fs, f)
		}

		conf := types.Config{Importer: importer.Default()}
		pkg, err := conf.Check("cmd/hello", fset, fs, nil)
		if err != nil {
			log.Fatal(err)
		}

		var m []Err
		for _, f := range fs {
			for _, v := range f.Scope.Objects {
				if v.Kind == ast.Con {
					d := v.Decl.(*ast.ValueSpec)
					m = append(m, Err{
						Name:    v.Name,
						Code:    pkg.Scope().Lookup(v.Name).(*types.Const).Val().String(),
						Comment: strings.TrimSuffix(d.Comment.Text(), "\n"),
					})
				}
			}
		}

		typesMap[t] = m
	}

	return typesMap
}

func (c *config) genString(consts []Err) []byte {
	const strTmp = `
    // Automatically-generated file. Do not edit
	package {{.pkg}}

	{{range $c :=.consts}}
    //{{$c.Name}}
    type HTTPError{{$c.Code}} struct {
	    ApiVersion string  ` + "`json:\"apiVersion\" example:\"{{$.version}}\"`" + `
	    Error  struct {
		    Code    int   ` + "`json:\"code\" example:\"{{$c.Code}}\"`" + `
		    Message string ` + "`json:\"message\" example:\"{{$c.Comment}}\"`" + `
	    }  ` + "`json:\"error\"`" + `
    }
	{{end}}
	`

	sort.Slice(consts, func(i, j int) bool {
		return consts[i].Code < consts[j].Code
	})

	pkgName := pkgInfo.Name
	if len(c.exportPath) != 0 {
		pkgName = path.Base(c.exportPath)
	}

	data := map[string]interface{}{
		"pkg":     pkgName,
		"consts":  consts,
		"version": c.apiVersion,
	}

	t, err := template.New("").Parse(strTmp)
	if err != nil {
		log.Fatal(err)
	}
	buff := bytes.NewBufferString("")
	err = t.Execute(buff, data)
	if err != nil {
		log.Fatal(err)
	}

	src, err := format.Source(buff.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	return src
}
