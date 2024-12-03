package web_scraping

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/ehmker/hockey_stats/internal/shared"
	"github.com/google/uuid"
)


func AddShotLocationsToDB (s shared.State, doc *goquery.Document, gameID string) {
	shotList := scrapeShotLocations(s, doc, gameID)

	for _, shot := range shotList {
		_, err := s.DB.CreateShot(context.Background(), shot)
		if err != nil {
			log.Panicf("error adding scoring summary: %v\n", err)
		}
	}
}



func scrapeShotLocations (s shared.State, doc *goquery.Document, ID string) []database.CreateShotParams{
	// doc, err := goquery.NewDocumentFromReader(f)
	// if err != nil {
	// 	log.Fatalf("Error parsing HTML: %v", err)
	// }
	var shotList []database.CreateShotParams

	doc.Find(".shotchart").Each(func(i int, div *goquery.Selection) {
		full_team_name := div.Find("h4").Text()
		short_name, err := s.DB.GetTeamShortName(context.Background(), full_team_name)
		if err != nil {
			log.Printf("unable to get team short name of [%v]: %v\n", full_team_name, err)
			return
		}
		div.Find("div").Children().Each(func(i int, shot *goquery.Selection) {
			shotParams := database.CreateShotParams{
				ID: uuid.New(),
				Gameid: ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Team: short_name,
				Goal: false,
			}
			//Extract if shot was a goal
			classification, _ := shot.Attr("class") // expect: "Goal" or "Shot"
			if strings.ToLower(classification) == "goal"{
				shotParams.Goal = true
			}

			//Extract the coordinates of where the shot was taken
			shot_location_string, _ := shot.Attr("style") // expect similar: "top: 204px; left: 80px"
			shot_coords := getShotCoordinates(shot_location_string)
			shotParams.XLoc = int32(shot_coords.x_loc)
			shotParams.YLoc = int32(shot_coords.y_loc)

			//Extract the player who took the shot
			player_string, _ := shot.Attr("title") // expect similar: "Saved Shot - Tom Wilson" 
			shotParams.Player = strings.Split(player_string, " - ")[1]

			shotList = append(shotList, shotParams)
		})
	})
	return shotList
}