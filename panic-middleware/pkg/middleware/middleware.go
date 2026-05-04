package middleware

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
)

var environment = os.Getenv("ENV")

func PanicMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				message := "Something went wrong!"
				if environment == "development" {
					stack := debug.Stack()
					hStack := addLinksToStack(stack)
					message = fmt.Sprintf("<h1>Panic occurred: %v</h1><pre>%s</pre>", err, string(hStack))
				}
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, message)
			}
		}()

		// new ResponseWriter of our own
		nw := &ResponseWriter{ResponseWriter: w}
		next.ServeHTTP(nw, r)
		nw.flush()
	}
}

func addLinksToStack(stack []byte) string {
	exp := regexp.MustCompile(`(\/(.*)):([0-9]+)`)
	filepathExp := regexp.MustCompile(`:([0-9]+)`)
	lineNumberExp := regexp.MustCompile(`(\/(.*)):`)

	res := exp.ReplaceAllStringFunc(string(stack), func(s string) string {

		filepath := filepathExp.Split(s, 2)
		lineNumber := lineNumberExp.Split(s, 2)

		return fmt.Sprintf(`<a href="/render?filename=%s&line=%s" target="_blank">%s</a>`, filepath[0], lineNumber[1], s)
	})

	return res
}

type ResponseWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	cp := make([]byte, len(b))
	copy(cp, b)
	rw.writes = append(rw.writes, cp)
	return len(b), nil
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *ResponseWriter) flush() error {
	totalLen := 0
	for _, write := range rw.writes {
		totalLen += len(write)
	}
	rw.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(totalLen))

	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}

	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter does not support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *ResponseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}
	flusher.Flush()
}
