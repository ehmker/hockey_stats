package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
)

func ScrapeGameResults () (database.CreateGameResultParams, error) {
	// Open the local HTML file
	file, err := os.Open("boxscore_example.htm")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	results := database.CreateGameResultParams {
		ID: "1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	

	// Find team names and scores
	doc.Find(".scorebox > div").Each(func(i int, s *goquery.Selection) {
		// Early return incase more than the expected items are found
		if i > 1 { 
			return
		}

		// Extract team name
		teamName := s.Find("strong a").Text()

		// Extract score
		score_string := s.Find(".scores .score").Text()
		var score_int int
		if score_string != "" {
			score_int, err = strconv.Atoi(score_string)
			if err != nil{
				log.Fatalf("error parsing int: %v", err)
			}
		}
		
		// Extract team record
		var record string
		s.Find("div").Each(func(i int, recordDiv *goquery.Selection) {
			text := strings.TrimSpace(recordDiv.Text())
			//Check if found text matches expected W-L-OT format
			if strings.Count(text, "-") == 2{
				record = text
			}
		})
		// Away Team
		if i == 0 { 
			results.AwayTeam = teamName
			results.AwayTeamScore = int32(score_int)
			results.AwayTeamRecord = record
		}
		// Home Team
		if i == 1 { 
			results.HomeTeam = teamName
			results.HomeTeamScore = int32(score_int)
			results.HomeTeamRecord = record
		}

		
	})

	// Find additional details in "scorebox_meta"
	// 0: Game Date Played
	// 1: Game Attendance
	// 2: Game Arena 
	// 3: Game Duration
	doc.Find(".scorebox_meta > div").Each(func(i int, s *goquery.Selection) {
		switch i {
		// Extract Date Played
		case 0:
			timeLayout := "January 02, 2006, 3:04 PM"
			parsedTime, err := time.Parse(timeLayout, s.Text())
			fmt.Println(parsedTime)
			if err != nil {
				fmt.Printf("unable to parse time for Date Played: %v\n", err)
				return
			}
			results.DatePlayed = parsedTime
		
		//Extract attendance
		//Expected format "Attendance: xx,xxx"
		case 1:
			attendance_string := strings.Split(s.Text(), ": ")[1]
			attendance_string = strings.Replace(attendance_string, ",", "", -1)
			attendance, err := strconv.Atoi(attendance_string)
			if err != nil {
				fmt.Printf("unable to parse int for Attendance: %v\n", err)
				return
			}
			results.Attendance = int32(attendance)
		
		//Extract arena name
		case 2:
			results.Arena = s.Text()
		
		//Extract game duration
		//Expected format "Game Duration: 2:28"
		case 3: 
			duration_string := strings.Split(s.Text(), ": ")[1]
			fmt.Println(duration_string)
			timeLayout := "3:04"
			parsedTime, err := time.Parse(timeLayout, duration_string)
			if err != nil {
				fmt.Printf("unable to parse time for duration: %v\n", err)
				return
			}
			results.Duration = parsedTime
		}
		
	})

	// Determine W/L/OTL result

	//Check if overtime was needed to determine winner
	var OTNeeded bool
	if doc.Find(".game_summary.nohover.current .teams .right").Last().Text() == "OT" {
		OTNeeded = true
	} else {
		OTNeeded = false
	}

	//
	if results.HomeTeamScore > results.AwayTeamScore {
		results.HomeTeamResult = "W"
		if OTNeeded {
			results.AwayTeamResult = "OTL"
		} else {
			results.AwayTeamResult = "L"
		}
	} else {
		results.AwayTeamResult = "W"
		if OTNeeded {
			results.HomeTeamResult = "OTL"
		} else {
			results.HomeTeamResult = "L"
		}
	}


	return results, nil
}