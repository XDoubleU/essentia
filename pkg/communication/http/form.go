package http

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

//nolint:gochecknoglobals //ok
var decoder = schema.NewDecoder()
var encoder = schema.NewEncoder()

func WriteForm(src any) (url.Values, error) {
	values := url.Values{}
	err := encoder.Encode(src, values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func ReadForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = decoder.Decode(dst, r.PostForm)
	if err != nil {
		return err
	}

	return nil
}
