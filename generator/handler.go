package generator

import (
	"context"
	"fmt"
	"github.com/dave/dst"
	"go/token"
	"strings"
)

type Handler interface {
	Handle(ctx context.Context, decoratedFile *dst.File) error
}

type CommentBasedHandler struct{}

func NewCommentBased() *CommentBasedHandler {
	return &CommentBasedHandler{}
}

func (c *CommentBasedHandler) Handle(ctx context.Context, decoratedFile *dst.File) error {
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
		addImports(decoratedFile, imports)
	}
	return nil
}
