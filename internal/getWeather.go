package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/craniacshencil/cachecast/utils"
)

func GetWeather(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	location := r.FormValue("location")
	if location != "" {
		apiEndpoint := fmt.Sprintf(
			"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s",
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
		utils.WriteJSON(w, 200, response)
	}
	// https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/[location]/[date1]/[date2]?key=YOUR_API_KEY
	// utils.WriteJSON(w, 200, r.PostForm)
}
