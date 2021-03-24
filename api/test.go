package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"shortener/model"
	"strconv"
)
import "shortener/service-grpc"
import "shortener/model/message"

func Test(c *gin.Context) {
	res := Forward(c)
	code, _ := strconv.Atoi(res.Status)
	c.JSON(code, res)
}

func Forward(c *gin.Context) *model.APIResponse {
	r := newForwardRequest(
		c.Request.Method,
		c.Request.URL.Path,
		getParams(c),
		getContentText(c),
		getHeaders(c),
	)

	res, _ := r.MakeRequest()

	fmt.Println(res)

	resp := &model.APIResponse{
		Status:    strconv.FormatInt(int64(res.GetStatus()), 10),
		Message:   res.GetMessage(),
		ErrorCode: res.GetErrorCode(),
	}
	_ = json.Unmarshal([]byte(res.GetContent()), &resp.Data)

	return resp
}

type ForwardRequest struct {
	Method  string            `json:"method" bson:"method"`
	Path    string            `json:"path" bson:"path"`
	Params  map[string]string `json:"params,omitempty" bson:"params,omitempty"`
	Headers map[string]string `json:"headers,headers" bson:"headers,omitempty"`
	Content string            `json:"content,omitempty" bson:"content,omitempty"`
}

func (f *ForwardRequest) MakeRequest() (*message.APIResponse, error){
	conn := service_grpc.PickConn(false)
	client := message.NewAPIServiceClient(conn.Conn)

	req := newgRPCResquest(f)

	res, err := client.Call(context.TODO(), req)

	fmt.Println(res)
	return res, err
}

func newgRPCResquest(req *ForwardRequest) *message.APIRequest {
	return &message.APIRequest{
		Path:    req.Path,
		Method:  req.Method,
		Content: req.Content,
		Params:  req.Params,
		Headers: req.Headers,
	}
}

func getParams(c *gin.Context) map[string]string {
	vals := c.Request.URL.Query()
	m := make(map[string]string)
	for key := range vals {
		m[key] = vals.Get(key)
	}
	return m
}

func getContentText(c *gin.Context) (body string) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}

	body = string(bodyBytes)

	return
}

func getHeaders(c *gin.Context) map[string]string {
	vals := c.Request.Header
	m := make(map[string]string)
	for key := range vals {
		m[key] = vals.Get(key)
	}
	return m
}

func newForwardRequest(method string, path string, params map[string]string, content string, headers map[string]string) *ForwardRequest {
	return &ForwardRequest{
		Method:  method,
		Path:    path,
		Params:  params,
		Content: content,
		Headers: headers,
	}
}

//var vals = req.context.Request().Header
//var m = make(map[string]string)
//for key := range vals {
//m[key] = vals.Get(key)
//}
