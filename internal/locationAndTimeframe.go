package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/craniacshencil/cachecast/utils"
)

func LocationAndTimeframe(
	w http.ResponseWriter,
	location string,
	startDateString string,
	endDateString string,
) {
	var apiEndpoint string
	var response interface{}

	// Converting startDateString and endDateString to time.Time
	// Checking whether dates valid or not
	var startDate, endDate time.Time
	startDate, err := time.Parse(time.DateOnly, startDateString)
	if err != nil {
		log.Println("ERR: While parsing start-date")
		log.Println(err)
		utils.WriteJSON(w, 404, err)
		return
	}

	endDate, err = time.Parse(time.DateOnly, endDateString)
	if err != nil {
		log.Println("ERR: While parsing start-date")
		log.Println(err)
		utils.WriteJSON(w, 404, err)
		return
	}

	// Error when endDate is before startDate
	if endDate.Before(startDate) {
		log.Println("ERR: end-date is before start-date")
		utils.WriteJSON(w, 404, "end-date is before start-date")
		return
	}

	// Error when startDate is too far into the future
	// - 1 Year after current year
	futureDeadline := time.Now().Add(time.Hour * 24 * 365)

	if endDate.After(futureDeadline) {
		log.Println("ERR: time-frame is far off in the future")
		utils.WriteJSON(w, 404, "time-frame is far off in the future")
		return
	}

	// Error when time-frame b/w startDate and endDate is too large
	// - More than 30 days
	timeFrameDeadline := startDate.Add(time.Hour * 24 * 30)
	if endDate.After(timeFrameDeadline) {
		log.Println("ERR: time-frame larger than 30 days")
		utils.WriteJSON(w, 404, "time-frame larger than 30 days")
		return
	}

	apiEndpoint = fmt.Sprintf(
		"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%s/%s?key=%s&unitGroup=metric&elements=temp,tempmin,tempmax,conditions,datetime&include=days",
		location,
		startDateString,
		endDateString,
		os.Getenv("API_KEY"),
	)

	res, err := http.Get(apiEndpoint)
	if err != nil {
		log.Println("ERR: While contacting third party API")
		log.Println(err)
		utils.WriteJSON(w, 404, err)
		return
	}

	err = utils.ParseBody(res, &response)
	if err != nil {
		log.Println("ERR: In response sent by third party API")
		log.Println(err)
		utils.WriteJSON(w, 404, err)
		return
	}
	utils.WriteJSON(w, 200, response)
}
