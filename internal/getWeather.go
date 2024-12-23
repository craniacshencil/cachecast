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
	if location == "" {
		log.Println("ERR: No location entered")
		utils.WriteJSON(w, 404, "no location was entered")
	} else if onlyDate != "" {
		// Case for when location and date1 is given
		LocationAndDay(w, location, onlyDate, onlyTime)
	} else {
		OnlyLocation(w, location)
	}
}
