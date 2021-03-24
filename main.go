package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"shortener/api"
	"shortener/model"
	service_grpc "shortener/service-grpc"
)

func main() {
	clientDB, err := model.ConnectDB()
	if err != nil {
		panic(err)
	}

	defer func() {
		// Close the connection once no longer needed
		err = clientDB.Disconnect(context.TODO())

		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Connection to MongoDB closed.")
		}
	}()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	routerAPIGroup := router.Group("/api")
	//routerAPIGroup.Use(func(c *ginContext) {
	//
	//})
	routerAPIGroup.GET("/getURL", api.GetURL)
	router.GET("/test", api.Test)
	router.GET("/s/:key", api.RedirectURL)

	//setHandlerGrpc(router.Routes())
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	_ = router.Run()
	// router.Run(":3000") for a hard coded port


	// ========== gRPC server setup ============
	service_grpc.SetupgRPCClient("localhost:8000", false)
	service_grpc.SetupgRPCServer(false)
}

func setHandler(router *gin.Engine, method, path string, handlers ...gin.HandlerFunc) {
	router.Handle(method, path, handlers...)
	//service_grpc.GrpcServer.Handlers[method + "://" + path] =
}

func setHandlerGrpc(routes gin.RoutesInfo) {
	//for _, route := range routes {
		//service_grpc.GrpcServer.Handlers[route.Method + "://" + route.Path] = route.

	//}
}