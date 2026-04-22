package helpers

import (
	"encoding/json"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Response(c *gin.Context, data domain.HttpResponse) {

	if logResponse, err := json.Marshal(data); err == nil {
		logrus.Info("Response: ", string(logResponse))
	}

	c.JSON(data.Code, data)
}
