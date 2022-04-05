package util

import (
	"strconv"
)

// IsInt is it int
func IsInt(str string) bool {
	_, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	return true
}
