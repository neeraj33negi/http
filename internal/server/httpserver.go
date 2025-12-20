package server

import (
	"fmt"
	"net"

	"github.com/neeraj33negi/http/internal/response"
)

type Server struct {
	closed bool
}

func (s *Server) handle(conn net.Conn) {
	// harcode for testing
	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text.plain\r\nContent-Length: 13\r\n\r\nHello World!\n"))
	response.WriteStatusLine(conn, 200)
	body := []byte("Hello World!\n")
	h := response.GetDefaultHeaders(len(body))
	response.WriteHeaders(conn, h)
	conn.Close()
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

func Serve(port uint) (*Server, error) {
	s := &Server{false}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	go run(s, listener)
	return s, nil
}
