package main

import (
	"context"
	"io"
	"os"

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
	var out io.Writer
	if args.DryRun {
		out = os.Stdout
	} else {
		f, err := os.Create(args.Output)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		out = f
	}
	var gen generator.Generate
	var handler generator.Handler
	switch {
	case args.CommentModeCmd != nil:
		handler = generator.NewCommentBased()
	}
	gen = generator.NewGenericGenerator(handler, parser)
	err := gen.Process(ctx, out)
	if err != nil {
		panic(err)
	}
}
