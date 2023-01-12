package helper

import (
	"strconv"

	"github.com/whitekid/goxp/fx"
)

func ParseIntDef(s string, defaultValue, minValue, maxValue int) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return fx.Min([]int{fx.Max([]int{value, minValue}), maxValue})
}
