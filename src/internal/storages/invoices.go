package storages

import (
	"context"

	"github.com/avalonprod/gasstrem/src/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoicesStorage struct {
	db *mongo.Collection
}

func NewInvoicesStorage(db *mongo.Database) *InvoicesStorage {
	return &InvoicesStorage{
		db: db.Collection(invoicesCollection),
	}
}

func (i *InvoicesStorage) Create(ctx context.Context, invoce models.Invoice) error {
	_, err := i.db.InsertOne(ctx, invoce)

	if err != nil {
		return err
	}

	return nil
}

func (i *InvoicesStorage) GetAllInvoceByUserId(ctx context.Context, userID primitive.ObjectID) ([]models.Invoice, error) {
	var invoces []models.Invoice

	res, err := i.db.Find(ctx, bson.M{"userID": userID})
	res.All(ctx, &invoces)
	if err != nil {
		return []models.Invoice{}, err
	}

	return invoces, nil
}
