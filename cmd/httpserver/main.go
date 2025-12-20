package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/neeraj33negi/http/internal/request"
	"github.com/neeraj33negi/http/internal/response"
	"github.com/neeraj33negi/http/internal/server"
)

const port = 42069

func main() {
	handlerFunc := func(w *response.Writer, r *request.Request) {
		headers := response.GetDefaultHeaders(0)
		var body []byte
		if r.RequestLine.RequestTarget == "/foo" {
			w.WriteStatusLine(response.StatusOk)
			body = []byte("bar")
		} else if r.RequestLine.RequestTarget == "/bar" {
			w.WriteStatusLine(response.StatusOk)
			body = []byte("baz")
		} else {
			w.WriteStatusLine(response.BadRequest)
			body = []byte("<html>.\n\n.no html for you")
		}
		headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeaders(headers)
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
