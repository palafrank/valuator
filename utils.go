package valuator

import "math"

func round(val float64) float64 {
	return math.Floor(val*100) / 100
}

func percentage(val float64) float64 {
	return math.Floor(val * 100)
}

func yoyCalc(past float64, curr float64) float64 {
	if past == 0 {
		return 0
	}
	return percentage((curr - past) / past)
}

func avgCalc(total float64, length float64) float64 {
	return math.Floor(total*100/length) / 100
}
