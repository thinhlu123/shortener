package db

import (
	"context"
	_var "github.com/thinhlu123/shortener/pkg/var"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	/* --- DB connection config --- */
	_maxPoolSize     uint64 = 200
	_minPoolSize     uint64 = 3
	_maxConnIdleTime        = time.Duration(600000000) // 600s
)

type Configuration struct {
	Address            []string
	Username           string
	Password           string
	AuthDB             string
	ReplicaSetName     string
	Ssl                bool
	SecondaryPreferred bool

	timezone *time.Location
}

type Client struct {
	client *mongo.Client
}

func (cl *Client) Connect(config Configuration) error {
	mongoClient := &options.ClientOptions{
		Auth: &options.Credential{
			AuthSource: config.AuthDB,
			Username:   config.Username,
			Password:   config.Password,
		},
		HeartbeatInterval: &_maxConnIdleTime,
		Hosts:             config.Address,
		LocalThreshold:    nil,
		MaxConnIdleTime:   &_maxConnIdleTime,
		MaxPoolSize:       &_maxPoolSize,
		MinPoolSize:       &_minPoolSize,
		ReplicaSet:        &config.ReplicaSetName,
		RetryReads:        &_var.Values.False,
		RetryWrites:       &_var.Values.False,
	}

	var err error
	cl.client, err = mongo.Connect(nil, mongoClient)
	if err != nil {
		return err
	}

	err = cl.client.Ping(nil, nil)
	if err != nil {
		return err
	}

	return err
}

func (cl *Client) Disconnect() {
	_ = cl.client.Disconnect(context.Background())
}

func (cl *Client) InitDatabase(dbName string, handler onConnectDBHandler) *Database {
	return &Database{
		Database:    cl.client.Database(dbName),
		DbName:      dbName,
		onConnectDb: handler,
	}
}

type onConnectDBHandler func(db *mongo.Database) error

type Database struct {
	Database *mongo.Database
	DbName   string

	onConnectDb onConnectDBHandler
}

func (db *Database) InitCollection() error {
	if err := db.onConnectDb(db.Database); err != nil {
		return err
	}

	return nil
}

type Collection struct {
	CollectionName string
	collection     *mongo.Collection
}

func (col *Collection) Init(colName string, db *mongo.Database) {
	col.CollectionName = colName
	col.collection = db.Collection(colName)
}

func (col *Collection) GetCollection() *mongo.Collection {
	return col.collection
}
