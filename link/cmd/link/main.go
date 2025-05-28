package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/net/html"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func traverse(node *html.Node) {
	log.Printf("node: %+v\n", node)
	for desc := range node.ChildNodes() {
		traverse(desc)
	}

	return
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

	contents, err := io.ReadAll(file)
	check(err)

	log.Println("read from file: \n" + string(contents))

	node, err := html.Parse(file)
	check(err)

	traverse(node)
	// fmt.Printf("First node is: %v \n", node.Namespace)
	fmt.Printf("First node is: %+v \n", node)
	// fmt.Printf("First node is: %#v \n", html.NodeType(node.Type))
}
