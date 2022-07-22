package lib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func TimeStampNow() int64 {
	l, _ := time.LoadLocation("Europe/Istanbul")
	return time.Now().In(l).Unix()
}

func ValidatePaginator(limit, offset string, maxlimit int64) (int64, int64) {
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

func GetAcceptedSortField(sortValue string) string {
	if strings.EqualFold(sortValue, "subject") {
		return "subject"
	}

	if strings.EqualFold(sortValue, "body") {
		return "body"
	}

	if strings.EqualFold(sortValue, "status") {
		return "status"
	}

	if strings.EqualFold(sortValue, "lastAnsweredAt") {
		return "lastAnsweredAt"
	}
	return ""
}

func GetAcceptedSortDirection(direction string) int {
	if strings.EqualFold(direction, "asc") || strings.EqualFold(direction, "1") {
		fmt.Println(direction)

		return 1
	}
	if strings.EqualFold(direction, "desc") || strings.EqualFold(direction, "dsc") || strings.EqualFold(direction, "-1") {
		return -1
	}
	return 0
}
