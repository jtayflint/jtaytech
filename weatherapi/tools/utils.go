package tools

import (
	"strconv"
)

func NumFromString(val string) (float64, error) {
	// fmt.Printf("value: %s.", val))
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}
