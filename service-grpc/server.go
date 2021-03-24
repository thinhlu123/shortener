package service_grpc

import (
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"net"
	"shortener/model"
	"shortener/model/message"
	"strconv"
	"time"
)

var GrpcServer *apiMessageServer

type Handler func(c *gin.Context)

type apiMessageServer struct {
	message.UnimplementedAPIServiceServer
	Handlers map[string]Handler
}

func (a *apiMessageServer) Call(ctx context.Context, request *message.APIRequest) (*message.APIResponse, error) {

	//fullPath := request.GetMethod() + "://" + request.GetPath()
	//if a.Handlers[fullPath] != nil {
	//	a.Handlers[fullPath]
	//}

	return nil, nil
}

func (a *apiMessageServer) GetURL(ctx context.Context, request *message.APIRequest) (*message.APIResponse, error) {

	var url model.URL
	params := request.GetParams()
	if _, ok := params["url"]; !ok {
		return &message.APIResponse{
			Status:  message.Status_INVALID,
			Message: "Missing url",
			Headers: nil,
		}, nil
	}

	url.Url = params["url"]
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
		Upsert:         model.HelperPtrBool(true),
		ReturnDocument: &after,
	}
	result := model.UrlCollection.FindOneAndUpdate(context.TODO(), filter, update, &opts)
	if result.Err() != nil {
		return &message.APIResponse{
			Status:  message.Status_INVALID,
			Message: result.Err().Error(),
			Headers: nil,
		}, nil
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
			return &message.APIResponse{
				Status:  message.Status_INVALID,
				Message: resultUpdate.Err().Error(),
				Headers: nil,
			}, nil
		}
		_ = resultUpdate.Decode(&resultURL)

		//c.JSON(http.StatusOK, model.APIResponse{
		//	Status:  http.StatusOK,
		//	Message: "Successfully",
		//	Data:    c.Request.Host + "/s/" + resultURL.ShortedUrl,
		//})
		//

		return &message.APIResponse{
			Status:  message.Status_OK,
			Message: "Successfully",
			Headers: nil,
			Content: request.GetPath() + "/s/" + resultURL.ShortedUrl,
		}, nil
	}

	return &message.APIResponse{
		Status:  message.Status_OK,
		Message: "Successfully",
		Headers: nil,
		Content: request.GetPath() + "/s/" + resultUrlData.ShortedUrl,
	}, nil

	//c.JSON(http.StatusOK, model.APIResponse{
	//	Status:  http.StatusOK,
	//	Message: "Successfully",
	//	Data:    c.Request.Host + "/s/" + resultUrlData.ShortedUrl,
	//})
	//return nil, nil
}

func (a *apiMessageServer) Test(ctx context.Context, request *message.APIRequest) (*message.APIResponse, error) {
	return nil, nil
}

func newGrpcRequest() {

}

func newGrpcResponse() {

}

func encode(num int64) string {
	data := strconv.FormatInt(num, 10)
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func SetupgRPCServer(isSecure bool) {
	lis, _ := net.Listen("tcp", "8000")
	var opts []grpc.ServerOption
	if isSecure {
		// opts config
	}
	grpcServer := grpc.NewServer(opts...)
	message.RegisterAPIServiceServer(grpcServer, GrpcServer)
	_ = grpcServer.Serve(lis)
}
