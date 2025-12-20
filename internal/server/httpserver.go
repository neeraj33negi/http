package server

import (
	"fmt"
	"net"

	"github.com/neeraj33negi/http/internal/request"
	"github.com/neeraj33negi/http/internal/response"
)

type HandlerError struct {
	Message string
}

type Handler func(w *response.Writer, r *request.Request)

type Server struct {
	closed  bool
	handler Handler
}

func (s *Server) handle(conn net.Conn) {
	// harcode for testing
	defer conn.Close()
	rWriter := response.NewWriter(conn)
	h := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		rWriter.WriteStatusLine(response.BadRequest)
		rWriter.WriteHeaders(h)
		return
	}

	s.handler(rWriter, r)
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}

func run(s *Server, listener net.Listener) error {
	for {
		if s.closed {
			return nil
		}

		conn, err := listener.Accept()
		if err != nil {
			return nil
		}
		go s.handle(conn)
	}
}

func Serve(port uint, h Handler) (*Server, error) {
	s := &Server{closed: false, handler: h}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	go run(s, listener)
	return s, nil
}
