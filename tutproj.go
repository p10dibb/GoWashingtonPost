package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type SitemapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles    []string `xml:"url>news>title"`
	Keywords  []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type NewsMap struct {
	Keywords string
	Location string
}

// type siteHolder struct {
// 	nm  map[string]NewsMap
// 	smi SitemapIndex
// }
type newsAggPage struct {
	Title string
	Time  string
	News  map[string]NewsMap
}

type sitePage struct {
	Title string

	Topics []newsAggPage
}

// func (s siteHolder) homepage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(3)
// 	fmt.Fprintf(w, "<h1>all sites<h1/>")

// 	for _, Location := range s.smi.Locations {

// 		str := strings.Split(Location, "https://www.washingtonpost.com/news-sitemaps/")
// 		str = strings.Split(str[1], ".xml")
// 		fmt.Fprintf(w, `<p><a href="%s">%s </a></p>`, Location, str[0])
// 	}
// 	fmt.Println(4)
// }

func (p newsAggPage) newsAggHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(5)
	t, _ := template.ParseFiles("baseTemplate.html")
	t.Execute(w, p)
	fmt.Println(6)

}

func getMap(s string, nm map[string]NewsMap) {
	var n News
	//newsMap := make(map[string]NewsMap)
	resp, _ := http.Get(s)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &n)
	for idx, _ := range n.Titles {
		nm[n.Titles[idx]] = NewsMap{n.Keywords[idx], n.Locations[idx]}
	}

}

func main() {

	start := time.Now()

	fmt.Println(1)
	fmt.Println(2)
	var s SitemapIndex

	newsMap := make(map[string]NewsMap)

	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml") //gets all

	bytes, _ := ioutil.ReadAll(resp.Body) //parses all into better format

	xml.Unmarshal(bytes, &s) // converts from binary to readable
	fmt.Println(3)
	for _, Location := range s.Locations {

		str := strings.Split(Location, "\n") //has a /n at begiinig for some reason

		go getMap(str[1], newsMap)

	}

	fmt.Println(4)

	p := newsAggPage{Title: "News stuff", News: newsMap, Time: time.Now().String()}
	elapsed := time.Since(start)
	println("Binomial took %s", elapsed)

	http.HandleFunc("/", p.newsAggHandler)

	http.ListenAndServe(":8000", nil) //not sure what this does
	println(7)

}
