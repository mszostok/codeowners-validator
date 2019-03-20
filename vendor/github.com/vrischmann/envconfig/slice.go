package envconfig

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type sliceTokenizer struct {
	err      error
	r        *bufio.Reader
	buf      bytes.Buffer
	inBraces bool
}

var eof = rune(0)

func newSliceTokenizer(str string) *sliceTokenizer {
	return &sliceTokenizer{
		r: bufio.NewReader(strings.NewReader(str)),
	}
}

func (t *sliceTokenizer) scan() bool {
	for {
		if t.err == io.EOF && t.buf.Len() == 0 {
			return false
		}

		ch := t.readRune()
		if ch == eof {
			return true
		}

		if ch == '{' {
			t.inBraces = true
		}
		if ch == '}' {
			t.inBraces = false
		}

		if ch == ',' && !t.inBraces {
			return true
		}

		// NOTE(vincent): we ignore the WriteRune error here because there is NO WAY
		// for WriteRune to return an error.
		// Yep. Seriously. Look here http://golang.org/src/bytes/buffer.go?s=7661:7714#L227
		_, _ = t.buf.WriteRune(ch)
	}
}

func (t *sliceTokenizer) readRune() rune {
	ch, _, err := t.r.ReadRune()
	if err != nil {
		t.err = err
		return eof
	}

	return ch
}

func (t *sliceTokenizer) text() string {
	str := t.buf.String()
	t.buf.Reset()

	return str
}

func (t *sliceTokenizer) Err() error {
	if t.err == io.EOF {
		return nil
	}
	return t.err
}
