package lib

import (
	"strconv"
	"time"
)

func TimeStampNow() int64 {
	l, _ := time.LoadLocation("Europe/Istanbul")
	return time.Now().In(l).Unix()
}

func ValidatePaginator(limit, offset string, maxlimit int) (int64, int64) {
	limInt, err := strconv.ParseInt(limit, 10, 64)
	max := int64(maxlimit)
	if err != nil || limInt <= 0 {
		limInt = 25
	}
	if limInt > max {
		limInt = max
	}
	offsetInt, err := strconv.ParseInt(offset, 10, 64)
	if err != nil || offsetInt < 0 {
		offsetInt = 0
	}

	return limInt, offsetInt
}
