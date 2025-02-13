// Package grapher provides a tool for easily creating graphs.
package grapher

import (
	"fmt"
	"slices"
	"strconv"
	"time"
)

// GraphType is used to define the type of graph you need.
type GraphType int

// Numeric is a union of all numeric types.
type Numeric interface {
	int | int64 | float64
}

const (
	// Normal GraphType provides a grapher where you can set points a certain dates.
	Normal GraphType = iota
	// Cumulative GraphType will accumulate all previous values on your graph.
	Cumulative GraphType = iota
	// CumulativeSameDate GraphType will accumulate values
	// with the same date on your graph.
	CumulativeSameDate GraphType = iota
)

// Grapher is used to easily create graphs.
type Grapher[T Numeric] struct {
	graphType   GraphType
	dateFormat  string
	dateStrings []string
	values      map[string][]T
}

// New returns a new Grapher.
func New[T Numeric](graphType GraphType, dateFormat string) *Grapher[T] {
	return &Grapher[T]{
		graphType:   graphType,
		dateFormat:  dateFormat,
		dateStrings: []string{},
		values:      make(map[string][]T),
	}
}

// AddPoint adds a new point to the graph.
func (grapher *Grapher[T]) AddPoint(date time.Time, value T, label string) {
	dateStr := date.Format(grapher.dateFormat)

	dateIndex := grapher.getDateIndex(dateStr, label)

	grapher.updateDays(dateIndex, value, label)
}

func (grapher *Grapher[T]) getDateIndex(dateStr string, label string) int {
	dateIndex := slices.Index(grapher.dateStrings, dateStr)

	if dateIndex == -1 {
		grapher.addDays(dateStr)
		dateIndex = slices.Index(grapher.dateStrings, dateStr)
	}

	if _, ok := grapher.values[label]; !ok {
		grapher.values[label] = append(
			grapher.values[label],
			*new(T),
		)
	}

	return dateIndex
}

func (grapher *Grapher[T]) addDays(dateStr string) {
	if len(grapher.dateStrings) == 0 {
		grapher.dateStrings = append(grapher.dateStrings, dateStr)
		return
	}

	dateDay, _ := time.Parse(grapher.dateFormat, dateStr)
	smallestDate, _ := time.Parse(grapher.dateFormat, grapher.dateStrings[0])
	largestDate, _ := time.Parse(
		grapher.dateFormat,
		grapher.dateStrings[len(grapher.dateStrings)-1],
	)

	i := smallestDate
	for i.After(dateDay) {
		i = i.AddDate(0, 0, -1)

		grapher.dateStrings = append(
			[]string{i.Format(grapher.dateFormat)},
			grapher.dateStrings...)

		for label := range grapher.values {
			grapher.values[label] = append(
				[]T{*new(T)},
				grapher.values[label]...)
		}
	}

	i = largestDate
	for i.Before(dateDay) {
		i = i.AddDate(0, 0, 1)

		grapher.dateStrings = append(
			grapher.dateStrings,
			i.Format(grapher.dateFormat),
		)

		indexOfI := slices.Index(
			grapher.dateStrings,
			i.Format(grapher.dateFormat),
		)

		for label := range grapher.values {
			grapher.values[label] = append(
				grapher.values[label],
				grapher.values[label][indexOfI-1],
			)
		}
	}
}

func (grapher *Grapher[T]) updateDays(dateIndex int, value T, label string) {
	switch grapher.graphType {
	case Normal:
		grapher.values[label][dateIndex] = value
	case CumulativeSameDate:
		grapher.values[label][dateIndex] += value
	case Cumulative:
		for i := dateIndex; i < len(grapher.dateStrings); i++ {
			grapher.values[label][i] += value
		}
	}
}

// ToStringSlices returns the graph as a string slice of dates and values.
func (grapher Grapher[T]) ToStringSlices() ([]string, map[string][]string) {
	strValues := make(map[string][]string)

	for label, values := range grapher.values {
		for _, value := range values {
			strValue := ""
			switch v := any(value).(type) {
			case int:
				strValue = strconv.Itoa(v)
			case int64:
				strValue = strconv.Itoa(int(v))
			case float64:
				strValue = fmt.Sprintf("%.2f", v)
			}

			strValues[label] = append(strValues[label], strValue)
		}
	}

	return grapher.dateStrings, strValues
}

// ToSlices returns the graph as a string slice of dates and a typed slice of values.
func (grapher Grapher[T]) ToSlices() ([]string, map[string][]T) {
	return grapher.dateStrings, grapher.values
}
