package sitemap

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var aSingleLink = `
<a href="/test123">Test123</a>
`

func createNode(content string) (n *html.Node) {
	n, err := html.Parse(strings.NewReader(content))
	if err != nil {
		panic(err)
	}
	return
}

func createAnchor(href, text string) *html.Node {
	attrs := make([]html.Attribute, 1)
	attrs[0] = html.Attribute{
		Key:       "href",
		Val:       href,
		Namespace: "",
	}

	textNode := createNode(text)

	n := &html.Node{
		Parent:      nil,
		FirstChild:  textNode,
		LastChild:   textNode,
		PrevSibling: nil,
		NextSibling: nil,
		Attr:        attrs,
		Namespace:   "",
		Data:        "a",
	}

	return n
}

const format = "expected: %+v; got: %+v"

// Write test for extractText
func TestExtractText(t *testing.T) {
	expected := "This is the expected text!"

	nodeContent := fmt.Sprintf("<h1>%s</h1>", expected)
	node := createNode(nodeContent)
	got := extractText(node)

	if expected != got {
		t.Errorf(format, expected, got)
	}
}

func TestExtractTextWithComment(t *testing.T) {
	expected := "This is some title!"
	nodeContent := fmt.Sprintf("<h2>%s<!--Some comment--></h2>", expected)
	node := createNode(nodeContent)
	got := extractText(node)

	if expected != got {
		t.Errorf(format, expected, got)
	}
}

const complexNodeContent = `
<div>
	<p >
		This text is very
		<span>important</span>, and
		this one isn't.
	</p>
</div>
`

func TestNestedNode(t *testing.T) {
	expected := "This text is very important, and this one isn't."
	complexNode := createNode(complexNodeContent)

	got := extractText(complexNode)

	if expected != got {
		t.Errorf(format, expected, got)
	}
}

// Write tests for traverse

// Write tests for HTMLFileToLinks

func equalLinks(l1 Link, l2 Link) bool {
	return l1.Href == l2.Href && l1.Text == l2.Text
}

func equalLinkSlices(ls1 []Link, ls2 []Link) bool {
	if len(ls1) != len(ls2) {
		return false
	}

	for i := 0; i < len(ls1); i++ {
		if !equalLinks(ls1[i], ls2[i]) {
			return false
		}
	}
	return true
}

func TestTraverse(t *testing.T) {
	r := strings.NewReader(aSingleLink)

	got, err := HTMLFileToLinks(r)

	expected := [1]Link{{Href: "/test123", Text: "Test123"}}

	if equalLinkSlices(expected[:], got) || err != nil {
		t.Errorf("expected %+v; got %+v", expected, got)
	}

}
