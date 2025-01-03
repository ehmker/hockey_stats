package web_scraping

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ehmker/hockey_stats/internal/database"
	"github.com/ehmker/hockey_stats/internal/shared"
)

type GameLink struct {
	Url string
	Gameid string
}

type Player struct {
	name string
	id string
}

type shotCoordinates struct{
	x_loc int
	y_loc int
}


// extracts the team short name from passed goalie table id string
func getTeamFromGoalieStatTableID (s string) string {
	//expected format "[team]_goalies"
	return strings.Split(s, "_")[0]
}

// extracts the team short name from passed player stat table id string
func getTeamFromPlayerStatTableID (id string) string {
	// Expected format "all_[team]_stats"
	return strings.Split(id, "_")[1]
}

// extracts int value from given "data-stat" value
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

func getTextStatFromCell_CanBeNull(stat string, s *goquery.Selection) sql.NullString {
	selection_string := "td[data-stat='" + stat + "']"
	str := strings.TrimSpace(s.Find(selection_string).Text())
	if str == "" {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		String: str,
		Valid: true,
	}
}

// helper function for extracting player detail
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



// extracts the player name and playerid from cell
func getPlayerDetailFromCell(player *goquery.Selection) Player {
	playerName := strings.TrimSpace(player.Text())
	playerHref, _ := player.Attr("href")
	playerId := path.Base(playerHref)
	playerId = strings.TrimSuffix(playerId, ".html")
	return Player{
		name: playerName,
		id: playerId,
	}
}




// Extracts the shot coordinates from given string
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

func intToBool(i interface{}) bool {
	switch i.(type) {
	case int32:
		return i != int32(0)
	default:
		log.Printf("unsupported type: %d", i)
		return false
	}
}

func AddGameToDB(s shared.State, game GameLink ) (database.CreateGameResultParams, error){
	resp, err := http.Get(game.Url)
	if err != nil {
		log.Println(err)
		return database.CreateGameResultParams{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		return database.CreateGameResultParams{}, err
	}

	// Removing all comment sections from the response 
	// Wanted information is kept in commented sections which prevents goquery from finding it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response body:", err)
		return database.CreateGameResultParams{}, err
	}
	bodyString := string(body)
	cleanedBody := strings.ReplaceAll(bodyString, "<!--", "")
	cleanedBody = strings.ReplaceAll(cleanedBody, "-->", "")


	doc, err := goquery.NewDocumentFromReader(strings.NewReader(cleanedBody))
	if err != nil {
		log.Printf("error getting document from reader: %v\n", err)
		return database.CreateGameResultParams{}, err
	}
	
	gameDetail, err := addGameResults(s, doc, game.Gameid)
	if err != nil {
		// early out if game already has been added to db
		log.Printf("error adding game results to database: %v\n", err)
		return database.CreateGameResultParams{}, err
	}
	addPlayerStats(s, doc, game.Gameid, gameDetail.Season)
	addGoalieStats(s, doc, game.Gameid, gameDetail.Season)
	addPenaltySummary(s, doc, game.Gameid)
	addScoringSummaryToDB(s, doc, game.Gameid)
	addShotLocationsToDB(s, doc, game.Gameid)
	return gameDetail, nil
}



func loadLocalTestFile (f string) *goquery.Document {
	file, err := os.Open(f)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}
	return doc
}
