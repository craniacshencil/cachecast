package internal

import (
	"errors"
	"time"
)

func validForm(location, onlyDate, onlyTime, startDate, endDate string) (err error) {
	// Error when location is not entered
	if location == "" {
		return errors.New("no location was entered")
	}

	// Error if user tries to hit both timeframe and daily weather updates
	if (onlyDate != "" || onlyTime != "") && (endDate != "" || startDate != "") {
		return errors.New(
			"you are trying both daily and timeframe weather updates. Choose one",
		)
	}

	// Error if user only fills time-field and not date
	if onlyDate == "" && onlyTime != "" {
		return errors.New("no date entered")
	}

	// Errors if one of start-date or end-date is not entered
	if startDate == "" && endDate != "" {
		return errors.New("no start-date was entered")
	} else if endDate == "" && startDate != "" {
		return errors.New("no end-date was entered")
	}

	// Form Validation for location+day
	if onlyDate != "" {
		return locationAndDayValidation(onlyDate, onlyTime)
	}
	// Form validation for location+timeFrame
	if startDate != "" && endDate != "" {
		return locationAndTimeframeValidation(startDate, endDate)
	}
	return nil
}

func locationAndDayValidation(onlyDate, onlyTime string) (err error) {
	date, err := validDate(onlyDate)
	if err != nil {
		return err
	}

	err = futureDeadline(date)
	if err != nil {
		return err
	}

	// Form Validation for location+day+time
	if onlyTime != "" {
		_, err = time.Parse(time.TimeOnly, onlyTime)
		if err != nil {
			return errors.New("couldn't parse time")
		}
	}
	return nil
}

func locationAndTimeframeValidation(startDate, endDate string) (err error) {
	startDateVal, err := validDate(startDate)
	if err != nil {
		return err
	}

	endDateVal, err := validDate(endDate)
	if err != nil {
		return err
	}

	// Error when endDate is before startDate
	if endDateVal.Before(startDateVal) {
		return errors.New("end-date is before start-date")
	}

	// Error when endDate is too far into the future
	// - 1 Year after current year
	err = futureDeadline(endDateVal)
	if err != nil {
		return err
	}

	// Error when time-frame b/w startDate and endDate is too large
	// - More than 30 days
	timeFrameDeadline := startDateVal.Add(time.Hour * 24 * 31)
	if endDateVal.After(timeFrameDeadline) {
		return errors.New("time-frame larger than 30 days")
	}
	return nil
}

func validDate(dateString string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		return date, err // can't use nil for some reason
	}
	return date, nil
}

func futureDeadline(date time.Time) error {
	futureDeadline := time.Now().Add(time.Hour * 24 * 365)

	if date.After(futureDeadline) {
		return errors.New("date is far off in the future")
	}
	return nil
}
