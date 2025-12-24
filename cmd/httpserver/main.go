package main

import (
	"crypto/sha256"
	"fmt"

	// "io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/neeraj33negi/http/internal/headers"
	"github.com/neeraj33negi/http/internal/request"
	"github.com/neeraj33negi/http/internal/response"
	"github.com/neeraj33negi/http/internal/server"
)

const port = 42069

func main() {
	handlerFunc := func(w *response.Writer, r *request.Request) {
		h := response.GetDefaultHeaders(0)
		var body []byte
		if r.RequestLine.RequestTarget == "/foo" {
			w.WriteStatusLine(response.StatusOk)
			body = []byte("bar")
		} else if r.RequestLine.RequestTarget == "/bar" {
			w.WriteStatusLine(response.StatusOk)
			body = []byte("baz")
		} else if strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin/stream") {
			target := r.RequestLine.RequestTarget
			res, err := http.Get("https://httpbin.org/" + target[len("/httpbin"):])
			if err != nil {
				w.WriteStatusLine(response.InternalServerError)
				h.Put("Content-Type", "text/html")
				body = []byte(fmt.Sprintf("<html><body>%s<body></html>", err.Error()))
			} else {
				w.WriteStatusLine(response.StatusOk)
				h.Delete("Content-Lenght")
				h.Put("Transfer-Encoding", "chunked")
				h.Put("Content-Type", "text/plain")
				h.Put("Trailer", "X-Content-SHA256,X-Content-Length")
				w.WriteHeaders(h)
				fullBody := []byte{}
				for {
					data := make([]byte, 32)
					_, err := res.Body.Read(data)
					if err != nil {
						break
					}
					w.WriteChunkedBody(data)
					fullBody = fmt.Appendf(fullBody, "%s", string(data))
				}
				checksumHash := sha256.Sum256(fullBody)
				hashStr := ""
				for _, b := range checksumHash {
					hashStr += string(b)
				}
				trailers := headers.NewHeaders()
				trailers.Put("X-Content-SHA256", hashStr)
				trailers.Put("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
				w.WriteHeaders(*trailers)
				w.WriteChunkedBodyDone()
				return
			}
		} else if r.RequestLine.RequestTarget == "/video" {
			h.Replace("Content-Type", "video/mp4")
			body, err := os.ReadFile("assets/vim.mp4")
			if err != nil {
				w.WriteStatusLine(response.InternalServerError)
			} else {
				w.WriteStatusLine(response.StatusOk)
			}
			h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
			w.WriteHeaders(h)
			w.WriteBody(body)
			return
		} else {
			w.WriteStatusLine(response.BadRequest)
			body = []byte("<html>.\n\n.no html for you")
		}
		h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeaders(h)
		w.WriteBody(body)
	}
	server, err := server.Serve(port, handlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
