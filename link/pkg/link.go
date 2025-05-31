package link

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

type Link struct {
	Text string
	Href string
}

func traverse(n *html.Node, links []Link) []Link {

	if n.Type == html.ElementNode && n.Data == "a" {
		fmt.Printf("found a link: %v\n", n.Attr)
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				fmt.Println("appending href", attr.Val)
				links = append(links, Link{Href: attr.Val, Text: "Test"})
			}
		}

	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = traverse(c, links)
	}

	return links
}

func HTMLFileToLinks(r io.Reader) ([]Link, error) {
	node, err := html.Parse(r)

	if err != nil {
		panic(err)
	}

	links := make([]Link, 0)

	links = traverse(node, links)

	return links, nil
}
