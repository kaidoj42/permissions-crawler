package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fastjson"
)

func crawlAndroid() string {
	res, err := http.Get("https://developer.android.com/reference/android/Manifest.permission")
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
		buffer.WriteString(output("android.permission."+code, desc))
	})

	return buffer.String()
}

func crawlIOS() string {
	//https://developer.apple.com/documentation/bundleresources/information_property_list/protected_resources
	res, err := http.Get("https://developer.apple.com/tutorials/data/documentation/bundleresources/information_property_list/protected_resources.json")

	if err != nil {
		log.Fatal("Could not connect to http server")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Could not read request body")
	}
	var p fastjson.Parser
	v, err := p.Parse(string(body))
	if err != nil {
		log.Fatalf("cannot parse json: %s", err)
	}

	var buffer bytes.Buffer
	// Visit all the items in the top object
	v.GetObject("references").Visit(func(k []byte, v *fastjson.Value) {
		role := string(v.GetStringBytes("role"))
		code := string(v.GetStringBytes("title"))
		desc := string(v.GetStringBytes("abstract", "0", "text"))
		if role == "symbol" {
			buffer.WriteString(output(code, desc))
		}

	})
	return buffer.String()
}

func output(code string, desc string) string {
	if code == "" {
		return ""
	}
	return fmt.Sprintf(`,[
		'title' => '%s',
		'id' => '%s',
		'tags' => ['permission'],
		'description' => '%s'
		]`, code, code, clean(desc))
}

func clean(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "'", "\\'")
	return input
}

func write(filename string, input string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(input); err != nil {
		log.Println(err)
	}
}

func main() {
	androidResults := crawlAndroid()
	log.Println("[Android] Done crawling permissions")
	if androidResults != "" {
		write("android_results.txt", androidResults)
	}
	log.Printf("[Android] Done writing %d characters to file \r\n", len(androidResults))

	iosResults := crawlIOS()
	log.Println("[iOS] Done crawling permissions")
	if iosResults != "" {
		write("ios_results.txt", iosResults)
	}
	log.Printf("[iOS] Done writing %d characters to file \r\n", len(iosResults))
}
