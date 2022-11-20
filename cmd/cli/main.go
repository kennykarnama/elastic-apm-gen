package main

import (
	"context"
	"github.com/alexflint/go-arg"
	"github.com/kennykarnama/elastic-apm-gen/generator"
)

type CommentModeCmd struct{}

var args struct {
	Input          string          `arg:"-i,--input,required" help:"input. can be file or directory"`
	CommentModeCmd *CommentModeCmd `arg:"subcommand:comment-mode" help:"do generation by comment"`
	DryRun         bool            `arg:"--dry-run" help:"if true, it will print to stdout"`
	Output         string          `arg:"-o,--output" help:"output"`
}

func main() {
	arg.MustParse(&args)
	ctx := context.Background()
	parser := generator.NewParser(args.Input)
	var out generator.Flusher
	if args.DryRun {
		out = &generator.StdOutFlusher{}
	} else {
		if args.Output == "" {
			args.Output = args.Input
		}
		out = generator.NewFileFlusher(args.Output)
	}
	var gen generator.Generate
	var handler generator.Handler
	switch {
	case args.CommentModeCmd != nil:
		handler = generator.NewCommentBased()
	}
	gen = generator.NewGenericGenerator(handler, parser, out)
	err := gen.Process(ctx)
	if err != nil {
		panic(err)
	}
}
