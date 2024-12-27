package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/craniacshencil/cachecast/utils"
)

func (c *CacheClient) LocationAndDay(
	w http.ResponseWriter,
	location string,
	dateString string,
	timeString string,
) {
	var apiEndpoint string
	var err error

	date, err := ValidDate(dateString)
	if err != nil {
		utils.WriteJSON(w, 404, err)
		return
	}

	err = FutureDeadline(date)
	if err != nil {
		utils.WriteJSON(w, 404, err.Error())
		return
	}

	if timeString != "" {
		// Format timeString according to APIs requirement and check whether valid
		timeString = timeString + ":00"
		_, err := time.Parse(time.TimeOnly, timeString)
		if err != nil {
			log.Println("ERR: Time string is not valid")
			log.Println(err)
			utils.WriteJSON(w, 404, "invalid time string")
			return
		}

		apiEndpoint = fmt.Sprintf(
			"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%sT%s?key=%s&unitGroup=metric&include=current&elements=temp,tempmin,tempmax,conditions,datetime",
			location,
			dateString,
			timeString,
			os.Getenv("API_KEY"),
		)
	} else {
		apiEndpoint = fmt.Sprintf(
			"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%s?key=%s&unitGroup=metric&elements=temp,tempmin,tempmax,conditions,datetime",
			location,
			dateString,
			os.Getenv("API_KEY"),
		)
	}
	res, err := http.Get(apiEndpoint)
	if err != nil {
		log.Println("ERR: While contacting third party API")
		log.Println(err)
		utils.WriteJSON(w, 404, err)
		return
	}

	var response interface{}
	err = utils.ParseBody(res, &response)
	if err != nil {
		log.Println("ERR: In response sent by third party API")
		log.Println(err)
		utils.WriteJSON(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, 200, response)
}
