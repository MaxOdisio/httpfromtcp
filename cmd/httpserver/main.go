package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/maxodisio/httpfromtcp/internal/headers"
	"github.com/maxodisio/httpfromtcp/internal/request"
	"github.com/maxodisio/httpfromtcp/internal/response"
	"github.com/maxodisio/httpfromtcp/internal/server"
)

const port = 42069

func main() {
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

func problemHandler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget

	if strings.HasPrefix(target, "/httpbin/") || target == "/httpbin" {
		handleHttpbinProxy(w, target)
		return
	}

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
	case "/video":
		handleVideo(w, req)
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

func handleHttpbinProxy(w *response.Writer, targetPath string) {
	trimmedPath := strings.TrimPrefix(targetPath, "/httpbin")
	targetURL := "https://httpbingo.org" + trimmedPath

	resp, err := http.Get(targetURL)
	if err != nil {
		w.WriteStatusLine(response.StatusError)
		body := respond500()
		h := response.GetDefaultHeaders(len(body))
		w.WriteHeaders(h)
		w.WriteBody(body)
		return
	}
	defer resp.Body.Close()

	resHeaders := headers.NewHeaders()
	for k, val := range resp.Header {
		if strings.EqualFold(k, "Content-Length") {
			continue
		}
		if len(val) > 0 {
			resHeaders.Set(k, strings.Join(val, ", "))
		}
	}

	resHeaders.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	resHeaders.Set("Transfer-Encoding", "chunked")
	resHeaders.Set("Connection", "close")

	w.WriteStatusLine(response.StatusCode(resp.StatusCode))
	w.WriteHeaders(resHeaders)

	var totalBytes int
	var fullBody []byte
	buf := make([]byte, 32)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			chunk := buf[:n]
			totalBytes += n
			fullBody = append(fullBody, chunk...)

			fmt.Printf("Read %d bytes from httpbingo.org\n", n)
			_, wErr := w.WriteChunkedBody(chunk)
			if wErr != nil {
				log.Printf("Error writing chunk to client: %v", wErr)
				return
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading response body from \"httpbingo.org\": %v", err)
			return
		}
	}

	hsh := sha256.Sum256(fullBody)
	stringHsh := hex.EncodeToString(hsh[:])

	trailers := headers.NewHeaders()
	trailers.Set("X-Content-SHA256", stringHsh)
	trailers.Set("X-Content-Length", strconv.Itoa(len(fullBody)))

	if err = w.WriteTrailers(trailers); err != nil {
		log.Printf("Error writing trailers: %v", err)
	}

	log.Printf("Proxy stream completed successfully. Total bytes transferred: %d", totalBytes)
}

func handleVideo(w *response.Writer, req *request.Request) {
	method := req.RequestLine.Method

	if method != "GET" {
		w.WriteStatusLine(response.StatusBadRequest)
		body := respond400()
		h := response.GetDefaultHeaders(len(body))
		h.Replace("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody(body)
		return
	}

	data, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		w.WriteStatusLine(response.StatusError)
		body := respond500()
		h := response.GetDefaultHeaders(len(body))
		h.Replace("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody(body)
		return
	}

	w.WriteStatusLine(response.StatusOK)
	body := data
	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "video/mp4")
	w.WriteHeaders(h)
	w.WriteBody(body)
}
