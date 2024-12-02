package web_scraping

import (
	"reflect"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
)

//Helper function to remove UUID, Created_At and Updated_At values for testing
func stripIDFromPenSummaries(penSummaries []database.CreatePenaltySummaryParams) []database.CreatePenaltySummaryParams {
	stripped := make([]database.CreatePenaltySummaryParams, len(penSummaries))
	for i, summary := range penSummaries {
		stripped[i] = database.CreatePenaltySummaryParams{
			Gameid: summary.Gameid,
			Period: summary.Period,
			Time: summary.Time,
			Team: summary.Team,
			Player: summary.Player,
			PlayerID: summary.PlayerID,
			Penalty: summary.Penalty,
			Pim: summary.Pim,
		}
	}
	return stripped
}

func Test_scrapePenaltySummary(t *testing.T) {
	parsePenTime := func (value string) time.Time{
		t, _ := time.Parse("04:05", value)
		return t
	}

	type args struct {
		doc *goquery.Document
		ID  string
	}

	tests := []struct {
		name string
		args args
		want []database.CreatePenaltySummaryParams
	}{
		{"Penalty Shot", 
			args{
				doc: loadLocalTestFile("/home/rehmke/workspace/github.com/ehmker/hockey_stats/test_files/penalty_shot_example.htm"),
				ID: "202410170MTL",
			},
			[]database.CreatePenaltySummaryParams{
				{Gameid: "202410170MTL", Period: "1st Period", Time: parsePenTime("09:34"), Team: "LAK", Player: "Andreas Englund", PlayerID: "engluan01", Penalty: "Illegal Check to Head", Pim: 2},
				{Gameid: "202410170MTL", Period: "2nd Period", Time: parsePenTime("04:28"), Team: "LAK", Player: "Trevor Lewis", PlayerID: "lewistr01", Penalty: "Roughing", Pim: 2},
				{Gameid: "202410170MTL", Period: "2nd Period", Time: parsePenTime("04:42"), Team: "MTL", Player: "Kirby Dach", PlayerID: "dachki01", Penalty: "Hooking", Pim: 2},
				{Gameid: "202410170MTL", Period: "2nd Period", Time: parsePenTime("10:03"), Team: "LAK", Player: "Caleb Jones", PlayerID: "jonesca01", Penalty: "Interference", Pim: 2},
				{Gameid: "202410170MTL", Period: "2nd Period", Time: parsePenTime("19:07"), Team: "LAK", Player: "Kevin Fiala", PlayerID: "fialake01", Penalty: "Interference", Pim: 2},
				{Gameid: "202410170MTL", Period: "3rd Period", Time: parsePenTime("02:51"), Team: "MTL", Player: "Arber Xhekaj", PlayerID: "xhekaar01", Penalty: "Holding", Pim: 2},
				{Gameid: "202410170MTL", Period: "3rd Period", Time: parsePenTime("07:10"), Team: "MTL", Player: "Lane Hutson", PlayerID: "hutsola01", Penalty: "Interference on breakaway (Penalty Shot)", Pim: 0},
				{Gameid: "202410170MTL", Period: "3rd Period", Time: parsePenTime("18:12"), Team: "LAK", Player: "Trevor Moore", PlayerID: "mooretr01", Penalty: "Tripping", Pim: 2},
				
			},
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scrapePenaltySummary(tt.args.doc, tt.args.ID)

			//remove unique identifiers 
			gotStripped := stripIDFromPenSummaries(got)


			if !reflect.DeepEqual(gotStripped, tt.want) {
				t.Errorf("scrapePenaltySummary()")
				for i, g:= range gotStripped{
					if !reflect.DeepEqual(g, tt.want[i]) {
						t.Errorf("fail record number: %v\nGot: %v\nWant: %v\n", i, g, tt.want[i])
					}
				}
				// t.Errorf("scrapePenaltySummary() = %v, \n\nwant %v", gotStripped, tt.want)
			}
		})
	}
}



