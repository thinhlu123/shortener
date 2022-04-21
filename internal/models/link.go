package models

import (
	"github.com/thinhlu123/shortener/pkg/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var LinkDB db.Collection

type Link struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedTime *time.Time         `json:"created_time,omitempty" bson:"c_t,omitempty"`
	UpdatedTime *time.Time         `json:"updated_time,omitempty" bson:"u_t,omitempty"`

	OriginalUrl string `json:"original_url,omitempty" bson:"original_url,omitempty"`
	ShortUrl    string `json:"short_url,omitempty" bson:"short_url,omitempty"`

	Username   string `json:"usr,omitempty" bson:"usr,omitempty"`
	ClickCount int64  `json:"click_count,omitempty" bson:"click_count,omitempty"`
}
