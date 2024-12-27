package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/craniacshencil/cachecast/utils"
)

func (c *CacheClient) OnlyLocation(w http.ResponseWriter, location string) {
	start := time.Now()
	var response utils.JSONWrapper

	// Check in cache first
	ctx := context.Background()
	getTransaction := c.rdsClient.Get(ctx, location)
	if getTransaction.Err() != nil {
		log.Println("cache miss")
	} else {
		log.Println("cache hit")
		binaryRedisHit, err := getTransaction.Bytes()
		if err != nil {
			log.Println("ERR: While converting cache-hit value to bytes")
			log.Println(err)
			// Find a way to break out of here to contact API or give error here?? idk
		}
		err = response.UnmarshalBinary(binaryRedisHit)
		if err != nil {
			log.Println("ERR: While unmarshalling bytes from cache-hit")
			log.Println(err)
			// Find a way to break out of here to contact API or give error here?? idk
		}
		log.Printf("Exec: %s", time.Since(start))
		utils.WriteJSON(w, 200, response.Data)
		return
	}

	// GET endpoint and store data
	apiEndpoint := fmt.Sprintf(
		"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s&unitGroup=metric&elements=temp,tempmin,tempmax,conditions,datetime&include=days",
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
	utils.ParseBody(res, &response.Data)

	if response.Data == nil {
		log.Println("ERR: Invalid location")
		utils.WriteJSON(w, http.StatusNotFound, "invalid location, check for errors")
		return
	}

	resBytes, err := response.MarshalBinary()
	if err != nil {
		log.Println("ERR: While serializing response for redis")
		log.Println(err)
		// Sending error doesn't make sense if caching does not work
		// App will still work using thirdparty APIs
	} else {
		// If data has been converted to byte array
		transactionStatus := c.rdsClient.Set(ctx, location, resBytes, time.Hour)
		if transactionStatus.Err() != nil {
			log.Println("ERR: While caching")
			log.Println(err)
			// Again failing to cache shouldn't stop flow
			// So no return
		}
	}

	log.Printf("Exec: %s", time.Since(start))
	utils.WriteJSON(w, 200, response.Data)
}
