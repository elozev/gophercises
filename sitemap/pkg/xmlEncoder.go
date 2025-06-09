package sitemap

import (
	"encoding/xml"
	"fmt"
	"os"
)

type UrlSet struct {
	XMLName xml.Name   `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	Urls    []EntryLoc `xml:"url"`
}

type EntryLoc struct {
	Loc string `xml:"loc"`
}

func (us *UrlSet) EncodeXML(host string) {
	out, err := xml.MarshalIndent(us, "", " ")
	if err != nil {
		panic(err)
	}

	xmlContents := append([]byte(xml.Header), out...)

	sitemapFile, err := os.Create(fmt.Sprintf("%s.xml", host))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file for host %s\n", host)
		panic(err)
	}

	_, err = sitemapFile.Write(xmlContents)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing bytes to file; err %+v\n", err)
		panic(err)
	}
}
