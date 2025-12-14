package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

var SUPPORTED_METHODS = []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"}

type RequestLine struct {
	HttpVersion   string
	Method        string
	RequestTarget string
}

func (rl *RequestLine) ValidHTTP() bool {
	return rl.HttpVersion == "1.1"
}

func (rl *RequestLine) ValidMethod() bool {
	return slices.Contains(SUPPORTED_METHODS, rl.Method)
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

var ERR_MALFORMED_REQUEST = fmt.Errorf("malformed request line")
var ERR_HTTP_VERSION_UNSUPPORTED = fmt.Errorf("http version not supported")
var ERR_METHOD_UNSUPPORTED = fmt.Errorf("method not supported")

var SEPARATOR = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to readStream"), err)
	}
	requestLine, _, err := parseRequestLine(string(req))
	if err != nil {
		return nil, err
	}
	return &Request{RequestLine: *requestLine}, nil
}

// Request line is first line of an HTTP request i.e. METHOD, Path, HttpVersion information line
func parseRequestLine(req string) (*RequestLine, string, error) {
	idx := strings.Index(req, SEPARATOR)
	if idx == -1 {
		return nil, req, ERR_MALFORMED_REQUEST
	}
	startLine := req[:idx]
	restMessage := req[idx+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, req, ERR_MALFORMED_REQUEST
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 {
		return nil, restMessage, ERR_MALFORMED_REQUEST
	}
	rl := RequestLine{
		HttpVersion:   httpParts[1],
		Method:        parts[0],
		RequestTarget: parts[1],
	}
	if !rl.ValidHTTP() {
		return nil, restMessage, ERR_HTTP_VERSION_UNSUPPORTED
	}

	if !rl.ValidMethod() {
		return nil, restMessage, ERR_METHOD_UNSUPPORTED
	}

	return &rl, restMessage, nil
}
