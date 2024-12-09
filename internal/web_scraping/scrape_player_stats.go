package web_scraping

import (
	"context"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/ehmker/hockey_stats/internal/shared"
	"github.com/google/uuid"
)


func addPlayerStats(s shared.State, doc *goquery.Document, gameID string, season int32) {
	playerStatLines := ScrapePlayerStats(doc, gameID)

	for _, statline := range playerStatLines{
		statline.Season = season
		_, err := s.DB.CreateSkaterGameStats(context.Background(), statline)
		if err != nil {
			log.Printf("error adding skater game stats: %v", err)
		}
	}
}

func ScrapePlayerStats (doc *goquery.Document, ID string) []database.CreateSkaterGameStatsParams {
	// doc, err := goquery.NewDocumentFromReader(f)
	// if err != nil {
	// 	log.Fatalf("Error parsing HTML: %v", err)
	// }
	
	var skaterStatsSlice []database.CreateSkaterGameStatsParams

	//stat tables are id as "all_[team]_stats"
	doc.Find("div[id^='all_'][id$='_skaters']").Each(func(i int, div *goquery.Selection) {
	// doc.Find("div[id='all_LAK_skaters']").Each(func(i int, div *goquery.Selection) {
		// Extract team short name
		var team string
		divID, exist := div.Attr("id")
		if exist == true{
			team = getTeamFromPlayerStatTableID(divID)
		}

		// Extract each player's statline for the game
		div.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
			p := getPlayerFromStatCell(row)
			skaterStats := database.CreateSkaterGameStatsParams{
				ID: uuid.New(),
				GameID: ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Team: team,
				PlayerName: p.name,
				PlayerID: p.id,
				Goals: getIntStatFromCell("goals", row),
				Assists: getIntStatFromCell("assists", row),
				Points: getIntStatFromCell("points", row),
				PlusMinus: getIntStatFromCell("plus_minus", row),
				PenMins: getIntStatFromCell("pen_min", row),
				GoalsEv: getIntStatFromCell("goals_ev", row),
				GoalsPp: getIntStatFromCell("goals_pp", row),
				GoalsSh: getIntStatFromCell("goals_sh", row),
				GoalsGw: getIntStatFromCell("goals_gw", row),
				AssistsEv: getIntStatFromCell("assists_ev", row),
				AssistsPp: getIntStatFromCell("assists_pp", row),
				AssistsSh: getIntStatFromCell("assists_sh", row),
				Shots: getIntStatFromCell("shots", row),
				Shifts: getIntStatFromCell("shifts", row),
				TimeOnIce: getTimeStatFromCell("time_on_ice", row),
				
			}
			
			skaterStatsSlice = append(skaterStatsSlice, skaterStats)
			})
	})
	return skaterStatsSlice
}