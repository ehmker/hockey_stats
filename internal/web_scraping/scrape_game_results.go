package web_scraping

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/ehmker/hockey_stats/internal/shared"
)

// Gets the game details and adds to DB, returns the season value for use in other Add functions
func addGameResults(s shared.State, doc *goquery.Document, gameID string) (database.CreateGameResultParams, error) {
	gameResult, err := ScrapeGameResults(doc, gameID)
	if err != nil {
		return database.CreateGameResultParams{}, err
	}
	_, err = s.DB.CreateGameResult(context.Background(), gameResult)
	if err != nil {
		return database.CreateGameResultParams{}, err
	}
	return gameResult, nil
}

func ScrapeGameResults (doc *goquery.Document, gameID string) (database.CreateGameResultParams, error) {
	results := database.CreateGameResultParams {
		ID: gameID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
		var err error
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
			timeLayout := "January 2, 2006, 3:04 PM"
			parsedTime, err := time.Parse(timeLayout, s.Text())
			// fmt.Println(parsedTime)
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
			// fmt.Println(duration_string)
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
	lastRow := doc.Find(".game_summary.nohover.current tbody tr").Last()

	// Step 2: Select the last td in that row
	str := lastRow.Find("td").Last().Text()
	str = strings.TrimSpace(str)
	// fmt.Println(str)


	//fmt.Println(str)
	if  str == "OT" || str == "SO"  {
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
	results.Season = setSeason(results.DatePlayed)

	return results, nil
}

//Sets the season based on the date the game was played.  
//If played before Aug 1st, the season is set to the same as the year the game was played
//otherwise the following year. 
func setSeason(d time.Time) int32 {
	new_season_start := time.Date(d.Year(), time.August, 1, 0, 0, 0, 0, time.UTC)
	if d.Before(new_season_start) {
		return int32(d.Year())
	}
	return int32(d.Year() + 1)
}