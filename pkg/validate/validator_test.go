package validate_test

import (
	"testing"

	"github.com/XDoubleU/essentia/pkg/validate"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	strVal   string
	intVal   int
	int64Val int64
	tzVal    string
}

func (ts TestStruct) Validate() *validate.Validator {
	v := validate.New()

	validate.Check(v, ts.strVal, validate.IsNotEmpty, "strVal")

	validate.Check(v, ts.intVal, validate.IsGreaterThanFunc(-1), "intVal")
	validate.Check(v, ts.intVal, validate.IsGreaterThanOrEqualFunc(1), "intVal")
	validate.Check(v, ts.intVal, validate.IsLesserThanFunc(4), "intVal")
	validate.Check(v, ts.intVal, validate.IsLesserThanOrEqualFunc(2), "intVal")

	validate.Check(v, ts.int64Val, validate.IsGreaterThanFunc(int64(-1)), "int64Val")
	validate.Check(v, ts.int64Val, validate.IsGreaterThanOrEqualFunc(int64(1)), "int64Val")
	validate.Check(v, ts.int64Val, validate.IsLesserThanFunc(int64(4)), "int64Val")
	validate.Check(v, ts.int64Val, validate.IsLesserThanOrEqualFunc(int64(2)), "int64Val")

	validate.Check(v, ts.tzVal, validate.IsValidTimeZone, "tzVal")

	return v
}

func TestAllOk(t *testing.T) {
	ts := TestStruct{
		strVal:   "hello",
		intVal:   1,
		int64Val: 1,
		tzVal:    "Europe/Brussels",
	}

	assert.True(t, ts.Validate().Valid())
}

func TestIsEmpty(t *testing.T) {
	ts := TestStruct{
		strVal:   "",
		intVal:   1,
		int64Val: 1,
		tzVal:    "Europe/Brussels",
	}

	errors := map[string]string{
		"strVal": "must be provided",
	}

	v := ts.Validate()

	assert.False(t, v.Valid())
	assert.Equal(t, errors, v.Errors)
}

func TestIsNotGT(t *testing.T) {
	ts := TestStruct{
		strVal:   "hello",
		intVal:   -1,
		int64Val: -1,
		tzVal:    "Europe/Brussels",
	}

	errors := map[string]string{
		"intVal":   "must be greater than -1",
		"int64Val": "must be greater than -1",
	}

	v := ts.Validate()

	assert.False(t, v.Valid())
	assert.Equal(t, errors, v.Errors)
}

func TestIsNotGTE(t *testing.T) {
	ts := TestStruct{
		strVal:   "hello",
		intVal:   0,
		int64Val: 0,
		tzVal:    "Europe/Brussels",
	}

	errors := map[string]string{
		"intVal":   "must be greater than or equal to 1",
		"int64Val": "must be greater than or equal to 1",
	}

	v := ts.Validate()

	assert.False(t, v.Valid())
	assert.Equal(t, errors, v.Errors)
}

func TestIsNotLT(t *testing.T) {
	ts := TestStruct{
		strVal:   "hello",
		intVal:   4,
		int64Val: 4,
		tzVal:    "Europe/Brussels",
	}

	errors := map[string]string{
		"intVal":   "must be lesser than 4",
		"int64Val": "must be lesser than 4",
	}

	v := ts.Validate()

	assert.False(t, v.Valid())
	assert.Equal(t, errors, v.Errors)
}

func TestIsNotLTE(t *testing.T) {
	ts := TestStruct{
		strVal:   "hello",
		intVal:   3,
		int64Val: 3,
		tzVal:    "Europe/Brussels",
	}

	errors := map[string]string{
		"intVal":   "must be lesser than or equal to 2",
		"int64Val": "must be lesser than or equal to 2",
	}

	v := ts.Validate()

	assert.False(t, v.Valid())
	assert.Equal(t, errors, v.Errors)
}

func TestIsNotValidTz(t *testing.T) {
	ts := TestStruct{
		strVal:   "hello",
		intVal:   1,
		int64Val: 1,
		tzVal:    "AAAAAAAAAAAAAAAAAAAAAA",
	}

	errors := map[string]string{
		"tzVal": "must be a valid IANA value",
	}

	v := ts.Validate()

	assert.False(t, v.Valid())
	assert.Equal(t, errors, v.Errors)
}
