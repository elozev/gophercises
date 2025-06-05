package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	link "github.com/elozev/gophercises/link/pkg"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	var filename = flag.String("filename", "ex1.html", "HTML file to parse!")
	flag.Parse()

	if filepath.Ext(*filename) != ".html" {
		log.Fatal("Only .html files are supported")
	}

	file, err := os.Open(*filename)
	check(err)
	defer file.Close()

	links, err := link.HTMLFileToLinks(file)
	check(err)

	fmt.Printf("links: %+v \n", links)
}
