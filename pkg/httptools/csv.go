package httptools

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
)

func WriteCSV(w http.ResponseWriter, filename string, data [][]string) error {
	w.Header().Set("content-type", "text/csv")
	w.Header().
		Set("content-disposition", fmt.Sprintf("attachment;filename=%s.csv", filename))

	csvWriter := csv.NewWriter(w)

	return csvWriter.WriteAll(data)
}

func ReadCSV(body io.Reader) ([][]string, error) {
	csvReader := csv.NewReader(body)
	return csvReader.ReadAll()
}
