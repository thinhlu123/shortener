package main

import (
	"context"
	"github.com/thinhlu123/shortener/config"
	"github.com/thinhlu123/shortener/internal/models"
	"github.com/thinhlu123/shortener/internal/server"
	cache2 "github.com/thinhlu123/shortener/pkg/cache"
	"github.com/thinhlu123/shortener/pkg/db"
	"github.com/thinhlu123/shortener/pkg/logger"
	_var "github.com/thinhlu123/shortener/pkg/var"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	log.Println("Starting server")

	// get config
	configPath := config.GetConfigPath(os.Getenv("config"))
	if err := config.GetConfig(configPath); err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	// init logger
	appLogger := logger.NewApiLogger(config.Conf)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		config.Conf.Server.AppVersion,
		config.Conf.Log.Level,
		config.Conf.Server.Mode,
		config.Conf.Server.SSL,
	)
	appLogger.Infof("Success parsed config: %#v", config.Conf.Server.AppVersion)
	appLogger.Infof("%v", config.Conf)

	// init mongodb client
	var client db.Client
	if err := client.Connect(db.Configuration{
		Address:  config.Conf.DB.DBHost,
		Username: config.Conf.DB.DBUser,
		Password: config.Conf.DB.DBPass,
		AuthDB:   config.Conf.DB.DBAuth,
		Ssl:      config.Conf.DB.EnableSSL,
	}); err != nil {
		appLogger.Fatalf("Fail to connect db: %v", err)
	}
	defer client.Disconnect()

	database := client.InitDatabase("shortener", func(db *mongo.Database) error {
		models.LinkDB.Init("link", db)
		models.UserDB.Init("user", db)

		// create indexLink for user collection
		indexUser := []mongo.IndexModel{
			{
				Keys: bson.M{"usr": 1},
				Options: &options.IndexOptions{
					Unique: &_var.Values.True,
				},
			},
		}
		_, err := models.UserDB.GetCollection().Indexes().CreateMany(context.Background(), indexUser)
		if err != nil {
			return err
		}

		// create indexLink for link collection
		expireTime := int32(60 * 60 * 24 * 365 * 2)
		opts := options.IndexOptions{
			ExpireAfterSeconds: &expireTime,
		}
		indexLink := []mongo.IndexModel{
			{
				Keys:    bson.M{"updated_time": 1},
				Options: &opts,
			},
			{
				Keys: bson.M{"short_url": 1},
			},
			{
				Keys: bson.M{"user_id": 1},
			},
		}
		_, err = models.LinkDB.GetCollection().Indexes().CreateMany(context.Background(), indexLink)
		if err != nil {
			return err
		}

		return nil
	})
	if err := database.InitCollection(); err != nil {
		log.Fatal(err)
		return
	}
	appLogger.Infof("Succeed init DB %#v", config.Conf.DB.DBName)

	// init redis cache client
	var cache cache2.Cache
	cache.Init()
	appLogger.Infof("Succeed init redis cache.")

	// init and run server
	s := server.NewServer(appLogger, &cache)
	appLogger.Fatal(s.Run())
}
