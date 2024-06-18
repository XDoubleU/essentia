package tools_test

import (
	"testing"
	"time"

	"github.com/XDoubleU/essentia/pkg/tools"
	"github.com/stretchr/testify/assert"
)

func TestStartOfDay(t *testing.T) {
	now := time.Now()
	startOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	assert.Equal(t, startOfDay, tools.StartOfDay(now))
}

func TestEndOfDay(t *testing.T) {
	now := time.Now()
	endOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		23, 59, 59, 999999999,
		now.Location(),
	)

	assert.Equal(t, endOfDay, tools.EndOfDay(now))
}

func TestTimeZoneIndependentTimeNow(t *testing.T) {
	now := time.Now()
	utcTimeZone, _ := time.LoadLocation("UTC")
	result := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		utcTimeZone,
	)

	assert.Equal(t, result, tools.TimeZoneIndependentTimeNow(now.Location().String()))
}
