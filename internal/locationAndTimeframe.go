package internal

import (
	"errors"
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
	var startDate, endDate time.Time
	var err error
	var apiEndpoint string
	var response interface{}

	// Converting startDateString and endDateString to time.Time
	// Checking whether dates valid or not
	startDate, err = ValidDate(startDateString)
	if err != nil {
		utils.WriteJSON(w, 404, err)
		return
	}

	endDate, err = ValidDate(endDateString)
	if err != nil {
		utils.WriteJSON(w, 404, err)
		return
	}

	// Error when endDate is before startDate
	if endDate.Before(startDate) {
		log.Println("ERR: end-date is before start-date")
		utils.WriteJSON(w, 404, "end-date is before start-date")
		return
	}

	// Error when endDate is too far into the future
	// - 1 Year after current year
	err = FutureDeadline(endDate)
	if err != nil {
		utils.WriteJSON(w, 404, err.Error())
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

func ValidDate(dateString string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		log.Println("ERR: While parsing:", dateString)
		log.Println(err)
		return date, err // can't use nil for some reason
	}
	return date, nil
}

func FutureDeadline(date time.Time) error {
	futureDeadline := time.Now().Add(time.Hour * 24 * 365)

	if date.After(futureDeadline) {
		log.Println("ERR: time-frame is far off in the future")
		return errors.New("date is far off in the future")
	}
	return nil
}
