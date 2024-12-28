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

func (c *CacheClient) OnlyLocation(w http.ResponseWriter, location string) {
	var response JSONWrapper
	start := time.Now()
	ctx := context.Background()
	// Add contexts with timeout and change fetchAPI code later

	err := c.searchCache(ctx, location, &response)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("cache-hit: %s", time.Since(start))
		utils.WriteJSON(w, http.StatusOK, &response.Data)
		return
	}

	statusCode, err := fetchlocationAPI(ctx, location, &response)
	if err != nil {
		utils.WriteJSON(w, statusCode, err.Error())
		return
	}
	utils.WriteJSON(w, statusCode, &response.Data)
	log.Printf("Exec: %s", time.Since(start))

	go func() {
		if err = c.storeInCache(ctx, location, &response); err != nil {
			log.Println(err)
		}
	}()
}

func fetchlocationAPI(
	ctx context.Context,
	location string,
	response *JSONWrapper,
) (status int, err error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	apiEndpoint := fmt.Sprintf(
		"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s&unitGroup=metric&elements=temp,tempmin,tempmax,conditions,datetime&include=days",
		location,
		os.Getenv("API_KEY"),
	)

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
	if response.Data == nil {
		err = errors.Join(err, errors.New("invalid location, check for errors"))
		return http.StatusNotFound, err
	}

	return 200, nil
}
