package main

import (
	"fmt"
	"github.com/kennykarnama/elastic-apm-gen/helper"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

var args struct {
	InputFile string `arg:"-i,--input-file" help:"input file"`
	DryRun    bool   `arg:"--dry-run" help:"if true, it will print to stdout"`
}

func main() {
	arg.MustParse(&args)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, args.InputFile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	decoratedFile, err := decorator.DecorateFile(fset, file)
	if err != nil {
		panic(err)
	}
	for _, imp := range decoratedFile.Imports {
		fmt.Println(imp.Name, imp.Path)
	}
	hasApm := false
	for _, f := range decoratedFile.Decls {
		fn, ok := f.(*dst.FuncDecl)
		exist := false
		var params []string
		if ok {
			fmt.Println("found function", fn.Name)
			// check if contains comment
			for _, stmt := range fn.Body.List {
				for _, dec := range stmt.Decorations().Start {
					if strings.Contains(dec, "// apm:startSpan") {
						hasApm = true
						exist = true
						params = strings.Split(dec, " ")
						break
					}
				}
			}
		}
		if exist && len(params) == 4 {
			fmt.Printf("Found APM params=%v\n", params[2:])

			a1 := dst.AssignStmt{
				// token.DEFINE is :=
				Tok: token.DEFINE,
				// left hand side has two identifiers, span and ctx
				Lhs: []dst.Expr{
					&dst.Ident{Name: "span"},
					&dst.Ident{Name: "ctx"},
				},
				// right hand is a call to function
				Rhs: []dst.Expr{
					&dst.CallExpr{
						// function is taken from a module 'tracing' by it's name
						Fun: &dst.SelectorExpr{
							X:   &dst.Ident{Name: "apm"},
							Sel: &dst.Ident{Name: "StartSpan"},
						},
						// function has two arguments
						Args: []dst.Expr{
							&dst.Ident{Name: "ctx"},
							// c.App.Context
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf(`"%s"`, params[2]),
							},
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf(`"%s"`, params[3]),
							},
						},
					},
				},
			}
			a3 := dst.DeferStmt{
				// what function call should be deferred?
				Call: &dst.CallExpr{
					// Finish from 'span' identifier
					Fun: &dst.SelectorExpr{
						X:   &dst.Ident{Name: "span"},
						Sel: &dst.Ident{Name: "End"},
					},
				},
			}
			fn.Body.List = append([]dst.Stmt{&a1, &a3}, fn.Body.List...)
		}
	}

	if hasApm {
		imp := "go.elastic.co/apm/v2"
		imports := map[string]string{
			"apm": imp,
		}
		helper.AddImports(decoratedFile, imports)
	}

	if err := decorator.Print(decoratedFile); err != nil {
		panic(err)
	}
	ToFile(decoratedFile)

}

func ToFile(decoratedFile *dst.File) {
	f, err := os.Create("refactored.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := decorator.Fprint(f, decoratedFile); err != nil {
		panic(err)
	}
}
