package response

import (
	"fmt"
	"io"

	"github.com/neeraj33negi/http/internal/headers"
)

type StatusCode int

const (
	StatusOk            = 200
	BadRequest          = 400
	InternalServerError = 500
)

var reasonPhrases = map[StatusCode]string{
	StatusOk:            fmt.Sprintf("HTTP/1.1 %d OK", StatusOk),
	BadRequest:          fmt.Sprintf("HTTP/1.1 %d Bad Request", BadRequest),
	InternalServerError: fmt.Sprintf("HTTP/1.1 %d Internal Server Error", InternalServerError),
}

func reasonPhraseFor(statuscode StatusCode) (string, error) {
	reason := reasonPhrases[statuscode]
	if len(reason) == 0 {
		return "", fmt.Errorf("invalid status code %d", statuscode)
	}
	return reason, nil
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reason, err := reasonPhraseFor(statusCode)
	if err != nil {
		return err
	}
	b := []byte(reason)
	b = fmt.Appendf(b, "\r\n")
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Put("Content-Length", fmt.Sprint(contentLen))
	headers.Put("Connection", "Close")
	headers.Put("Content-Type", "text/plain")
	return *headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	b := []byte{}
	for _, k := range headers.FieldLines() {
		b = fmt.Appendf(b, "%s:%s\r\n", k, headers.Get(k))
	}
	b = fmt.Appendf(b, "\r\n")
	_, err := w.Write(b)
	if err != nil {
		return err
	}
	return nil
}
