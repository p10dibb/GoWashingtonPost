package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SitemapIndex struct {
	Locations []Location `xml:"sitemap"`
}
type Location struct {
	Loc string `xml:"loc"`
}

func (l Location) String() string {
	return fmt.Sprintf(l.Loc)
}

func (s SitemapIndex) homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>all sites<h1/>")

	for _, Location := range s.Locations {
		str := strings.Split(Location.String(), "https://www.washingtonpost.com/news-sitemaps/")
		str = strings.Split(str[1], ".xml")
		fmt.Fprintf(w, `<p><a href="%s">%s </a></p>`, Location, str[0])
	}

}

func main() {
	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml") //gets all

	bytes, _ := ioutil.ReadAll(resp.Body) //parses all into better format

	resp.Body.Close() // closes since we got it all

	var s SitemapIndex

	xml.Unmarshal(bytes, &s)
	//fmt.Println(s.Locations)

	/*for _, Location := range s.Locations {
		fmt.Printf("\n%s", Location)
	}*/

	http.HandleFunc("/", s.homepage)
	http.ListenAndServe(":8001", nil)
}
