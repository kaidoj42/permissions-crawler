package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func crawl(URL string) string {
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal("Could not connect to http server")
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var buffer bytes.Buffer

	doc.Find(".constants tr").Each(func(i int, s *goquery.Selection) {
		code := s.Find("td").Next().Find("code").First().Text()
		desc := s.Find("p").First().Text()
		buffer.WriteString(output(code, desc))
	})

	return buffer.String()
}

func output(code string, desc string) string {
	if code == "" {
		return ""
	}
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

func write(input string) {
	f, err := os.Create("results.txt")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(input); err != nil {
		log.Println(err)
	}
}

func main() {
	results := crawl("https://developer.android.com/reference/android/Manifest.permission")
	fmt.Println("Done crawling")
	if results != "" {
		write(results)
	}
	fmt.Printf("Done writing %d characters \r\n", len(results))
}
