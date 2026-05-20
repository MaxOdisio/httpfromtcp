package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqLine, err := parseRequestLine(b)
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *reqLine}, nil
}

func parseRequestLine(b []byte) (*RequestLine, error) {
	i := bytes.Index(b, []byte("\r\n"))
	if i == -1 {
		return nil, errors.New("Wrong request, missing CRLF.")
	}
	line := string(b[:i])
	s := strings.Fields(line)
	if len(s) != 3 {
		return nil, errors.New("Wrong request line format.")
	}

	httpVersion := s[2]
	reqTarget := s[1]
	method := s[0]

	for _, r := range method {
		if r < 'A' || r > 'Z' {
			return nil, errors.New("Method must contain only capital alphabetic characters.")
		}
	}

	httpParts := strings.Split(httpVersion, "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, errors.New("Only HTTP 1.1 version supported.")
	}

	return &RequestLine{
		HttpVersion:   httpParts[1],
		RequestTarget: reqTarget,
		Method:        method,
	}, nil
}
