package main
type Weather struct {
	Properties struct {
		Timeseries []struct {
			Time string `json:"time"`
			Data struct {
				Instant struct {
					Details struct {
						AirTemperature float64 `json:"air_temperature"`
						WindSpeed      float64 `json:"wind_speed"`
					} `json:"details"`
				} `json:"instant"`
				Next12Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
				} `json:"next_12_hours"`
			} `json:"data"`
		} `json:"timeseries"`
	} `json:"properties"`
}
