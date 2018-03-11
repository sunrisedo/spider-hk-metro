package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector()
	// c := colly.NewCollector(
	// 	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	// 	colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	// )

	// On every a element which has href attribute call callback
	c.OnHTML("div", func(e *colly.HTMLElement) {
		// link := e.Attr("href")
		// Print link
		value := e.ChildText(".subway-menu")
		log.Println("data:", value)
		fmt.Printf("Link found: %q\n", e.Text)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		// c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("http://map.baidu.com/?subwayShareId=hongkong,2912")
}
