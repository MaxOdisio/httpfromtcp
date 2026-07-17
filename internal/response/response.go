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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", statusCode, statusCode)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	var r []byte

	for k, val := range h {
		r = append(r, k...)
		r = append(r, ": "...)
		r = append(r, val...)
		r = append(r, "\r\n"...)
	}

	r = append(r, "\r\n"...)

	_, err := w.Write(r)
	return err
}
