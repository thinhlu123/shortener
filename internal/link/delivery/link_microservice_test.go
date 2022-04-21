package delivery

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thinhlu123/shortener/config"
	"github.com/thinhlu123/shortener/internal/link/mocks"
	pd "github.com/thinhlu123/shortener/internal/link/pb"
	"github.com/thinhlu123/shortener/internal/models"
	cache2 "github.com/thinhlu123/shortener/pkg/cache"
	"github.com/thinhlu123/shortener/pkg/db"
	"github.com/thinhlu123/shortener/pkg/logger"
	_var "github.com/thinhlu123/shortener/pkg/var"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/metadata"
	"log"
	"testing"
)

func beforeTest() {
	if err := config.GetConfig("../../../config/config_test"); err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	var client db.Client
	if err := client.Connect(db.Configuration{
		Address:  config.Conf.DB.DBHost,
		Username: config.Conf.DB.DBUser,
		Password: config.Conf.DB.DBPass,
		AuthDB:   config.Conf.DB.DBAuth,
		Ssl:      config.Conf.DB.EnableSSL,
	}); err != nil {
		return
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
}

func TestLinkMicroservice_CreateLink(t *testing.T) {
	beforeTest()

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//mockRepo := mocks.NewMockLinkRepository(ctrl)
	mockUsecase := mocks.NewMockLinkUsecases(ctrl)

	appLogger := logger.NewApiLogger(config.Conf)
	appLogger.InitLogger()
	var cache cache2.Cache
	cache.Init()

	linkService := NewLinkMicroservice(appLogger, mockUsecase)

	md := metadata.New(map[string]string{
		"Authorization": "test",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	mockUsecase.EXPECT().CreateShortLink(gomock.Any(), gomock.Any()).Return("", nil)
	//mockRepo.EXPECT().CreateShortLink(gomock.Any(), gomock.Any()).Return(nil)

	link, err := linkService.CreateLink(ctx, &pd.CreateLinkReq{
		Link: "http://localhost/avaaxvasf",
	})
	require.NoError(t, err)
	require.NotEmpty(t, link)
}

func TestLinkMicroservice_GetLink(t *testing.T) {

}
