package repository

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/thinhlu123/shortener/config"
	"github.com/thinhlu123/shortener/internal/models"
	"github.com/thinhlu123/shortener/pkg/db"
	_var "github.com/thinhlu123/shortener/pkg/var"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

func beforeTest() {
	if err := config.GetConfig("../../config/config_test"); err != nil {
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

func TestLinkRepository_CreateAndGetShortLink(t *testing.T) {
	beforeTest()
	t.Parallel()

	linksRepository := NewLinkRepo()
	mockLink := models.Link{
		OriginalUrl: "http://localhost/test",
		ShortUrl:    "http://localhost/aXnsjd",
		Username:    "test",
		ClickCount:  0,
	}
	err := linksRepository.CreateShortLink(context.Background(), mockLink)
	require.NoError(t, err)

	shortLink, err := linksRepository.GetLink(context.Background(), models.Link{
		OriginalUrl: "http://localhost/test",
	})
	require.NoError(t, err)
	require.Equal(t, mockLink.ShortUrl, shortLink)
}

func TestLinkRepository_IncreaseClickCount(t *testing.T) {
	beforeTest()
	t.Parallel()

	linksRepository := NewLinkRepo()
	mockLink := models.Link{
		OriginalUrl: "http://localhost/test2",
		ShortUrl:    "http://localhost/aXnsjd12",
		Username:    "test",
		ClickCount:  0,
	}
	err := linksRepository.CreateShortLink(context.Background(), mockLink)
	require.NoError(t, err)

	filter := models.Link{
		OriginalUrl: "http://localhost/test2",
	}
	err = linksRepository.IncreaseClickCount(context.Background(), filter)
	require.NoError(t, err)

	var link models.Link
	err = models.LinkDB.GetCollection().FindOne(context.Background(), filter).Decode(&link)
	require.NoError(t, err)
	require.NotNil(t, link)
	require.Equal(t, link.ClickCount, 1)
}
