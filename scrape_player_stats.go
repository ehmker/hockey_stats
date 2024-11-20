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

func getTeamFromStatTableID (id string) string {
	return strings.Split(id, "_")[1]
}

func getIntStatFromCell (stat string, s *goquery.Selection) int32 {
	selection_string := "td[data-stat='" + stat + "']"
	str := s.Find(selection_string).Text()
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("error converting %s: %v\n", stat, err)
		return 0
	}
	return int32(num)
}

func getPlayerFromStatCell (s *goquery.Selection) Player {
	return getPlayerDetailFromCell(s.Find("td[data-stat='player'] a"))
}

// Extract the time in MM:SS format and returns the time as seconds 
func getTimeStatFromCell (stat string, s *goquery.Selection) int32 {
	selection_string := "td[data-stat='" + stat + "']"
	time_string := strings.Split(s.Find(selection_string).Text(),":")
	minutes, err  := strconv.Atoi(time_string[0])
	if err != nil {
		log.Printf("error converting minute string to int [%v]: %v", time_string[0], err)
		return 0
	}
	seconds, err := strconv.Atoi(time_string[1])
	if err != nil {
		log.Printf("error converting seconds string to int [%v]: %v", time_string[1], err)
		return 0
	}

	return int32(minutes * 60 + seconds)

}

func AddPlayerStats(s state) {
	file, err := os.Open("example_2.htm")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	playerStatLines := ScrapePlayerStats(file)

	for _, statline := range playerStatLines{
		_, err = s.db.CreateSkaterGameStats(context.Background(), statline)
		if err != nil {
			log.Printf("error adding skater game stats: %v", err)
		}
	}
}

func ScrapePlayerStats (f *os.File) []database.CreateSkaterGameStatsParams {
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}
	
	var skaterStatsSlice []database.CreateSkaterGameStatsParams

	//stat tables are id as "all_[team]_stats"
	doc.Find("div[id^='all_'][id$='_skaters']").Each(func(i int, div *goquery.Selection) {
	// doc.Find("div[id='all_LAK_skaters']").Each(func(i int, div *goquery.Selection) {
		// Extract team short name
		var team string
		divID, exist := div.Attr("id")
		if exist == true{
			team = getTeamFromStatTableID(divID)
		}

		// Extract each player's statline for the game
		div.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
			p := getPlayerFromStatCell(row)
			skaterStats := database.CreateSkaterGameStatsParams{
				ID: uuid.New(),
				Gameid: "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Team: team,
				PlayerName: p.name,
				Playerid: p.id,
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