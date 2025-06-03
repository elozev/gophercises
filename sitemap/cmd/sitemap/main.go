package main

import (
	"flag"
	sitemap "sitemap/pkg"
)

func main() {
	var url = flag.String("url", "", "website to build a sitemap for")
	flag.Parse()

	if *url == "" {
		panic("-url is required")
	}

	sitemap.RetrieveSiteMap(*url)
}
