package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/maxodisio/httpfromtcp/internal/request"
	"github.com/maxodisio/httpfromtcp/internal/response"
	"github.com/maxodisio/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	problemHandler := func(w *response.Writer, req *request.Request) {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			w.WriteStatusLine(response.StatusBadRequest)
			body := respond400()
			h := response.GetDefaultHeaders(len(body))
			h.Replace("Content-Type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody(body)
			return
		case "/myproblem":
			w.WriteStatusLine(response.StatusError)
			body := respond500()
			h := response.GetDefaultHeaders(len(body))
			h.Replace("Content-Type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody(body)
			return
		default:
			w.WriteStatusLine(response.StatusOK)
			body := respond200()
			h := response.GetDefaultHeaders(len(body))
			h.Replace("Content-Type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody(body)
			return
		}
	}

	srv, err := server.Serve(port, problemHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer srv.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func respond200() []byte {
	return []byte(`
		<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
	`)
}

func respond400() []byte {
	return []byte(`
	<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
	`)
}

func respond500() []byte {
	return []byte(`
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
	`)
}
