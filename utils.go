package valuator

import (
	"math"
)

func round(val float64) float64 {
	return math.Floor(val*100) / 100
}

func percentage(val float64) float64 {
	return math.Floor(val * 100)
}

func yoyCalc(past float64, curr float64, pc bool) float64 {
	if past == 0 {
		return 0
	}
	if pc && past != 0 {
		return percentage((curr - past) / past)
	}
	return round(curr - past)
}

func avgCalc(total float64, length float64) float64 {
	return math.Floor(total*100/length) / 100
}

func contains(key int, db []int) bool {
	for _, d := range db {
		if d == key {
			return true
		}
	}
	return false
}
