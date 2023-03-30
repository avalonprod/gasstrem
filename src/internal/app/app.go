package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avalonprod/gasstrem/src/internal/config"
	"github.com/avalonprod/gasstrem/src/internal/handlers"
	"github.com/avalonprod/gasstrem/src/internal/server"
	"github.com/avalonprod/gasstrem/src/internal/services"
	"github.com/avalonprod/gasstrem/src/internal/storages"
	"github.com/avalonprod/gasstrem/src/packages/auth"
	"github.com/avalonprod/gasstrem/src/packages/database/mongodb"
	"github.com/avalonprod/gasstrem/src/packages/email/smtp"
	"github.com/avalonprod/gasstrem/src/packages/hash"
	"github.com/avalonprod/gasstrem/src/packages/logger"
)

const configPath string = "configs"

func Run() {
	cfg, err := config.Init(configPath)

	if err != nil {
		logger.Error(err)
		return
	}
	mongoClient, err := mongodb.NewClient(cfg.Mongo.URL, cfg.Mongo.Username, cfg.Mongo.Password)

	mongodb := mongoClient.Database(cfg.Mongo.Database)

	if err != nil {
		logger.Error(err)

		return
	}
	// Packages init
	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Password, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		logger.Error(err)

		return
	}
	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)
	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	storages := storages.NewStorages(mongodb)
	services := services.NewServices(&services.Options{
		Storages:        storages,
		EmailSender:     emailSender,
		EmailConfig:     cfg.Email,
		Hasher:          hasher,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.JWT.RefreshTokenTTL,
	})
	handlers := handlers.NewHandlers(services, tokenManager)

	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server %s\n", err.Error())
		}
	}()
	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err.Error())
	}
}
