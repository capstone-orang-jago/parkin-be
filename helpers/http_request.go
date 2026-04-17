package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DecodeJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(v); err != nil {
		if err == io.EOF {
			return errors.New("request body is empty")
		}
		return err
	}

	return nil
}
