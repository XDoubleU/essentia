package tools_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/tools"
)

type Random struct {
}

func ignoreError(value string, _ error) string {
	return value
}

func ignoreValue(_ string, err error) error {
	return err
}

func TestAnyToString(t *testing.T) {
	assert.Equal(t, "string", ignoreError(tools.AnyToString("string")))
	assert.Equal(t, "1", ignoreError(tools.AnyToString(1)))
	assert.Equal(t, "1", ignoreError(tools.AnyToString(int64(1))))
	assert.Error(
		t,
		errors.New("undefined type"),
		ignoreValue(tools.AnyToString(Random{})),
	)
}
