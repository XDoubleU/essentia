package grapher_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/XDoubleU/essentia/pkg/grapher"
	"github.com/stretchr/testify/assert"
)

func TestGrapherCumulative(t *testing.T) {
	grapher := grapher.New[int](grapher.Cumulative, "2006-01-02")

	dateNow := time.Now().UTC()
	for i := 0; i < 10; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, i), 1, "data")
	}

	for i := 1; i < 10; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, -1*i), 1, "data")
	}

	dateSlice, valueSlice := grapher.ToStringSlices()

	assert.Equal(t, 19, len(dateSlice))
	assert.Equal(t, 19, len(valueSlice["data"]))

	for i := 0; i < 19; i++ {
		assert.Equal(
			t,
			time.Now().UTC().AddDate(0, 0, i-9).Format("2006-01-02"),
			dateSlice[i],
		)
		assert.Equal(t, fmt.Sprint(i+1), valueSlice["data"][i])
	}
}

func TestGrapherNormal(t *testing.T) {
	grapher := grapher.New[int](grapher.Normal, "2006-01-02")

	dateNow := time.Now().UTC()
	for i := 0; i < 10; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, i), i, "data")
	}

	dateSlice, valueSlice := grapher.ToStringSlices()

	assert.Equal(t, 10, len(dateSlice))
	assert.Equal(t, 10, len(valueSlice["data"]))

	for i := 0; i < 10; i++ {
		assert.Equal(
			t,
			time.Now().UTC().AddDate(0, 0, i).Format("2006-01-02"),
			dateSlice[i],
		)
		assert.Equal(t, fmt.Sprint(i), valueSlice["data"][i])
	}
}
