package internal

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

/* func parseHourWeather(weatherData interface{}) (parsedWeather *WeatherForHour, err error) {
	parsedWeather, ok := weatherData.(*WeatherForHour)
	if !ok {
		return nil, errors.New("couldn't assert weather by hour data to struct")
	}
	return parsedWeather, nil
}

func parseDayWeather(weatherData interface{}) (parsedWeather *Weather, err error) {
	parsedWeather, ok := weatherData.(*Weather)
	if !ok {
		return nil, errors.New("couldn't assert daily weather data to struct")
	}
	return parsedWeather, nil
} */
