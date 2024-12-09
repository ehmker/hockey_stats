package web_scraping

import (
	"log"
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
)

func Test_scrapeGoalieStats(t *testing.T) {
	type args struct {
		doc *goquery.Document
		ID  string
	}
	tests := []struct {
		name string
		args args
		want []database.CreateGoalieStatsParams
	}{
		{
			name: "Pulled Goalie Test",
			args: args{
				doc: loadLocalTestFile("../../test_files/202412030BUF_pulled_goalie_example.htm"),
				ID: "202412030BUF",
			},
			want: []database.CreateGoalieStatsParams{},

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scrapeGoalieStats(tt.args.doc, tt.args.ID)
			for _, s := range got{
				log.Printf("got: %v\n", s)
			}
			if got := scrapeGoalieStats(tt.args.doc, tt.args.ID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scrapeGoalieStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
