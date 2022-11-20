package generator

import (
	"context"
)

type Generate interface {
	Process(ctx context.Context) error
}

type GenericGenerator struct {
	handler Handler
	parser  Parser
	flusher Flusher
}

func NewGenericGenerator(handler Handler, parser Parser, flusher Flusher) *GenericGenerator {
	return &GenericGenerator{handler: handler, parser: parser, flusher: flusher}
}

func (g *GenericGenerator) Process(ctx context.Context) error {
	decoratedFile, err := g.parser.Get(ctx)
	if err != nil {
		return err
	}
	err = g.handler.Handle(ctx, decoratedFile)
	if err != nil {
		return err
	}
	err = g.flusher.Write(ctx, decoratedFile)
	if err != nil {
		return err
	}
	return nil
}
