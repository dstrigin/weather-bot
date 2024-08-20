package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Weather struct {
	Current     int8
	FeelsLike   int8
	Description string
}

func GetInfoByCoords(lat, lon string) (map[string]interface{}, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", lat, lon, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Can't get weather with API")
		return nil, err
	}
	defer resp.Body.Close()

	byteData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(byteData, &data)
	return data, err
}

func GetWeather(data map[string]interface{}) *Weather {
	var weather Weather

	if deg, ok := data["main"].(map[string]interface{}); ok {
		currKelvin := deg["temp"].(float64)
		weather.Current = KelvinToCelsius(currKelvin)

		feelsLikeKelvin := deg["feels_like"].(float64)
		weather.FeelsLike = KelvinToCelsius(feelsLikeKelvin)
	}

	if weatherInfo, ok := data["weather"].([]interface{}); ok {
		if len(weatherInfo) > 0 {
			if desc, ok := weatherInfo[0].(map[string]interface{}); ok {
				weather.Description = desc["description"].(string)
			}
		}
	}

	return &weather
}

func KelvinToCelsius(kelvinDegrees float64) int8 {
	celsiusDegrees := kelvinDegrees - 273
	return int8(celsiusDegrees)
}

func GetInfoForBot() string {
	godotenv.Load()

	kv, err := GetInfoByCoords("47.235714", "39.701504")
	if err != nil {
		log.Printf("API error: %v", err)
		return "Сейчас не удалось получить данные. Попробуйте еще раз через пару секунд . . ."
	}
	w := GetWeather(kv)

	weatherInfo := fmt.Sprintf("Сейчас на улице %v градусов, ощущается как %v", w.Current, w.FeelsLike)
	return weatherInfo
}
