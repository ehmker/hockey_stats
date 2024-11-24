package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/google/uuid"
)

func AddGoalieStats(s state) {
	file, err := os.Open("example_2.htm")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	goalieStats := scrapeGoalieStats(file)

	for _, statline := range goalieStats {
		_, err := s.db.CreateGoalieStats(context.Background(), statline)
		if err != nil {
			log.Println(err)
		}
	}
}

func scrapeGoalieStats (f *os.File) []database.CreateGoalieStatsParams{
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	var goalieStatsSlice []database.CreateGoalieStatsParams

	doc.Find("table[id$='_goalies']").Each(func(i int, table *goquery.Selection) {
		team_table_id, _ := table.Attr("id")
		team := getTeamFromGoalieStatTableID(team_table_id)


		table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
			if row.Find("td").First().Text() == "Empty Net" {
				// skip records for empty net
				return
			}
			goalie := getPlayerFromStatCell(row)
			goalieStat := database.CreateGoalieStatsParams{
				ID: uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Gameid: "1",
				Team: team,
				PlayerName: goalie.name,
				Playerid: goalie.id,
				Decision: getTextStatFromCell_CanBeNull("decision",row),
				GoalsAgainst: getIntStatFromCell("goals_against", row),
				ShotsAgainst: getIntStatFromCell("shots_against", row),
				Saves: getIntStatFromCell("saves", row),
				Shutout: intToBool(getIntStatFromCell("shutouts", row)), // shutout either 1 or 0
				PenMins: getIntStatFromCell("pen_min", row),
				TimeOnIce: getTimeStatFromCell("time_on_ice", row),
			}


			goalieStatsSlice = append(goalieStatsSlice, goalieStat)
		})
	})
	return goalieStatsSlice
}