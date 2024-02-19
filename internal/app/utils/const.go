package utils

import (
	"fmt"
	"time"
)

var TimeLocation *time.Location
var RFC3339 = "2006-01-02T15:04:05.000Z07:00"

func init() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	TimeLocation = loc
	fmt.Printf("use time location of %+v (now:%s)\n", loc, time.Now().Format(RFC3339)) // TODO Log this
}
