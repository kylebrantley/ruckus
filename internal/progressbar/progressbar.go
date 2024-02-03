package progressbar

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Option func(p *ProgressBar)

type ProgressBar struct {
	current int
	total   int
	writer  io.Writer
	empty   int
}

func New(total int, options ...Option) *ProgressBar {
	p := &ProgressBar{
		total:   total,
		current: 0,
		writer:  os.Stdout,
		empty:   50,
	}

	for _, o := range options {
		o(p)
	}

	return p
}

// WithWriter sets the output writer (defaults to os.StdOut)
func WithWriter(w io.Writer) Option {
	return func(p *ProgressBar) {
		p.writer = w
	}
}

func (p *ProgressBar) Print(current int) {
	p.current = current
	percent := float64(p.current) / float64(p.total)
	filled := int(percent * 50)
	p.empty = 50 - filled

	if p.current == p.total {
		fmt.Fprintf(
			p.writer,
			"\r[%s%s] %d%%\n",
			strings.Repeat("■", filled),
			strings.Repeat(" ", p.empty),
			int(percent*100),
		)
	} else {
		fmt.Fprintf(
			p.writer,
			"\r[%s%s] %d%%",
			strings.Repeat("■", filled),
			strings.Repeat(" ", p.empty),
			int(percent*100),
		)
	}
}
