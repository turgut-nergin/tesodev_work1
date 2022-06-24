package lib

import "time"

func TimeStampNow() int64 {
	l, _ := time.LoadLocation("Europe/Istanbul")
	return time.Now().In(l).Unix()
}
