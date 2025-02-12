package m3u8

import (
	"bufio"
	"io"
	"iter"
)

type line = []byte

type Lines interface {
	Lines() iter.Seq[line]
}

type bufLines struct {
	*bufio.Scanner
}

func (b bufLines) Lines() iter.Seq[line] {
	return func(yield func(line) bool) {
		for b.Scan() {
			if !yield(b.Bytes()) {
				return
			}
		}
	}
}

func NewLines(r io.Reader) Lines {
	return bufLines{bufio.NewScanner(r)}
}

func AppendLines(sl []byte, lines Lines) []byte {
	for line := range lines.Lines() {
		sl = AppendLine(sl, line)
	}
	return sl
}
