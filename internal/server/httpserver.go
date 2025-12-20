package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/neeraj33negi/http/internal/request"
	"github.com/neeraj33negi/http/internal/response"
)

type HandlerError struct {
	Message string
}

type Handler func(w io.Writer, r *request.Request) *HandlerError

type Server struct {
	closed  bool
	handler Handler
}

func (s *Server) handle(conn net.Conn) {
	// harcode for testing
	defer conn.Close()
	writer := bytes.NewBuffer([]byte{})
	h := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.BadRequest)
		response.WriteHeaders(conn, h)
		return
	}
	handlerError := s.handler(writer, r)
	if handlerError != nil {
		response.WriteStatusLine(conn, response.InternalServerError)
		response.WriteHeaders(conn, h)
		return
	}
	response.WriteStatusLine(conn, 200)
	body := writer.Bytes()
	h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
	response.WriteHeaders(conn, h)
	conn.Write(body)
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
