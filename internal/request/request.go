package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

const bufferSize = 8
const crlf = "\r\n"

type requestState int

const (
	requestInitialized requestState = iota
	requestDone
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestInitialized:
		n, reqLine, err := parseRequestLine(data)
		if err != nil {
			// something actually went wrong
			return 0, err
		}
		if n == 0 {
			// just need more data
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.state = requestDone
		return n, nil
	case requestDone:
		return 0, errors.New("Error: trying to read data in a done state.")
	default:
		return 0, errors.New("Error: uknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		state: requestInitialized,
	}

	for req.state != requestDone {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		readToIndex += numBytesRead
		if err != nil {
			if err == io.EOF {
				req.parse(buf[:readToIndex]) // fa il parse dei dati rimasti
				if req.state != requestDone {
					return nil, errors.New("incomplete request: unexpected EOF")
				}
				break
			}
			return nil, err
		}

		bytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[bytesParsed:])
		readToIndex -= bytesParsed
	}
	return req, nil
}

func parseRequestLine(b []byte) (int, *RequestLine, error) {
	i := bytes.Index(b, []byte(crlf))
	if i == -1 {
		return 0, nil, nil
	}
	line := string(b[:i])
	reqLine, err := requestLineFromString(line)
	if err != nil {
		return 0, nil, err
	}
	return i + 2, reqLine, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	s := strings.Fields(str)
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
