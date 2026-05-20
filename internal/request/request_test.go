package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestFromReader(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test: Good GET Request line
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(err)
	require.NotNil(r)
	assert.Equal("GET", r.RequestLine.Method)
	assert.Equal("/", r.RequestLine.RequestTarget)
	assert.Equal("1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(err)
	require.NotNil(r)
	assert.Equal("GET", r.RequestLine.Method)
	assert.Equal("/coffee", r.RequestLine.RequestTarget)
	assert.Equal("1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	r, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(err)
	assert.Nil(r)
}

func TestParseRequestLine(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test: Good Request line
	r, err := parseRequestLine([]byte("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(err)
	assert.Equal("GET", r.Method)
	assert.Equal("/", r.RequestTarget)
	assert.Equal("1.1", r.HttpVersion)

	// Test: Good Request line with path
	r, err = parseRequestLine([]byte("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(err)
	assert.Equal("GET", r.Method)
	assert.Equal("/coffee", r.RequestTarget)
	assert.Equal("1.1", r.HttpVersion)

	// Test: Good POST Request with path
	r, err = parseRequestLine([]byte("POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(err)
	assert.Equal("POST", r.Method)
	assert.Equal("/coffee", r.RequestTarget)
	assert.Equal("1.1", r.HttpVersion)

	// Test: Invalid number of parts in request line
	_, err = parseRequestLine([]byte("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(err)
	_, err = parseRequestLine([]byte("GET HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(err)

	// Test: Invalid method (out of order) Request line
	r, err = parseRequestLine([]byte("/coffee GET HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(err)

	// Test: Invalid version in Request line
	r, err = parseRequestLine([]byte("GET /coffee HTTP/2\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(err)
}
