package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/google/uuid"
)


func AddShotLocationsToDB (s state) {
	// Open the local HTML file
	file, err := os.Open("example_2.htm")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	shotList := scrapeShotLocations(file, s)

	for _, shot := range shotList {
		_, err = s.db.CreateShot(context.Background(), shot)
		if err != nil {
			log.Panicf("error adding scoring summary: %v\n", err)
		}
	}
}



func scrapeShotLocations (f *os.File, s state) []database.CreateShotParams{
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}
	var shotList []database.CreateShotParams

	doc.Find(".shotchart").Each(func(i int, div *goquery.Selection) {
		full_team_name := div.Find("h4").Text()
		short_name, err := s.db.GetTeamShortName(context.Background(), full_team_name)
		if err != nil {
			log.Printf("unable to get team short name of [%v]: %v\n", full_team_name, err)
			return
		}
		div.Find("div").Children().Each(func(i int, shot *goquery.Selection) {
			shotParams := database.CreateShotParams{
				ID: uuid.New(),
				Gameid: "1",
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