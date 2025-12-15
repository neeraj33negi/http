package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test: Valid headers
func TestValidHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nApi: rkey\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "rkey", headers.Get("Api"))
	assert.Equal(t, 36, n)
	assert.True(t, done)
}

// Test: Valid headers Case Insensitive
func TestValidHeaderCaseInsensitive(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, 25, n)
	assert.True(t, done)
}

// Test: Invalid characters headers
func TestInvalidCharsInHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

// Test: Invalid spacing header
func TestInvalidSpacingHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

// Test: Multiple value headers
func TestMultipleValueHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nApi: rkey\r\nApi: okey\r\nApi: pkey\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "rkey,okey,pkey", headers.Get("Api"))
	assert.Equal(t, 58, n)
	assert.True(t, done)
}
