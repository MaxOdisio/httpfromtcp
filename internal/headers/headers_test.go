package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(err)
	require.NotNil(headers)
	assert.Equal("localhost:42069", headers["host"])
	assert.Equal(23, n)
	assert.False(done)

	// Test: Valid single header with capital letters on key
	headers = NewHeaders()
	data = []byte("HoSt: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(err)
	require.NotNil(headers)
	assert.Equal("localhost:42069", headers["host"])
	assert.Equal(23, n)
	assert.False(done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host:            localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(err)
	require.NotNil(headers)
	assert.Equal("localhost:42069", headers["host"])
	assert.Equal(34, n)
	assert.False(done)

	// Test: Valid 2 headers with existing headers
	headers = map[string]string{"host": "localhost:42069"}
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(err)
	require.NotNil(headers)
	assert.Equal("localhost:42069", headers["host"])
	assert.Equal("curl/7.81.0", headers["user-agent"])
	assert.Equal(25, n)
	assert.False(done)

	// Test: Valid 2 headers with same existing key
	headers = map[string]string{"set-person": "lane-loves-go"}
	data = []byte("Set-Person: prime-loves-zig\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(err)
	require.NotNil(headers)
	assert.Equal("lane-loves-go, prime-loves-zig", headers["set-person"])
	assert.Equal(29, n)
	assert.False(done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(err)
	assert.Equal(2, n)
	assert.True(done)

	// Test: Invalid spacing header before key
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(err)
	assert.Equal(0, n)
	assert.False(done)

	// Test: Invalid spacing header between key and colon
	headers = NewHeaders()
	data = []byte("Host : localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(err)
	assert.Equal(0, n)
	assert.False(done)

	// Test: Invalid character in key
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(err)
	assert.Equal(0, n)
	assert.False(done)
}
