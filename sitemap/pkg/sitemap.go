package sitemap

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func retrievePage(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("failed to retrieve %s! err: %v", url, err)
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	return body, err
}

func parseLinks(body []byte) {
	links := link.HTMLFileToLinks(body)

	fmt.Println(links)
}

func RetrieveSiteMap(baseurl string) {
	body, _ := retrievePage(baseurl)
	// fmt.Println(string(body))
	parseLinks(body)
}
