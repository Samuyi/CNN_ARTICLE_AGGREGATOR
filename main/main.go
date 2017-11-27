package main

import (
	//"bufio"
	//"os"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
)

type SiteMapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

type ArticleIndex struct {
	Articles []string `xml:"url>loc"`
}

var ArticleMap map[string][]string
var wg sync.WaitGroup

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)

}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://edition.cnn.com/sitemaps/sitemap-index.xml")
	check(err)

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	check(err)

	ArticleMap = make(map[string][]string)

	var site SiteMapIndex
	xml.Unmarshal(bytes, &site)

	for _, location := range site.Locations {
		wg.Add(1)
		go goLink(location)

	}
	wg.Wait()
	t, err := template.ParseFiles("/home/samuyi/projects/Go_Networking/src/body.html")
	check(err)

	t.Execute(w, ArticleMap)

}

func goLink(location string) {
	resp, err := http.Get(location)
	check(err)

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	check(err)

	var articleLocation ArticleIndex
	xml.Unmarshal(bytes, &articleLocation)
	key := location[36:]
	ArticleMap[key] = articleLocation.Articles
	defer wg.Done()
}
