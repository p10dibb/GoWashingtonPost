package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
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
	Keywords string
	Location string
}

// type siteHolder struct {
// 	nm  map[string]NewsMap
// 	smi SitemapIndex
// }
type newsAggPage struct {
	Title string
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

func newsAggHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(2)
	var s SitemapIndex
	var n News

	var sp sitePage
	sp.Title = "washington post"

	//var sp sitePage

	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml") //gets all

	bytes, _ := ioutil.ReadAll(resp.Body) //parses all into better format

	xml.Unmarshal(bytes, &s) // converts from binary to readable
	fmt.Println(3)
	for _, Location := range s.Locations {
		//println(Location)
		newsMap := make(map[string]NewsMap)

		str := strings.Split(Location, "\n") //has a /n at begiinig for some reason

		resp, _ := http.Get(str[1])
		str = strings.Split(str[1], "https://www.washingtonpost.com/news-sitemaps/")
		str = strings.Split(str[1], ".xml")

		bytes, _ := ioutil.ReadAll(resp.Body) //parses all into better format
		//resp1.Body.Close()                      // closes since we got it all
		xml.Unmarshal(bytes, &n)
		for idx, _ := range n.Titles {
			newsMap[n.Titles[idx]] = NewsMap{n.Keywords[idx], n.Locations[idx]}
		}

		// for key, value := range newsMap {
		// 	println(key)
		// 	value = value
		// }

		str = strings.Split(str[1], "https://www.washingtonpost.com/news-sitemaps/")
		str = strings.Split(str[1], ".xml")

		var p newsAggPage
		p.Title = str[0]

		for key, value := range newsMap {
			p.News[key] = value
			delete(newsMap, key)
		}
		sp.Topics = append(sp.Topics, p)
		// for key, val := range newsMap {
		// 	val = val
		// 	delete(newsMap, key)
		// }
	}
	fmt.Println(5)

	//p := newsAggPage{Title: "News stuff", News: newsMap}
	t, _ := template.ParseFiles("upTemp.html")
	t.Execute(w, sp)
	fmt.Println(6)

}

func main() {

	fmt.Println(1)
	http.HandleFunc("/", newsAggHandler)
	http.ListenAndServe(":8000", nil)
	// http.HandleFunc("/", g.homepage)

}
