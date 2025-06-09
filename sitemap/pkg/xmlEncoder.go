package sitemap

import (
	"encoding/xml"
	"log"
)

type UrlSet []Entry

type Entry struct {
	Url struct {
	} `xml:"url"`
}

type EntryLoc struct {
	Loc string `xml:"loc"`
}

func (us *UrlSet) EncodeXML() {
	xmlContents, err := xml.MarshalIndent(us, "", " ")
	if err != nil {
		panic(err)
	}

	log.Println(string(xmlContents))
}
