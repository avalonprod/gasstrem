package services

import (
	"context"
	"time"

	"github.com/avalonprod/gasstrem/src/internal/config"
	"github.com/avalonprod/gasstrem/src/internal/models"
	"github.com/avalonprod/gasstrem/src/internal/storages"
	"github.com/avalonprod/gasstrem/src/packages/auth"
	"github.com/avalonprod/gasstrem/src/packages/email"
	"github.com/avalonprod/gasstrem/src/packages/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users interface {
	UsersSignUp(ctx context.Context, input UserSignUpInput) error
	UsersSignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Invoices interface {
	CreateInvoice(ctx context.Context, input InvoiceInput) error
	GetAllInvoceByUserId(ctx context.Context, userID primitive.ObjectID) ([]models.Invoice, error)
}

type Emails interface {
	SendUserVerificationEmail(VerificationEmailInput) error
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Services struct {
	Users    Users
	Invoices Invoices
}

type Options struct {
	Storages        *storages.Storages
	EmailSender     email.Sender
	EmailConfig     config.EmailConfig
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(options *Options) *Services {
	emailsService := NewEmailService(options.EmailSender, options.EmailConfig)
	usersService := NewUsersService(options.Storages.Users, options.Hasher, emailsService, options.TokenManager, options.AccessTokenTTL, options.RefreshTokenTTL)
	return &Services{
		Users:    usersService,
		Invoices: NewInvoicesService(options.Storages.Invoices),
	}
}
