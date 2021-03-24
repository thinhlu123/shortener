package api

import (
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"shortener/model"
	"strconv"
	"time"
)

func GetURLgRPC(c *gin.Context) {

}

func GetURL(c *gin.Context) {
	var url model.URL
	err := c.BindJSON(&url)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if url.Url == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Message: "Missing Url",
			Data:    nil,
		})
		return
	}

	url.CreatedTime = time.Now()

	urlBson, _ := model.ToBsonDoc(url)
	filter, _ := model.ToBsonDoc(model.URL{
		Url: url.Url,
	})
	update := bson.M{
		"$setOnInsert": urlBson,
	}

	after := options.After
	opts := options.FindOneAndUpdateOptions{
		Upsert: model.HelperPtrBool(true),
		ReturnDocument: &after,

	}
	result := model.UrlCollection.FindOneAndUpdate(context.TODO(), filter, update, &opts)
	if result.Err() != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Message: result.Err().Error(),
			Data:    nil,
		})
		return
	}
	var resultUrlData model.URL
	_ = result.Decode(&resultUrlData)

	if resultUrlData.ShortedUrl == "" {
		//CountDB
		var countModel *model.Count

		doc, _ := model.ToBsonDoc(model.Count{
			NumberUrl: 1,
		})
		countResult := model.CountCollection.FindOneAndUpdate(context.TODO(), bson.M{}, bson.M{"$inc": doc})
		_ = countResult.Decode(&countModel)
		countNum := countModel.NumberUrl + 1

		count := encode(countNum)

		var resultURL *model.URL
		updateBson, _ := model.ToBsonDoc(model.URL{
			ShortedUrl: count,
		})
		updateData := bson.M{
			"$set": updateBson,
		}
		filter, _ := model.ToBsonDoc(model.URL{
			ID: resultUrlData.ID,
		})

		after := options.After
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}
		resultUpdate := model.UrlCollection.FindOneAndUpdate(context.TODO(), filter, updateData, &opt)
		if resultUpdate.Err() != nil {
			c.JSON(http.StatusBadRequest, model.APIResponse{
				Message: resultUpdate.Err().Error(),
				Data:    nil,
			})
			return
		}
		err = resultUpdate.Decode(&resultURL)

		c.JSON(http.StatusOK, model.APIResponse{
			Message: "Successfully",
			Data:    c.Request.Host + "/s/" + resultURL.ShortedUrl,
		})
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Message: "Successfully",
		Data:    c.Request.Host + "/s/" + resultUrlData.ShortedUrl,
	})
}

func encode(num int64) string {
	data := strconv.FormatInt(num, 10)
	return base64.StdEncoding.EncodeToString([]byte(data))
}
