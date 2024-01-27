package utils

import (
	"log"
	"strconv"
	"strings"
	"time"
)

// GetTimeAsString - get time as a string
func GetTimeAsString(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000-0700")
}

// Now - get now as a string
func Now() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000-0700")
}

// ParseDatetime - parse datetime string into datetime object
func ParseDatetime(datetime string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000-0700", datetime)
}

// CurrentHourMin - Get current hour and min
func CurrentHourMin() (hour int, min int) {
	loc, loadErr := time.LoadLocation("Etc/UTC")
	if loadErr != nil {
		log.Printf("[CurrentHourMin] Unable to load location")
		log.Print(loadErr)
		return 0, 0
	}

	nowHour, nowMin, _ := time.Now().In(loc).Clock()
	return nowHour, nowMin
}

// ParseTimeFieldForHourMin - Parse time field for hour and minute
func ParseTimeFieldForHourMin(timeField string) (hour int, minute int) {
	if timeField == "" {
		return 0, 0
	}

	// extract
	parts := strings.Split(timeField, ":")

	// get rid of Z
	cleanedMinute := strings.Replace(parts[1], "Z", "", 1)

	// convert
	newHour, newHourErr := strconv.Atoi(parts[0])
	newMin, newMinErr := strconv.Atoi(cleanedMinute)

	if newHourErr != nil || newMinErr != nil {
		log.Print("[ParseTimeFieldForHourMin] Error getting hour or min")
		log.Print(newHourErr)
		log.Print(newMinErr)
	}

	return newHour, newMin
}
