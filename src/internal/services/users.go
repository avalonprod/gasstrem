package services

import (
	"context"
	"errors"
	"time"

	"github.com/avalonprod/gasstrem/src/internal/models"
	"github.com/avalonprod/gasstrem/src/internal/storages"
	"github.com/avalonprod/gasstrem/src/packages/auth"
	"github.com/avalonprod/gasstrem/src/packages/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSignUpInput struct {
	Name     string
	Surname  string
	Email    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type UsersService struct {
	storages        storages.Users
	hasher          hash.PasswordHasher
	emailsService   Emails
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(storages storages.Users, hasher hash.PasswordHasher, emialsService Emails, tokenManager auth.TokenManager, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *UsersService {
	return &UsersService{
		storages:        storages,
		hasher:          hasher,
		emailsService:   emialsService,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (u *UsersService) UsersSignUp(ctx context.Context, input UserSignUpInput) error {
	if input.Name == "" || input.Surname == "" || input.Email == "" || input.Password == "" {
		return models.ErrUserNotCompleteData
	}
	passwordHash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:           input.Name,
		Surname:        input.Surname,
		Email:          input.Email,
		Password:       passwordHash,
		RegisteredTime: time.Now(),
		LastVisitTime:  time.Now(),
		Verification:   false,
	}

	isDublicate := u.storages.IsDuplicateUserEmail(ctx, input.Email)

	if isDublicate {
		return models.ErrUserAlreadyExists
	}

	err = u.storages.Create(ctx, user)
	if err != nil {
		return err
	}

	err = u.emailsService.SendUserVerificationEmail(VerificationEmailInput{Name: input.Name, Email: input.Email, VerificationCode: "333"})
	if err != nil {
		return models.ErrUserSendEmail
	}
	return nil
}

func (u *UsersService) UsersSignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	if input.Email == "" || input.Password == "" {
		return Tokens{}, models.ErrUserNotCompleteData
	}
	passwordHash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := u.storages.GetByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, models.ErrorUserNotFound) {
			return Tokens{}, err
		}
		return Tokens{}, err
	}

	return u.createSession(ctx, user.ID)

}

func (u *UsersService) createSession(ctx context.Context, userID primitive.ObjectID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = u.tokenManager.NewJWT(userID.Hex(), u.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = u.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}
	session := models.Session{
		RefreshToken: res.RefreshToken,
		ExpiresTime:  time.Now().Add(u.refreshTokenTTL),
	}

	err = u.storages.SetSession(ctx, userID, session)
	return res, err
}

func (u *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	user, err := u.storages.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return u.createSession(ctx, user.ID)
}
