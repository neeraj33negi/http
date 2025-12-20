package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/neeraj33negi/http/internal/request"
	"github.com/neeraj33negi/http/internal/server"
)

const port = 42069

func main() {
	handlerFunc := func(w io.Writer, r *request.Request) *server.HandlerError {
		if r.RequestLine.RequestTarget == "/foo" {
			w.Write([]byte("bar"))
		} else {
			w.Write([]byte("sup fool"))
		}
		return nil
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
