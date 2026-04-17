package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/sirupsen/logrus"
)

func Response(data domain.HttpResponse, w http.ResponseWriter) {
	response, err := json.Marshal(data)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(`{"status":"error","message":"Internal Server Error"}`))
		if err != nil {
			return
		}
		return
	}

	go func() {
		logResponse, _ := json.Marshal(data)
		logrus.Info("Response: ", string(logResponse))
	}()

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(data.Code)

	_, err = w.Write(response)
	if err != nil {
		return
	}

}
