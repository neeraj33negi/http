package headers

import (
	"bytes"
	"fmt"
)

var CRLF = []byte("\r\n")
var MALFORMED_HEADER = fmt.Errorf("malformed field line")

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) < 2 {
		return "", "", MALFORMED_HEADER
	}

	key := parts[0]
	value := bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(key, []byte(" ")) {
		return "", "", MALFORMED_HEADER
	}
	return string(key), string(value), nil
}

func (h Headers) Parse(b []byte) (int, bool, error) {
	read := 0
	done := false
	for !done {
		idx := bytes.Index(b[read:], CRLF)
		if idx == -1 {
			break
		}

		if idx == 0 {
			read += len(CRLF)
			done = true
			break
		}
		k, v, err := parseHeader(b[read : read+idx])
		if err != nil {
			return 0, done, err
		}
		h[k] = v
		read += idx + len(CRLF)
	}

	return read, done, nil
}
