package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/craniacshencil/cachecast/utils"
)

func (c *CacheClient) GetWeather(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var response JSONWrapper
	var cacheStatus, elapsed string
	start := time.Now()
	ctx := context.Background()
	location := r.FormValue("location")
	onlyDate := r.FormValue("only-date")
	onlyTime := r.FormValue("only-time")
	if onlyTime != "" {
		onlyTime = onlyTime + ":00"
	}
	startDate := r.FormValue("start-date")
	endDate := r.FormValue("end-date")

	// Form validation
	err := validForm(location, onlyDate, onlyTime, startDate, endDate)
	if err != nil {
		displayWeather(w, Weather{}, "", "", err, http.StatusBadRequest)
		return
	}

	// Get API URL and cacheKey
	apiEndpoint := getApiURL(location, onlyDate, onlyTime, startDate, endDate)
	cacheKey := getCacheKey(location, onlyDate, onlyTime, startDate, endDate)

	err = c.searchCache(ctx, cacheKey, &response)
	if err != nil {
		log.Println(err.Error())
		// Don't need to end flow because of cache failing
	} else {
		elapsed = time.Since(start).String()
		cacheStatus = "cache-hit"
		log.Printf("Exec: %s", elapsed)
		displayWeather(w, response.Data, cacheStatus, elapsed, nil, http.StatusOK)
		return
	}

	statusCode, err := fetchData(ctx, apiEndpoint, &response)
	if err != nil {
		displayWeather(w, Weather{}, "", "", err, statusCode)
		return
	}

	go func() {
		if err = c.storeInCache(ctx, cacheKey, &response); err != nil {
			log.Println(err)
		}
	}()

	elapsed = time.Since(start).String()
	cacheStatus = "cache-miss"
	log.Printf("Exec: %s", elapsed)
	displayWeather(w, response.Data, cacheStatus, elapsed, nil, http.StatusOK)
}

func getApiURL(location, onlyDate, onlyTime, startDate, endDate string) (apiEndpoint string) {
	if startDate != "" && endDate != "" {
		// URL for timeFrame case
		apiEndpoint = fmt.Sprintf(
			"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%s/%s?key=%s&unitGroup=metric&elements=temp,tempmin,tempmax,conditions,datetime&include=days",
			location,
			startDate,
			endDate,
			os.Getenv("API_KEY"),
		)
	} else if onlyTime != "" {
		// URL for single day with time specified
		apiEndpoint = fmt.Sprintf(
			"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%sT%s?key=%s&unitGroup=metric&include=current&elements=temp,tempmin,tempmax,conditions,datetime",
			location,
			onlyDate,
			onlyTime,
			os.Getenv("API_KEY"),
		)
	} else if onlyDate != "" {
		// URL for single day with no time specified
		apiEndpoint = fmt.Sprintf(
			"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%s?key=%s&unitGroup=metric&elements=temp,tempmin,tempmax,conditions,datetime&include=days",
			location,
			onlyDate,
			os.Getenv("API_KEY"),
		)
	} else {
		// URL for only location specified
		apiEndpoint = fmt.Sprintf(
			"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s&unitGroup=metric&elements=temp,tempmin,tempmax,conditions,datetime&include=days",
			location,
			os.Getenv("API_KEY"),
		)
	}
	return apiEndpoint
}

func getCacheKey(location, onlyDate, onlyTime, startDate, endDate string) (cacheKey string) {
	if startDate != "" && endDate != "" {
		cacheKey = fmt.Sprintf("%s_%s_%s", location, startDate, endDate)
	} else if onlyTime != "" {
		cacheKey = fmt.Sprintf("%s_%s_%s", location, onlyDate, onlyTime)
	} else if onlyDate != "" {
		cacheKey = fmt.Sprintf("%s_%s", location, onlyDate)
	} else {
		cacheKey = location
	}
	return cacheKey
}

func fetchData(
	ctx context.Context,
	apiEndpoint string,
	response *JSONWrapper,
) (status int, err error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	req, err := http.NewRequestWithContext(timeout, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		err = errors.Join(err, errors.New("while defining request"))
		return http.StatusBadRequest, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.Join(err, errors.New("request timed-out"))
		return http.StatusRequestTimeout, err
	}

	utils.ParseBody(res, &response.Data)
	// Invalid location
	if response.Data.ResolvedAddress == "" {
		err = errors.Join(err, errors.New("invalid location, check for errors"))
		return http.StatusNotFound, err
	}

	return 200, nil
}
