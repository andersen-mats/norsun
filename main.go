package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)
	
func main() {
	var d uint8
	if len(os.Args) >= 2 {
		var parsed uint64
		parsed, err := strconv.ParseUint(os.Args[1], 10, 64)
		if err != nil {
			fmt.Printf("Error parsing argument: %v\n", err)
			os.Exit(1)
		}
		d = uint8(parsed)
	}

	const apiURL string = "https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=60.10&lon=9.58"
	const userAgent string = "for educational use / github.com/andersen-mats/norsun"
	client := &http.Client{}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error performing request: %v\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Error | Status code: %d\n", res.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		os.Exit(1)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Printf("Error parsing JSON data: %v\n", err)
		os.Exit(1)
	}

	// curTemp := weather.Properties.Timeseries[0].Data.Instant.Details.AirTemperature
	// curWindSpeed := weather.Properties.Timeseries[0].Data.Instant.Details.WindSpeed
	// curSymbol := weather.Properties.Timeseries[0].Data.Next12Hours.Summary.SymbolCode

	// curStatus := fmt.Sprintf("Now: %s, %.0fC, %.0fm/s", curSymbol, curTemp, curWindSpeed)
	// fmt.Println(curStatus)

	var i uint8
	days := weather.Properties.Timeseries
	var prev string 

	for _, day := range days {
		date, err := time.Parse(time.RFC3339, day.Time)
		if err != nil {
			fmt.Printf("Error parsing time: %v\n", err)
		}
		dateStr := date.Format("01/02")
		if dateStr == prev {
			prev = dateStr
			continue
		}

		if date.Before(time.Now()) {
			dateStr = "Now: "
		}
		prev = dateStr


		dayTemp := day.Data.Instant.Details.AirTemperature
		var dayTempStr string
		if dayTemp <= 0 {
			dayTempStr = color.CyanString(fmt.Sprintf("%.0fC", dayTemp))
		} else if dayTemp >= 20 {
			dayTempStr = color.RedString(fmt.Sprintf("%.0f", dayTemp))
		}
		
		message := fmt.Sprintf("%s %s, %.0fm/s - %s",
			dateStr,
			dayTempStr,
			day.Data.Instant.Details.WindSpeed,
			strings.Split(day.Data.Next12Hours.Summary.SymbolCode, "_")[0],
		)

		fmt.Println(message)

		i++
		if i == d {
			break
		}
	}
}
