package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/maxodisio/httpfromtcp/internal/request"
	"github.com/maxodisio/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he *HandlerError) Write(w io.Writer) error {
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		return err
	}

	msgBytes := []byte(he.Message)
	h := response.GetDefaultHeaders(len(msgBytes))
	err = response.WriteHeaders(w, h)
	if err != nil {
		return err
	}

	_, err = w.Write(msgBytes)
	return err
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
		handler:  handler,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error parsing request: %v\n", err)
		return
	}

	buf := new(bytes.Buffer)
	handlerErr := s.handler(buf, req)
	if handlerErr != nil {
		err := handlerErr.Write(conn)
		if err != nil {
			log.Printf("Error writing handler error to connection: %v\n", err)
		}
		return
	}

	h := response.GetDefaultHeaders(buf.Len())

	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Printf("Error writing status line: %v\n", err)
		return
	}

	err = response.WriteHeaders(conn, h)
	if err != nil {
		log.Printf("Error writing headers: %v\n", err)
		return
	}

	_, err = buf.WriteTo(conn)
	if err != nil {
		log.Printf("Error writing body to connection: %v\n", err)
		return
	}
}
