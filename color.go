package main

import (
	"math/rand"
	"strconv"
)

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func RandomRGB() RGB {
	return RGB{
		Red:   uint8(rand.Intn(255)),
		Green: uint8(rand.Intn(255)),
		Blue:  uint8(rand.Intn(255)),
	}
}

func HexToRGB(hex string) (RGB, error) {
	var rgb RGB
	values, err := strconv.ParseUint(string(hex), 16, 32)
	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   uint8(values >> 16),
		Green: uint8((values >> 8) & 0xFF),
		Blue:  uint8(values & 0xFF),
	}

	return rgb, nil
}
