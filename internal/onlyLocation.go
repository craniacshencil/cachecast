package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/craniacshencil/cachecast/utils"
)

func OnlyLocation(w http.ResponseWriter, location string) {
	apiEndpoint := fmt.Sprintf(
		"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s&unitGroup=metric",
		location,
		os.Getenv("API_KEY"),
	)
	res, err := http.Get(apiEndpoint)
	if err != nil {
		log.Println("ERR: While contacting third party api")
		log.Println(err)
		utils.WriteJSON(w, http.StatusNotFound, err)
		return
	}
	var response interface{}
	utils.ParseBody(res, &response)
	if response == nil {
		log.Println("ERR: Invalid location")
		utils.WriteJSON(w, http.StatusNotFound, "invalid location, check for errors")
		return
	}
	utils.WriteJSON(w, 200, response)
	// https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/[location]/[date1]/[date2]?key=YOUR_API_KEY
	// utils.WriteJSON(w, 200, r.PostForm)
}
