package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// StdinReader manages input with IFS-based tokenisation.
type StdinReader struct {
	interp  *Interpreter
	scanner *bufio.Scanner
	file    *os.File
	buf     []string // pending tokens
}

func newStdinReader(interp *Interpreter) *StdinReader {
	r := &StdinReader{interp: interp}
	r.scanner = bufio.NewScanner(os.Stdin)
	return r
}

// NextToken returns the next input token (split by IFS or whitespace).
func (r *StdinReader) NextToken(ifs string) string {
	for len(r.buf) == 0 {
		line, ok := r.readLine()
		if !ok {
			return ""
		}
		r.buf = splitByIFS(line, ifs)
	}
	tok := strings.TrimSpace(r.buf[0])
	r.buf = r.buf[1:]
	return tok
}

func (r *StdinReader) readLine() (string, bool) {
	if r.scanner.Scan() {
		return r.scanner.Text(), true
	}
	return "", false
}

// EOF returns true if stdin is exhausted.
func (r *StdinReader) EOF() bool {
	return !r.scanner.Scan()
}

func splitByIFS(line, ifs string) []string {
	if ifs == "" {
		return []string{line}
	}
	parts := strings.Split(line, ifs)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// SetFile redirects input to a file.
func (r *StdinReader) SetFile(path string) error {
	if path == "" {
		// Restore stdin
		if r.file != nil {
			_ = r.file.Close()
			r.file = nil
		}
		r.scanner = bufio.NewScanner(os.Stdin)
		r.buf = nil
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	if r.file != nil {
		_ = r.file.Close()
	}
	r.file = f
	r.scanner = bufio.NewScanner(f)
	r.buf = nil
	return nil
}

// StdoutWriter manages output with optional file redirection.
type StdoutWriter struct {
	w    io.Writer
	file *os.File
}

func newStdoutWriter() *StdoutWriter {
	return &StdoutWriter{w: os.Stdout}
}

func (w *StdoutWriter) Write(b []byte) {
	w.w.Write(b) //nolint:errcheck
}

func (w *StdoutWriter) WriteString(s string) {
	_, _ = fmt.Fprint(w.w, s)
}

// SetFile redirects output to a file.
func (w *StdoutWriter) SetFile(path, mode string) error {
	if path == "" {
		if w.file != nil {
			_ = w.file.Close()
			w.file = nil
		}
		w.w = os.Stdout
		return nil
	}
	flag := os.O_WRONLY | os.O_CREATE
	if mode == "at" {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return err
	}
	if w.file != nil {
		_ = w.file.Close()
	}
	w.file = f
	w.w = f
	return nil
}
