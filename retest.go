package main

import (
	"log"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	// Find and visit all links
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// c.OnHTML(`<table class="table-d table-topline">`, func(e *colly.HTMLElement) {

	// c.OnHTML("#sbw_map", func(e *colly.HTMLElement) {
	// 	// start := e.ChildText(".firstTrain")
	// 	value := e.ChildText(".subway-menu")
	// 	log.Println("data:", value)
	// 	// end := e.ChildText("td")
	// 	fmt.Printf("%v\n", *e)
	// 	// fmt.Printf("ChildText: name:%s\nText:%s\nobj:%v\n", e.Name, e.Text, *e)

	// })
	// c.onRep
	detailCollector := c.Clone()
	c.OnHTML("#sbw_map", func(e *colly.HTMLElement) {
		// If attribute class is this long string return from callback
		// As this a is irrelevant
		log.Println("data：", *e)
		log.Println("data2：", e.ChildText(".subway-menu"))
		if e.Attr("class") == "subway-menu" {
			return
		}
		link := e.Attr("href")
		// If link start with browse or includes either signup or login return from callback
		log.Println("link：", link)
		// start scaping the page under the link found
		e.Request.Visit(link)
	})
	// nHTML("#currencies-all tbody tr", func(e *colly.HTMLElement) {
	detailCollector.OnHTML(".subway-menu", func(e *colly.HTMLElement) {
		// If attribute class is this long string return from callback
		// As this a is irrelevant
		log.Println("data：", *e)
		log.Println("data2：", e.ChildText(".subway-menu"))
		if e.Attr("class") == "subway-menu" {
			return
		}
		link := e.Attr("href")
		// If link start with browse or includes either signup or login return from callback
		log.Println("link：", link)
		// start scaping the page under the link found
		e.Request.Visit(link)
	})
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })
	// c.Visit("http://www.mtr.com.hk/ch/customer/services/service_hours_search.php?query_type=search&station=3")
	c.Visit("http://map.baidu.com/?subwayShareId=hongkong,2912/")
}
