package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

// func main() {
// 	c := colly.NewCollector()

// 	// Find and visit all links
// 	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
// 		link := e.Attr("href")

// 		fmt.Printf("ChildText: %s\n", e.ChildText("RYGLineStatus"))
// 		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
// 		e.Request.Visit(link)
// 	})

// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting", r.URL)
// 	})

// 	// c.Visit("http://go-colly.org/")
// 	c.Visit("http://www.mtr.com.hk/ch/customer/main/index.html")

// }

func main() {
	c := colly.NewCollector()

	// Find and visit all links
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// c.OnHTML(`<table class="table-d table-topline">`, func(e *colly.HTMLElement) {
	c.OnHTML("tbody tr td", func(e *colly.HTMLElement) {
		// start := e.ChildText(".firstTrain")
		// end := e.ChildText("td")
		fmt.Printf("%s\n", e.Text)
		// fmt.Printf("ChildText: name:%s\nText:%s\nobj:%v\n", e.Name, e.Text, *e)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// c.Visit("http://go-colly.org/")
	c.Visit("http://www.mtr.com.hk/ch/customer/services/service_hours_search.php?query_type=search&station=3")

}
