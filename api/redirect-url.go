package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"shortener/model"
)

func RedirectURL(c *gin.Context) {
	key := c.Param("key")
	filter, _ := model.ToBsonDoc(model.URL{
		ShortedUrl:  key,
	})

	result := model.UrlCollection.FindOne(context.TODO(), filter)
	var redirectURL model.URL
	err := result.Decode(&redirectURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Message:   err.Error(),
		})
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Message:   "OK",
		Data:      redirectURL.Url,
	})
}