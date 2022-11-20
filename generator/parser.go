package generator

import (
	"context"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"go/parser"
	"go/token"
	"os"
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
	// create backup
	bak, err := os.OpenFile(fmt.Sprintf("%s.apm_gen.bak", fp.infile), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return nil, err
	}
	defer bak.Close()
	err = decorator.Fprint(bak, decoratedFile)
	if err != nil {
		return nil, err
	}
	return decoratedFile, nil
}
