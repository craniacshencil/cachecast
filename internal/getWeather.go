package internal

import (
	"errors"
	"net/http"

	"github.com/craniacshencil/cachecast/utils"
)

func (c *CacheClient) GetWeather(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	location := r.FormValue("location")
	onlyDate := r.FormValue("only-date")
	onlyTime := r.FormValue("only-time")
	startDate := r.FormValue("start-date")
	endDate := r.FormValue("end-date")

	// Form validation
	err := validForm(location, onlyDate, onlyTime, startDate, endDate)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err.Error())
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

func validForm(location, onlyDate, onlyTime, startDate, endDate string) (err error) {
	// Error when location is not entered
	if location == "" {
		return errors.New("no location was entered")
	}

	// Error if user only fills time-field and not date
	if onlyDate == "" && onlyTime != "" {
		return errors.New("no date entered")
	}

	// Error if user tries to hit both timeframe and daily weather updates
	if onlyDate != "" && (endDate != "" || startDate != "") {
		return errors.New(
			"you are trying both daily and timeframe weather updates. Choose one",
		)
	}
	// Errors if one of start-date or end-date is not entered
	if startDate == "" && endDate != "" {
		return errors.New("no start-date was entered")
	} else if endDate == "" && startDate != "" {
		return errors.New("no end-date was entered")
	}

	return nil
}
