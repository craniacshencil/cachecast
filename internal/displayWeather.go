package internal

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type Weather struct {
	ResolvedAddress   string  `json:"resolvedAddress"`
	Latitude          float32 `json:"latitude"`
	Longitude         float32 `json:"longitude"`
	QueryCost         int     `json:"queryCost"`
	Tzoffset          int     `json:"tzoffset"`
	Timezone          string  `json:"timezone"`
	Days              []Day   `json:"days"`
	CurrentConditions *Hour   `json:"currentConditions,omitempty"`
}

type Day struct {
	Conditions string  `json:"conditions"`
	Datetime   string  `json:"datetime"`
	Temp       float32 `json:"temp"`
	Tempmax    float32 `json:"tempmax"`
	Tempmin    float32 `json:"tempmin"`
}

type Hour struct {
	Conditions string  `json:"conditions"`
	Datetime   string  `json:"datetime"`
	Temp       float32 `json:"temp"`
}

type Cache struct {
	Cachestatus string
	Reqtime     string
}
type DisplayData struct {
	WeatherData  Weather
	CacheData    Cache
	ErrorMessage string
}

func displayWeather(
	w http.ResponseWriter,
	data Weather,
	cacheStatus string,
	reqtime string,
	onlyErrorMessage error,
	statusCode int,
) {
	cacheData := Cache{
		Cachestatus: cacheStatus,
		Reqtime:     reqtime,
	}
	completeErrorMessage := ""
	if onlyErrorMessage != nil {
		completeErrorMessage = fmt.Sprintf("ERROR %d: %s", statusCode, onlyErrorMessage.Error())
	}
	displayData := DisplayData{
		WeatherData:  data,
		CacheData:    cacheData,
		ErrorMessage: completeErrorMessage,
	}

	t, err := template.ParseFiles("./web/index.html")
	if err != nil {
		log.Println("ERR: While creating template for HTML file")
		log.Println(err)
		return
	}

	err = t.Execute(w, displayData)
	if err != nil {
		log.Println("ERR: While executing template")
		log.Println(err)
		return
	}
}
