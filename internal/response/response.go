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
	writerStateChunking
	writerStateDone
)

const crlf = "\r\n"

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

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != writerStateHeadersWritten && w.state != writerStateChunking {
		return 0, fmt.Errorf("invalid state for writing chunked body.")
	}
	if len(p) == 0 {
		return 0, nil
	}

	h := fmt.Sprintf("%x%s", len(p), crlf)
	_, err := w.writer.Write([]byte(h))
	if err != nil {
		return 0, err
	}

	n, err := w.writer.Write(p)
	if err != nil {
		return n, err
	}

	_, err = w.writer.Write([]byte(crlf))
	if err != nil {
		return n, err
	}

	w.state = writerStateChunking
	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != writerStateHeadersWritten && w.state != writerStateChunking {
		return 0, fmt.Errorf("invalid state for completing chunked body")
	}

	doneChunk := fmt.Sprintf("0%s%s", crlf, crlf)
	n, err := w.writer.Write([]byte(doneChunk))
	if err == nil {
		w.state = writerStateDone
	}

	return n, err
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
