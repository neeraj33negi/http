package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

var CRLF = []byte("\r\n")
var MALFORMED_HEADER = fmt.Errorf("malformed field line")
var MALFORMED_HEADER_NAME = fmt.Errorf("malformed field name")
var validTokens = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Print() {
	fmt.Println("Headers:")
	for k, v := range h.headers {
		fmt.Printf("  - %s: %s\n", k, v)
	}
}

func (h *Headers) FieldLines() []string {
	var fieldLines []string
	for k := range h.headers {
		fieldLines = append(fieldLines, k)
	}
	return fieldLines
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Put(name, value string) {
	h.headers[strings.ToLower(name)] = value
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

func (h *Headers) Parse(b []byte) (int, bool, error) {
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
		if !validFieldName(k) {
			return 0, done, MALFORMED_HEADER_NAME
		}
		if len(h.Get(k)) == 0 {
			h.Put(k, v)
		} else {
			h.Put(k, fmt.Sprintf("%s,%s", h.Get(k), v))
		}
		read += idx + len(CRLF)
	}
	return read, done, nil
}

func validFieldName(s string) bool {
	if len(s) < 1 {
		return false
	}
	valid := true
	for _, ch := range s {
		if (ch >= 'A' && ch <= 'Z') ||
			(ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') ||
			slices.Contains(validTokens, byte(ch)) {
			continue
		} else {
			valid = false
		}
		if !valid {
			break
		}
	}
	return valid
}
