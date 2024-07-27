package http

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
)

// WriteCSV writes the provided data as a CSV file with
// the provided filename to a [http.ResponseWriter].
func WriteCSV(w http.ResponseWriter, filename string, data [][]string) error {
	w.Header().Set("content-type", "text/csv")
	w.Header().
		Set("content-disposition", fmt.Sprintf("attachment;filename=%s.csv", filename))

	csvWriter := csv.NewWriter(w)

	return csvWriter.WriteAll(data)
}

// ReadCSV reads the returned CSV file from a [http.Response.Body].
func ReadCSV(body io.Reader) ([][]string, error) {
	csvReader := csv.NewReader(body)
	return csvReader.ReadAll()
}
