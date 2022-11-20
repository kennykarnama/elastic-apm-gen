package generator

import (
	"context"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"os"
)

type Flusher interface {
	Write(ctx context.Context, decoratedFile *dst.File) error
}

type FileFlusher struct {
	out string
}

func NewFileFlusher(out string) *FileFlusher {
	return &FileFlusher{out: out}
}

func (ff *FileFlusher) Write(ctx context.Context, decoratedFile *dst.File) error {
	f, err := os.OpenFile(fmt.Sprintf(ff.out), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	return decorator.Fprint(f, decoratedFile)
}

type StdOutFlusher struct{}

func (sf *StdOutFlusher) Write(ctx context.Context, decoratedFile *dst.File) error {
	return decorator.Fprint(os.Stdout, decoratedFile)
}
