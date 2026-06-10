package request

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func TestRequestFromReader(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test: Good GET Request line
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(err)
	require.NotNil(r)
	assert.Equal("GET", r.RequestLine.Method)
	assert.Equal("/", r.RequestLine.RequestTarget)
	assert.Equal("1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader = &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err = RequestFromReader(reader)
	require.NoError(err)
	require.NotNil(r)
	assert.Equal("GET", r.RequestLine.Method)
	assert.Equal("/coffee", r.RequestLine.RequestTarget)
	assert.Equal("1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	r, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(err)
	assert.Nil(r)

	// Test: Standard Headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(err)
	require.NotNil(r)
	assert.Equal("localhost:42069", r.Headers["host"])
	assert.Equal("curl/7.81.0", r.Headers["user-agent"])
	assert.Equal("*/*", r.Headers["accept"])

	// Test: Malformed Header
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Error(err)

	// Test: Duplicate Headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nUser-Agent: my-app\r\nUser-Agent: curl/7.81.0\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(err)
	require.NotNil(r)
	assert.Equal("my-app, curl/7.81.0", r.Headers["user-agent"])

	// Test: Case Insensitive Headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nHoST: localhost:42069\r\nuser-AGEnt: curl/7.81.0\r\nAccePT: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(err)
	require.NotNil(r)
	assert.Equal("localhost:42069", r.Headers["host"])
	assert.Equal("curl/7.81.0", r.Headers["user-agent"])
	assert.Equal("*/*", r.Headers["accept"])

	// Test: Missing End of Headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Error(err)

	// Test: Empty Headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(err)
	require.NotNil(r)
	require.Empty(r.Headers)
	assert.Equal("GET", r.RequestLine.Method)
}

func TestParseRequestLine(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test: Good Request line
	_, r, err := parseRequestLine([]byte("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(err)
	assert.Equal("GET", r.Method)
	assert.Equal("/", r.RequestTarget)
	assert.Equal("1.1", r.HttpVersion)

	// Test: Good Request line with path
	_, r, err = parseRequestLine([]byte("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(err)
	assert.Equal("GET", r.Method)
	assert.Equal("/coffee", r.RequestTarget)
	assert.Equal("1.1", r.HttpVersion)

	// Test: Good POST Request with path
	_, r, err = parseRequestLine([]byte("POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(err)
	assert.Equal("POST", r.Method)
	assert.Equal("/coffee", r.RequestTarget)
	assert.Equal("1.1", r.HttpVersion)

	// Test: Invalid number of parts in request line
	_, _, err = parseRequestLine([]byte("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(err)
	_, _, err = parseRequestLine([]byte("GET HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(err)

	// Test: Invalid method (out of order) Request line
	_, r, err = parseRequestLine([]byte("/coffee GET HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(err)

	// Test: Invalid version in Request line
	_, r, err = parseRequestLine([]byte("GET /coffee HTTP/2\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(err)
}
