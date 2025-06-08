package sitemap

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	link "github.com/elozev/gophercises/link/pkg"
)

func check(err error, critical bool) {
	if err != nil {
		log.Printf("error: %v", err)
		if critical {
			os.Exit(-1)
		}
	}
}

func getPage(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("failed to retrieve %s! err: %v", url, err)
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	return body, err
}

func cleanLinkPrefix(url string, host string) string {
	matcher := fmt.Sprintf("(^\\/)|(^http(s)?:\\/\\/%s\\/)", host)
	rgx := regexp.MustCompile(matcher)

	return "/" + rgx.ReplaceAllString(url, "")
}

func isLinkInternal(l link.Link, host string) bool {
	matcher := fmt.Sprintf("(^\\/)|(^http(s)?:\\/\\/%s)", host)
	rgx := regexp.MustCompile(matcher)

	return rgx.MatchString(l.Href)
}

func onlyInternal(links []link.Link, baseURL string) (filtered []link.Link) {
	pURL, err := url.Parse(baseURL)
	check(err, true)
	host := pURL.Host

	filtered = make([]link.Link, 0, len(links))

	for _, l := range links {
		if isLinkInternal(l, host) {
			l.Href = cleanLinkPrefix(l.Href, host)
			filtered = append(filtered, l)
		}
	}

	return
}

func parseLinks(body io.Reader) (links []link.Link) {
	links, err := link.HTMLFileToLinks(body)
	if err != nil {
		panic(err)
	}
	return
}

func getPageInternalLinks(url string) []link.Link {
	body, err := getPage(url)
	check(err, true)
	br := bytes.NewReader(body)
	links := parseLinks(br)
	internal := onlyInternal(links, url)
	for _, i := range internal {
		log.Println(i.Href)
	}
	return internal
}

func traverse(baseURL string, path string, visited map[string]bool) map[string]bool {
	url, err := url.Parse(baseURL + path)
	check(err, true)

	log.Printf("Traversing %s\n", url.String())

	links := getPageInternalLinks(url.String())

	visited[path] = true

	for _, l := range links {
		if _, ok := visited[l.Href]; !ok {
			traverse(baseURL, l.Href, visited)
		}
	}

	return visited
}

func cleanUrl(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	check(err, true)
	return parsedURL.Scheme + "://" + parsedURL.Host
}

func RetrieveSiteMap(baseURL string) {
	cURL := cleanUrl(baseURL)

	visited := make(map[string]bool)

	v := traverse(cURL, "/", visited)

	fmt.Printf("visited: %+v\ntotal links: %d\n", v, len(v))

}
