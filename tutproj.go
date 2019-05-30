package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SitemapIndex struct {
	Locations []string `xml:"sitemap,loc"`
}
type News struct {
	Titles    []string `xml:"url>news>title"`
	Keywords  []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

func (s SitemapIndex) homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>all sites<h1/>")

	for _, Location := range s.Locations {
		str := strings.Split(Location, "https://www.washingtonpost.com/news-sitemaps/")
		str = strings.Split(str[1], ".xml")
		fmt.Fprintf(w, `<p><a href="%s">%s </a></p>`, Location, str[0])
	}

}

func main() {
	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml") //gets all

	bytes, _ := ioutil.ReadAll(resp.Body) //parses all into better format

	resp.Body.Close() // closes since we got it all

	var s SitemapIndex

	xml.Unmarshal(bytes, &s) // converts from binary to readable

	http.HandleFunc("/", s.homepage)
	http.ListenAndServe(":8001", nil)
}
