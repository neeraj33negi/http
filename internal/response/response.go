package response

import (
	"fmt"
	"io"

	"github.com/neeraj33negi/http/internal/headers"
)

type StatusCode int

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	reason, err := reasonPhraseFor(statusCode)
	if err != nil {
		return err
	}
	b := []byte(reason)
	b = fmt.Appendf(b, "\r\n")
	_, err = w.writer.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	b := []byte{}
	for _, k := range headers.FieldLines() {
		b = fmt.Appendf(b, "%s:%s\r\n", k, headers.Get(k))
	}
	b = fmt.Appendf(b, "\r\n")
	_, err := w.writer.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.writer.Write(p)
}

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

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	fmt.Println(string(string(p)))
	body := fmt.Sprintf("%x\r\n", len(p))
	_, err := w.WriteBody([]byte(body))
	if err != nil {
		return 0, err
	}
	return w.WriteBody(p)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	body := fmt.Sprintf("%x\r\n\r\n", 0)
	return w.WriteBody([]byte(body))
}
