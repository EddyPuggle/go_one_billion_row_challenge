package main

import (
	"math"
)

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

func round(value float64) float64 {
	return math.Round(value*10.0) / 10.0
}
