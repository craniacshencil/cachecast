package internal

import (
	"log"
	"net/http"

	"github.com/craniacshencil/cachecast/utils"
)

func GetWeather(w http.ResponseWriter, r *http.Request) {
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
		LocationAndTimeframe(w, location, startDate, endDate)
	} else if onlyDate != "" {
		// Case for when location and date1 is given
		LocationAndDay(w, location, onlyDate, onlyTime)
	} else {
		OnlyLocation(w, location)
	}
}
