package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"panic-middleware/pkg/middleware"
	"strconv"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleGet)
	mux.HandleFunc("GET /panic", handlePanic)
	mux.HandleFunc("GET /render", handleRender)

	log.Println("Running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.PanicMiddleware(mux)))
}

func handleRender(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filename := query.Get("filename")
	if filename == "" {
		http.Error(w, "filename is a required query param", http.StatusBadRequest)
		return
	}

	lineNumberString := query.Get("line")
	line := 0
	if lineNumberString != "" {
		line, _ = strconv.Atoi(lineNumberString)
	}

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed opening file: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")

	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}
	htmlOptions := []html.Option{html.WithClasses(false), html.WithLineNumbers(true)}

	if line != 0 {
		htmlOptions = append(htmlOptions, html.HighlightLines([][2]int{{line, line}}))
	}

	formatter := html.New(htmlOptions...)

	iterator, _ := lexer.Tokenise(nil, string(fileContents))

	formatter.Format(w, style, iterator)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func handlePanic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("x-api-key", "123456")
	w.Write([]byte("successful response"))
	w.WriteHeader(http.StatusOK)
	panic("ooops, something went wrong")
}
