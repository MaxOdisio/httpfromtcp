package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/maxodisio/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK         StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusError      StatusCode = 500
)

// Implementing the method String() allows fmt package to automatically call it if string is required (%s)
// AKA Stringer interface
func (sc StatusCode) String() string {
	switch sc {
	case StatusOK:
		return "OK"
	case StatusBadRequest:
		return "Bad Request"
	case StatusError:
		return "Internal Server Error"
	default:
		return "Unknown Status Code"
	}
}

type Writer struct {
	writer io.Writer
	state  writerState
}

type writerState int

const (
	writerStateInit writerState = iota
	writerStateStatusWritten
	writerStateHeadersWritten
	writerStateDone
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateInit {
		return fmt.Errorf("status line already written or invalid sequence.")
	}

	_, err := fmt.Fprintf(w.writer, "HTTP/1.1 %d %s\r\n", statusCode, statusCode)
	if err == nil {
		w.state = writerStateStatusWritten
	}

	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.state != writerStateStatusWritten {
		return fmt.Errorf("you have to write status line before headers.")
	}

	var r []byte

	for k, val := range h {
		r = append(r, k...)
		r = append(r, ": "...)
		r = append(r, val...)
		r = append(r, "\r\n"...)
	}

	r = append(r, "\r\n"...)

	_, err := w.writer.Write(r)
	if err == nil {
		w.state = writerStateHeadersWritten
	}
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateHeadersWritten && w.state != writerStateDone {
		return 0, fmt.Errorf("you have to write headers before body.")
	}
	w.state = writerStateDone

	b, err := w.writer.Write(p)
	return b, err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}
