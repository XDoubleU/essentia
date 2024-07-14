package tools

import "time"

// StartOfDay returns the start of day of the specified date.
func StartOfDay(dateTime time.Time) time.Time {
	output := time.Date(
		dateTime.Year(),
		dateTime.Month(),
		dateTime.Day(),
		0, 0, 0, 0,
		dateTime.Location(),
	)

	return output
}

// EndOfDay returns the end of day of the specified date.
func EndOfDay(dateTime time.Time) time.Time {
	output := time.Date(
		dateTime.Year(),
		dateTime.Month(),
		dateTime.Day(),
		23, 59, 59, 999999999,
		dateTime.Location(),
	)

	return output
}

// TimeZoneIndependentTimeNow returns the time in the
// provided time zone but changes the time zone to UTC.
func TimeZoneIndependentTimeNow(locationTimeZone string) time.Time {
	timeZone, _ := time.LoadLocation(locationTimeZone)
	utcTimeZone, _ := time.LoadLocation("UTC")
	now := time.Now().In(timeZone)
	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		utcTimeZone,
	)
}
