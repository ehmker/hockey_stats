package web_scraping

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/ehmker/hockey_stats/internal/shared"
	"github.com/google/uuid"
)


func AddScoringSummaryToDB(s shared.State) {
	// Open the local HTML file
	file, err := os.Open("example_2.htm")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	scoringSummary := scrapeScoringSummary(file)

	for _, score := range scoringSummary{
		_, err = s.DB.CreateScoringSummary(context.Background(), score)
		if err != nil {
			log.Panicf("error adding scoring summary: %v\n", err)
		}
	}

}

func scrapeScoringSummary(f *os.File) []database.CreateScoringSummaryParams {
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	var scoringSummarySlice []database.CreateScoringSummaryParams
	
	//Extract the first scoring period from table 
	period := doc.Find("#scoring thead th").Text()

	doc.Find("#scoring tbody tr").Each(func(i int, row *goquery.Selection) {
		// Check if row is a header row (by presence of <th> instead of <td>)
		if row.Find("th").Length() > 0 {
			// It's a header row, update the period information
			period = row.Find("th").Text()
		} else {

			scoringSum := database.CreateScoringSummaryParams {
				ID: uuid.New(),
				Gameid: "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Period: period,
			}


			// It's a data row, extract time, team, player, etc.
			// Expected format
			// time_scored | team | *optional* EN | scoring player | *optional* assist1, assist2

			// Extract time
			timeLayout := "04:05"
			time_scored_string := row.Find("td").Eq(0).Text()
			time_scored, err := time.Parse(timeLayout, time_scored_string)
			if err != nil{
				fmt.Printf("unable to parse time for time scored: %v\n", err)
					return
			}
			scoringSum.Time = time_scored

			// Extract Team 
			team := row.Find("td").Eq(1).Text()
			scoringSum.Team = team 

			// Extract if Empty Net Goal 
			if strings.TrimSpace(row.Find("td").Eq(2).Text()) == "EN" {
				scoringSum.EmptyNet = true
			} else {
				scoringSum.EmptyNet = false
			}

			// Extract goal scorer information
			var scoringPlayers []Player
			goal_scorer := getPlayerDetailFromCell(row.Find("td").Eq(3).Find("a"))
			scoringPlayers = append(scoringPlayers, goal_scorer)
			
			// Extract assist player information
			row.Find("td").Eq(4).Find("a").Each(func(i int, s *goquery.Selection) {
				assist_player := getPlayerDetailFromCell(s)
				scoringPlayers = append(scoringPlayers, assist_player)
			})

			for i, p := range scoringPlayers {
				switch i {
				case 0: // Goal Scorer
					scoringSum.ScoringPlayer = p.name
					scoringSum.ScoringPlayerID = p.id
				
				case 1: // Primary Assist
					scoringSum.FirstAssist = sql.NullString{
						String: p.name,
						Valid: true,
					}
					scoringSum.FirstAssistID = sql.NullString{
						String: p.id,
						Valid: true,
					}
				case 2: // Secondary Assist
					scoringSum.SecondAssist = sql.NullString{
						String: p.name,
						Valid: true,
					}
					scoringSum.SecondAssistID = sql.NullString{
						String: p.id,
						Valid: true,
					}
				}
			}
			scoringSummarySlice = append(scoringSummarySlice, scoringSum)
			}
		})

		return scoringSummarySlice
	}


