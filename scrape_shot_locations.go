package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/google/uuid"
)

type shotCoordinates struct{
	x_loc int
	y_loc int
}

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

func getShotCoordinates (shot_string string) shotCoordinates {
	// expect similar: "top: 204px; left: 80px"
	s := strings.Split(shot_string, ";") // ["top: 204px", "left: 80px"]
	x_coord_str := strings.Split(s[1], ": ")[1] //"80px"
	x_coord_str = strings.TrimSuffix(x_coord_str, "px") //"80"
	x_coord_int, err := strconv.Atoi(x_coord_str)
	if err != nil {
		log.Printf("error converting x coord to int [%v]: %v", x_coord_str, err)
		return shotCoordinates{}
	}

	y_coord_str := strings.Split(s[0], ": ")[1] //"204px"
	y_coord_str = strings.TrimSuffix(y_coord_str, "px") //"204"
	y_coord_int, err := strconv.Atoi(y_coord_str)
	if err != nil {
		log.Printf("error converting y coord to in [%v]: %v", y_coord_str, err)
		return shotCoordinates{}
	}
	return shotCoordinates{
		x_loc: x_coord_int,
		y_loc: y_coord_int,
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