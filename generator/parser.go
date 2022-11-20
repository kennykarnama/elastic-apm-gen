package generator

import (
	"context"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"go/parser"
	"go/token"
)

type Parser interface {
	Get(ctx context.Context) (*dst.File, error)
}

type FileParser struct {
	infile string
}

func NewParser(infile string) *FileParser {
	return &FileParser{infile: infile}
}

func (fp *FileParser) Get(ctx context.Context) (*dst.File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fp.infile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	decoratedFile, err := decorator.DecorateFile(fset, file)
	if err != nil {
		return nil, err
	}
	return decoratedFile, nil
}
