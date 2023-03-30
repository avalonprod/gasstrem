package storages

import (
	"context"
	"errors"
	"time"

	"github.com/avalonprod/gasstrem/src/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersStorage struct {
	db *mongo.Collection
}

func NewUsersStorage(db *mongo.Database) *UsersStorage {
	return &UsersStorage{
		db: db.Collection(usersCollection),
	}
}

func (u *UsersStorage) Create(ctx context.Context, user models.User) error {
	_, err := u.db.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (u *UsersStorage) GetByCredentials(ctx context.Context, email, password string) (models.User, error) {
	var user models.User

	if err := u.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, models.ErrorUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (u *UsersStorage) SetSession(ctx context.Context, userID primitive.ObjectID, session models.Session) error {
	_, err := u.db.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"session": session, "lastVisitTime": time.Now()}})

	return err
}

func (u *UsersStorage) GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error) {
	var user models.User

	if err := u.db.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresTime":  bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, models.ErrorUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (u *UsersStorage) IsDuplicateUserEmail(ctx context.Context, email string) bool {
	var usr models.User
	if isDuplicate := u.db.FindOne(ctx, bson.M{"email": email}).Decode(&usr); isDuplicate != nil {
		return false
	}
	return true
}
