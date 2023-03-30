package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 10 * time.Second

func NewClient(url, username, password string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	opts := options.Client().ApplyURI(url)

	if username != "" && password != "" {
		opts.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	client, err := mongo.NewClient(opts)

	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
