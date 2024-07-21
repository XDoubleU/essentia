package httptools_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/pkg/httptools"
)

func TestCSV(t *testing.T) {
	res := httptest.NewRecorder()

	data := [][]string{
		{
			"h1",
			"h2",
		},
		{
			"c11",
			"c12",
		},
	}

	writeCSV := func(w http.ResponseWriter, _ *http.Request) {
		err := httptools.WriteCSV(w, "test", data)
		require.Nil(t, err)
	}
	http.HandlerFunc(writeCSV).ServeHTTP(res, nil)

	records, err := httptools.ReadCSV(res.Body)

	require.Nil(t, err)
	assert.Equal(t, data, records)
}

func TestJSON(t *testing.T) {
	res := httptest.NewRecorder()

	data := [][]string{
		{
			"h1",
			"h2",
		},
		{
			"c11",
			"c12",
		},
	}

	writeJSON := func(w http.ResponseWriter, _ *http.Request) {
		err := httptools.WriteJSON(w, http.StatusOK, data, nil)
		require.Nil(t, err)
	}
	http.HandlerFunc(writeJSON).ServeHTTP(res, nil)

	var result [][]string
	err := httptools.ReadJSON(res.Body, &result)

	require.Nil(t, err)
	assert.Equal(t, data, result)
}
