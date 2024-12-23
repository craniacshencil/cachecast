package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/craniacshencil/cachecast/utils"
)

func LocationAndDay(w http.ResponseWriter, location string, date string, time string) {
	var apiEndpoint string

	if time != "" {
		// Format time according to APIs requirement
		time = time + ":00"
		apiEndpoint = fmt.Sprintf(
			"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%sT%s?key=%s&unitGroup=metric&include=current",
			location,
			date,
			time,
			os.Getenv("API_KEY"),
		)
	} else {
		apiEndpoint = fmt.Sprintf(
			"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%s?key=%s&unitGroup=metric",
			location,
			date,
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
