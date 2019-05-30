package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SitemapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles    []string `xml:"url>news>title"`
	Keywords  []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type NewsMap struct {
	keywords string
	Location string
}

type siteHolder struct {
	nm  map[string]NewsMap
	smi SitemapIndex
}

func (s siteHolder) homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(3)
	fmt.Fprintf(w, "<h1>all sites<h1/>")

	for _, Location := range s.smi.Locations {

		str := strings.Split(Location, "https://www.washingtonpost.com/news-sitemaps/")
		str = strings.Split(str[1], ".xml")
		fmt.Fprintf(w, `<p><a href="%s">%s </a></p>`, Location, str[0])
	}
	fmt.Println(4)
}

func main() {
	var s SitemapIndex
	var n News
	newsMap := make(map[string]NewsMap)

	var g siteHolder

	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml") //gets all

	bytes, _ := ioutil.ReadAll(resp.Body) //parses all into better format

	//resp.Body.Close() // closes since we got it all

	xml.Unmarshal(bytes, &s) // converts from binary to readable
	fmt.Println(1)
	for _, Location := range s.Locations {
		str := strings.Split(Location, "\n")
		//fmt.Printf("%s", str[1])
		resp, _ := http.Get(str[1])
		//gets all
		bytes, _ := ioutil.ReadAll(resp.Body) //parses all into better format
		//resp1.Body.Close()                      // closes since we got it all
		xml.Unmarshal(bytes, &n)
		for idx, _ := range n.Titles {
			newsMap[n.Titles[idx]] = NewsMap{n.Keywords[idx], n.Locations[idx]}
		}

	}
	fmt.Println(2)
	/*
		for idx, data := range newsMap {
			fmt.Println("\n\n\n", idx)
			fmt.Println("\n", data.keywords)
			fmt.Println("\n", data.Location)

		}*/

	g.nm = newsMap
	g.smi = s
	http.HandleFunc("/", g.homepage)
	http.ListenAndServe(":8000", nil)

}
