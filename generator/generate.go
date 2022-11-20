package generator

import (
	"context"
	"github.com/dave/dst/decorator"
	"io"
)

type Generate interface {
	Process(ctx context.Context, out io.Writer) error
}

type GenericGenerator struct {
	handler Handler
	parser  Parser
}

func NewGenericGenerator(handler Handler, parser Parser) *GenericGenerator {
	return &GenericGenerator{handler: handler, parser: parser}
}

func (g *GenericGenerator) Process(ctx context.Context, out io.Writer) error {
	decoratedFile, err := g.parser.Get(ctx)
	if err != nil {
		return err
	}
	err = g.handler.Handle(ctx, decoratedFile)
	if err != nil {
		return err
	}
	err = decorator.Fprint(out, decoratedFile)
	if err != nil {
		return err
	}
	return nil
}
