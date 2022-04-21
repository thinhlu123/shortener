package repository

import (
	"context"
	"github.com/thinhlu123/shortener/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LinkRepository struct{}

func NewLinkRepo() *LinkRepository {
	return &LinkRepository{}
}

func (LinkRepository) CreateShortLink(ctx context.Context, item models.Link) error {
	now := time.Now()
	item.CreatedTime = &now
	item.UpdatedTime = &now

	_, err := models.LinkDB.GetCollection().InsertOne(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (LinkRepository) GetLink(ctx context.Context, item models.Link) (string, error) {
	var link models.Link
	err := models.LinkDB.GetCollection().FindOne(ctx, item).Decode(&link)
	if err != nil {
		return "", err
	}

	return link.OriginalUrl, nil
}

func (LinkRepository) IncreaseClickCount(ctx context.Context, item models.Link) error {
	_, err := models.LinkDB.GetCollection().UpdateOne(ctx, item, bson.D{
		{"$inc", bson.D{{"click_count", 1}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}
