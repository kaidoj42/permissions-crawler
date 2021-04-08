package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func crawl(URL string) {
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal("Could not connect to http server")
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".constants tr").Each(func(i int, s *goquery.Selection) {
		code := s.Find("td").Next().Find("code").Text()
		desc := s.Find("p").First().Text()
		fmt.Println(output(code, desc))
	})
}

func output(code string, desc string) string {
	return fmt.Sprintf(`,[
		'title' => 'android.permission.%s',
		'id' => 'android.permission.%s',
		'tags' => ['permission'],
		'description' => '%s'
		]`, code, code, clean(desc))
}

func clean(input string) string {
	return strings.TrimSpace(input)
}

func main() {
	crawl("https://developer.android.com/reference/android/Manifest.permission")
}
