package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
)

type parseState int

const (
	initalized parseState = iota
	done
)

var SUPPORTED_METHODS = []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"}
var ERR_MALFORMED_REQUEST = fmt.Errorf("malformed request line")
var ERR_HTTP_VERSION_UNSUPPORTED = fmt.Errorf("http version not supported")
var ERR_METHOD_UNSUPPORTED = fmt.Errorf("method not supported")
var SEPARATOR = []byte("\r\n")

type RequestLine struct {
	HttpVersion   string
	Method        string
	RequestTarget string
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	parseState  parseState
}

func (rl *RequestLine) ValidHTTP() bool {
	return rl.HttpVersion == "1.1"
}

func (rl *RequestLine) ValidMethod() bool {
	return slices.Contains(SUPPORTED_METHODS, rl.Method)
}

func (r *Request) parse(data []byte) (int, error) {
	requestLine, n, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, nil
	}
	r.parseState = done
	r.RequestLine = *requestLine
	return n, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{parseState: initalized}
	buff := make([]byte, 1024)
	buffLen := 0
	for request.parseState != done {
		n, err := reader.Read(buff[buffLen:])
		if err != nil {
			return nil, errors.Join(fmt.Errorf("unable to readStream"), err)
		}
		buffLen += n
		buffN, err := request.parse(buff[:buffLen+n])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[buffN:buffLen])
		buffLen -= buffN
	}

	return &request, nil
}

// Request line is first line of an HTTP request i.e. METHOD, Path, HttpVersion information line
func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	startLine := b[:idx]
	// restMessage := b[idx+len(SEPARATOR):]

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERR_MALFORMED_REQUEST
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 {
		return nil, len(startLine), ERR_MALFORMED_REQUEST
	}
	rl := RequestLine{
		HttpVersion:   string(httpParts[1]),
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
	}
	if !rl.ValidHTTP() {
		return nil, len(startLine), ERR_HTTP_VERSION_UNSUPPORTED
	}

	if !rl.ValidMethod() {
		return nil, len(startLine), ERR_METHOD_UNSUPPORTED
	}

	return &rl, len(startLine), nil
}
