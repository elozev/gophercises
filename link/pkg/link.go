package link

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Text string
	Href string
}

func extractTextHelper(n *html.Node, text string) string {
	if n.Type == html.TextNode {
		if len(text) > 0 && !strings.HasSuffix(text, " ") {
			text = text + " "
		}

		newText := strings.TrimSpace(n.Data)
		matched, err := regexp.Match("^(,|!|;|\\.).*", []byte(newText))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to match new data %s; err: %v\n", n.Data, err)
		}

		if matched {
			text = strings.TrimRight(text, " ")
		}

		sanitised := strings.ReplaceAll(newText, "\n", "")
		sanitised = strings.ReplaceAll(sanitised, "\t", " ")
		sanitised = strings.ReplaceAll(sanitised, "  ", " ")

		text = text + strings.TrimSpace(sanitised)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text = extractTextHelper(c, text)
	}

	return strings.TrimSpace(text)
}

func extractText(n *html.Node) string {
	return extractTextHelper(n, "")
}

func traverse(n *html.Node, links []Link) []Link {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				text := extractText(n)
				links = append(links, Link{Href: attr.Val, Text: text})
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
