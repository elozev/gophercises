package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	link "link/pkg"

	"golang.org/x/net/html"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func traverse(node *html.Node) {
	if node.FirstChild != nil {
		fmt.Printf("type: <%v> | first child: <%v> \n", node.DataAtom, node.FirstChild.DataAtom)
	} else {
		fmt.Printf("type: <%v> \n", node.DataAtom)
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverse(c)
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

	// contents, err := io.ReadAll(file)
	// check(err)

	// log.Println("read from file: \n" + string(contents))
	// fmt.Println("------------------------")

	// node, err := link.HTMLFileToLinks(file)

	links, err := link.HTMLFileToLinks(file)
	check(err)

	fmt.Printf("list of links: %v \n", links)

	// traverse(node)
}
