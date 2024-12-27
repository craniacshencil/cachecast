package internal

import (
	"log"
	"net/http"

	"github.com/craniacshencil/cachecast/utils"
	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	rdsClient *redis.Client
}

func NewCacheClient(client *redis.Client) *CacheClient {
	return &CacheClient{rdsClient: client}
}

func (c *CacheClient) GetWeather(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	location := r.FormValue("location")
	onlyDate := r.FormValue("only-date")
	onlyTime := r.FormValue("only-time")
	startDate := r.FormValue("start-date")
	endDate := r.FormValue("end-date")

	// Error when location is not entered
	if location == "" {
		log.Println("ERR: No location entered")
		utils.WriteJSON(w, 404, "no location was entered")
		return
	}

	// Error if user only fills time-field and not date
	if onlyDate == "" && onlyTime != "" {
		log.Println("ERR: No date entered, only time entered")
		utils.WriteJSON(w, 404, "no date entered")
		return
	}

	// Error if user tries to hit both timeframe and daily weather updates
	if onlyDate != "" && (endDate != "" || startDate != "") {
		log.Println("ERR: You can get only one weather update at a time")
		utils.WriteJSON(
			w,
			404,
			"you are trying both daily and timeframe weather updates. Choose one",
		)
		return
	}

	// Errors if one of start-date or end-date is not entered
	if startDate == "" && endDate != "" {
		log.Println("ERR: No start-date entered, only end-date entered")
		utils.WriteJSON(w, 404, "no start-date was entered")
		return
	} else if endDate == "" && startDate != "" {
		log.Println("ERR: No end-date entered, only start-date entered")
		utils.WriteJSON(w, 404, "no end-date was entered")
		return
	}

	// Handling the three cases
	if startDate != "" && endDate != "" {
		c.LocationAndTimeframe(w, location, startDate, endDate)
	} else if onlyDate != "" {
		// Case for when location and date1 is given
		c.LocationAndDay(w, location, onlyDate, onlyTime)
	} else {
		c.OnlyLocation(w, location)
	}
}
