package utils

import (
	"os"
	"strconv"
)

func IsCI() bool {
	v, _ := strconv.ParseBool(os.Getenv("CI"))
	return v
}
