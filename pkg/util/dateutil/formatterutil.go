package dateutil

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"time"
)

func AddDateString(year, month, day int, date, datePattern string) (string, error) {
	if date == "" {
		return date, nil
	}
	endCreatedTime, err := time.Parse(datePattern, date)
	if err != nil {
		log.Error(fmt.Sprintf("Error converting string to date %s: ", datePattern), err)
		return "", err
	}

	return endCreatedTime.AddDate(year, month, day).Format(datePattern), nil
}
