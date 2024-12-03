package shared

import (
	"context"
	"log"
	"time"
)

func (s State) GetLastScrapedDate () time.Time {
	lastScraped, err := s.DB.GetLastScrapedDateFromDB(context.Background())
	if err != nil {
		log.Printf("error getting last scraped date from DB: %v\n", err)
		return time.Now()
	}
	
	lastScraped_time, err := time.Parse("2006/01/02", lastScraped)
	if err != nil {
		log.Printf("error parsing last scraped date as time: %v\n", err)
		return time.Now()
	}

	return lastScraped_time
}

func (s State) SetLastScrapedDate () {
	lastScrapedGameDate, err := s.DB.GetDateOfLastResult(context.Background())
	if err != nil {
		log.Printf("error getting last game date: %v\n", err)
		lastScrapedGameDate = time.Now()
	}

	err = s.DB.UpdateLastScrapedDate(context.Background(), lastScrapedGameDate.Format("2006/01/02"))
	if err != nil{
		log.Printf("error updating last scraped date in DB: %v\n", err)
	}
}