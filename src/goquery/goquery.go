package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

/*
	<article class="clearfix">
		<div class="thumb">
			<a href="https://www.metalsucks.net/2019/11/18/album-review-blood-incantation-reveal-the-hidden-history-of-the-human-race/"><img src="https://www.metalsucks.net/wp-content/uploads/2019/10/bloodincantation-hiddenhistory-150x150.jpg"></a>
		</div>
		<div class="content-block">
			<a class="header-xs" href="https://www.metalsucks.net/2019/11/18/album-review-blood-incantation-reveal-the-hidden-history-of-the-human-race/">Blood Incantation</a>
			<i>Hidden History of the Human Race</i>
			<div class="rating">
				<span>Rating</span>
				<img src="https://www.metalsucks.net/wp-content/themes/metalsucks.v5/images/ratings/rating-50.svg">
			</div>
		</div>
	</article>
*/
func ExampleScrape() {
	doc, err := goquery.NewDocument("http://metalsucks.net")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}

func main() {
	ExampleScrape()
}
