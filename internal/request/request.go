package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/neeraj33negi/http/internal/headers"
)

type parseState int

const (
	Initalized parseState = iota
	StateHeaders
	StateError
	StateBody
	Done
)

var SUPPORTED_METHODS = []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"}
var ERR_MALFORMED_REQUEST = fmt.Errorf("malformed request line")
var ERR_HTTP_VERSION_UNSUPPORTED = fmt.Errorf("http version not supported")
var ERR_METHOD_UNSUPPORTED = fmt.Errorf("method not supported")
var ERR_REQUEST_STATE_ERR = fmt.Errorf("request parse error")
var ERR_INVALID_CONTENT_LENGTH = fmt.Errorf("invalid content-length")
var SEPARATOR = []byte("\r\n")

type RequestLine struct {
	HttpVersion   string
	Method        string
	RequestTarget string
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        []byte
	parseState  parseState
}

func newRequest() *Request {
	return &Request{
		parseState: Initalized,
		Headers:    headers.NewHeaders(),
	}
}

func (r *Request) Print() {
	fmt.Println("Request line:")
	fmt.Printf("  - Method: %s\n", r.RequestLine.Method)
	fmt.Printf("  - Target: %s\n", r.RequestLine.RequestTarget)
	fmt.Printf("  - Version: %s\n", r.RequestLine.HttpVersion)
	r.Headers.Print()
}

func (rl *RequestLine) ValidHTTP() bool {
	return rl.HttpVersion == "1.1"
}

func (rl *RequestLine) ValidMethod() bool {
	return slices.Contains(SUPPORTED_METHODS, rl.Method)
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		remainingData := data[read:]
		if len(remainingData) == 0 {
			break
		}
		switch r.parseState {
		case StateError:
			return 0, ERR_REQUEST_STATE_ERR
		case Initalized:
			requestLine, n, err := parseRequestLine(remainingData)
			if err != nil {
				r.parseState = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.parseState = StateHeaders
			r.RequestLine = *requestLine
			read += n
		case StateHeaders:
			n, done, err := r.Headers.Parse(remainingData)
			if err != nil {
				return 0, err
			}
			read += n
			if n == 0 {
				break outer
			}
			if done {
				if r.Headers.Get("content-length") == "" || r.Headers.Get("content-length") == "0" {
					r.parseState = Done
				} else {
					r.parseState = StateBody
				}
			}
		case StateBody:
			contentLength, err := strconv.Atoi(r.Headers.Get("content-length"))
			if err != nil {
				return 0, ERR_INVALID_CONTENT_LENGTH
			}
			if contentLength == 0 {
				r.parseState = Done
				break
			}
			remainingBodyLength := min(contentLength-len(r.Body), len(remainingData))
			r.Body = append(r.Body, remainingData[:remainingBodyLength]...)
			read += remainingBodyLength
			if contentLength == len(r.Body) {
				r.parseState = Done
			}
			if contentLength < len(r.Body) {
				return 0, ERR_INVALID_CONTENT_LENGTH
			}
		case Done:
			break outer
		default:
			panic("solar flare probably")
		}
	}
	return read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buff := make([]byte, 1024)
	buffLen := 0
	for request.parseState != Done {
		n, err := reader.Read(buff[buffLen:])
		if err != nil {
			return nil, errors.Join(fmt.Errorf("unable to readStream"), err)
		}
		buffLen += n
		buffN, err := request.parse(buff[:buffLen])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[buffN:buffLen])
		buffLen -= buffN
	}

	// request.Print()
	return request, nil
}

// Request line is first line of an HTTP request i.e. METHOD, Path, HttpVersion information line
func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	startLine := b[:idx]
	read := idx + len(SEPARATOR)

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

	return &rl, read, nil
}
