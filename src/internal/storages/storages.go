package storages

import (
	"context"

	"github.com/avalonprod/gasstrem/src/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	Create(ctx context.Context, user models.User) error
	GetByCredentials(ctx context.Context, email, password string) (models.User, error)
	SetSession(ctx context.Context, userID primitive.ObjectID, session models.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error)
	IsDuplicateUserEmail(ctx context.Context, email string) bool
}

type Invoices interface {
	Create(ctx context.Context, invoce models.Invoice) error
}

type Storages struct {
	Users    Users
	Invoices Invoices
}

func NewStorages(db *mongo.Database) *Storages {
	return &Storages{
		Users:    NewUsersStorage(db),
		Invoices: NewInvoicesStorage(db),
	}
}
