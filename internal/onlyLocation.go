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
	var response utils.JSONWrapper
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

	statusCode, err := fetchlocationAPI(location, &response)
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

func (c *CacheClient) searchCache(
	ctx context.Context,
	location string,
	response *utils.JSONWrapper,
) (err error) {
	getTransaction := c.rdsClient.Get(ctx, location)
	if getTransaction.Err() != nil {
		return errors.New("cache-miss")
	}

	binaryRedisHit, err := getTransaction.Bytes()
	if err != nil {
		return errors.Join(err, errors.New("while converting cache-hit values to bytes"))
	}

	err = response.UnmarshalBinary(binaryRedisHit)
	if err != nil {
		return errors.Join(err, errors.New("while unmarshalling bytes from cache-hit"))
	}
	return nil
}

func fetchlocationAPI(
	location string,
	response *utils.JSONWrapper,
) (status int, err error) {
	apiEndpoint := fmt.Sprintf(
		"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s&unitGroup=metric&elements=temp,tempmin,tempmax,conditions,datetime&include=days",
		location,
		os.Getenv("API_KEY"),
	)

	res, err := http.Get(apiEndpoint)
	if err != nil {
		err = errors.Join(err, errors.New("while contacting third party api"))
		return http.StatusNotFound, err
	}
	utils.ParseBody(res, &response.Data)

	if response.Data == nil {
		err = errors.Join(err, errors.New("invalid location, check for errors"))
		return http.StatusNotFound, err
	}

	return 200, nil
}

func (c *CacheClient) storeInCache(
	ctx context.Context,
	location string,
	response *utils.JSONWrapper,
) (err error) {
	resBytes, err := response.MarshalBinary()
	if err != nil {
		return errors.Join(errors.New("while serializing response for redis"), err)
	}
	transactionStatus := c.rdsClient.Set(ctx, location, resBytes, time.Hour)
	if transactionStatus.Err() != nil {
		return errors.Join(
			errors.New("while storing response in redis"),
			transactionStatus.Err(),
		)
	}
	return nil
}
