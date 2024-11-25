package web_scraping

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/ehmker/hockey_stats/internal/shared"
	"github.com/google/uuid"
)

func AddPenaltySummary(s shared.State) {
	// Open the local HTML file
	file, err := os.Open("example_2.htm")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	penaltySummary := scrapePenaltySummary(file)

	for _, penalty := range penaltySummary{
		_, err := s.DB.CreatePenaltySummary(context.Background(), penalty)
		if err != nil {
			log.Println("unable to add to penalty_summaries: ", err)
		}
	}

}

func scrapePenaltySummary (f *os.File) []database.CreatePenaltySummaryParams{
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	var penaltySummarySlice []database.CreatePenaltySummaryParams
	
	//Extract the first penalty period from table 
	period := doc.Find("#penalty thead th").Text()

	doc.Find("#penalty tbody tr").Each(func(i int, row *goquery.Selection) {
		// Check if row is a header row (by presence of <th> instead of <td>)
		if row.Find("th").Length() > 0 {
			// It's a header row, update the period information
			period = row.Find("th").Text()
		} else {

			penaltySum := database.CreatePenaltySummaryParams {
				ID: uuid.New(),
				Gameid: "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Period: period,
			}

			// Extract time
			timeLayout := "04:05"
			pen_time_string := row.Find("td").Eq(0).Text()
			pen_time, err := time.Parse(timeLayout, pen_time_string)
			if err != nil{
				fmt.Printf("unable to parse time for time scored: %v\n", err)
					return
			}
			penaltySum.Time = pen_time

			// Extract team
			penaltySum.Team = row.Find("td").Eq(1).Text()

			// Extract player
			player := getPlayerDetailFromCell(row.Find("td").Eq(2).Find("a"))
			//if bench infraction no player ID.  set ID to bench as well
			if player.name == "Bench" {
				player.id = "Bench"
			}
			penaltySum.Player = player.name
			penaltySum.PlayerID = player.id

			// Extract Infraction
			penaltySum.Penalty = row.Find("td").Eq(3).Text()

			// Extract PIMs
			pim_string := row.Find("td").Eq(4).Text()  //ex "2 min"
			pim_int, err := strconv.Atoi(strings.Split(pim_string, " ")[0])
			if err != nil {
				log.Panicln("unable to convert PIMs to int: ", err)
			}
			penaltySum.Pim = int32(pim_int)

			penaltySummarySlice = append(penaltySummarySlice, penaltySum)

		}
	})

	return penaltySummarySlice
}