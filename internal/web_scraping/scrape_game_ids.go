package web_scraping

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeGameLinks (url string) []GameLink {

	// file, err := os.Open("example_pages/games_of_day.htm")
	// if err != nil {
	// 	log.Fatalf("Error opening file: %v", err)
	// }
	// defer file.Close()

	// Load the HTML document
	// doc, err := goquery.NewDocumentFromReader(file)
	// if err != nil {
	// 	log.Fatalf("Error parsing HTML: %v", err)
	// }


	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return []GameLink{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		return []GameLink{}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("error getting document from reader: %v\n", err)
	}

	var gameLinks []GameLink
	game_summary_div := doc.Find(".game_summaries").Last()

	game_summary_div.Find("td.right.gamelink a").Each(func(i int, link_tag *goquery.Selection) {
		link, exist := link_tag.Attr("href")
		if exist { 
			var url string
			if strings.HasPrefix(link, "https://www.hockey-reference.com"){
				url = link
			} else {
				url = "https://www.hockey-reference.com" + link
			}
			gameLinks = append(gameLinks, GameLink{
				Url: url,
				Gameid: strings.TrimSuffix(path.Base(link), ".html"),
			})
		}
	})

	return gameLinks
}