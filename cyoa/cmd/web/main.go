package main

import (
	cyoa "cyoa/src"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	filename = flag.String("filename", "gopher.json", "Filename to be used as a base for the story program!")
	port     = flag.Int("port", 8080, "Default port to run the program on!")
	cli      = flag.Bool("cli", false, "Run the program as CLI program!")
)

func main() {
	flag.Parse()

	var sh cyoa.StoryHandler

	// Open the gopher.json file
	file, err := os.Open(*filename)
	cyoa.Check(err)
	defer file.Close()

	stories, err := cyoa.JsonStory(file)
	cyoa.Check(err)
	sh.Stories = stories

	if !*cli {
		mux := http.NewServeMux()

		mux.HandleFunc("/", sh.IndexHandler)
		mux.HandleFunc("/story/{chapter}", sh.StoryHandler)

		log.Printf("Listening and serving on port %d\n", *port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
	} else {
		sh.CliHandler()
	}
}
