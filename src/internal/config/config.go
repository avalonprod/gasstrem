package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	defaultHTTPPort           = "8000"
	defaultHTTPWriteTimeout   = 10 * time.Second
	defaultHTTPReadTimeout    = 10 * time.Second
	defaultHTTPMaxHeaderBytes = 1
	defaultAccessTokenTTL     = 15 * time.Minute
	defaultRefreshTokenTTL    = 24 * time.Hour * 30
)

type (
	Config struct {
		HTTP     HTTPConfig
		SMTP     SMTPConfig
		Mongo    MongoConfig
		CacheTTL time.Duration `mapstructure:"ttl"`
		Email    EmailConfig
		Auth     AuthConfig
	}

	MongoConfig struct {
		URL      string
		Database string
		Username string
		Password string
	}

	HTTPConfig struct {
		Host           string        `mapstructure:"host"`
		Port           string        `mapstructure:"port"`
		ReadTimeout    time.Duration `mapstructure:"readTimeout"`
		WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
	}

	AuthConfig struct {
		JWT          JWTConfig
		PasswordSalt string
	}
	EmailSubjects struct {
		Verification string `mapstructure:"verification"`
	}

	EmailConfig struct {
		Templates EmailTemplates
		Subjects  EmailSubjects
	}

	EmailTemplates struct {
		Verification string
	}

	SMTPConfig struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		From     string `mapstructure:"from"`
		Password string
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}
)

func Init(configsDir string) (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	SetDefault()
	if err := parseConfigFile(configsDir); err != nil {
		return nil, err
	}
	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	setFromEnv(&cfg)
	return &cfg, nil
}

func parseConfigFile(ConfigsDir string) error {
	viper.AddConfigPath(ConfigsDir)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("cashe.ttl", &cfg.CacheTTL); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("smtp", &cfg.SMTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("email.templates", &cfg.Email.Templates); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("email.subjects", &cfg.Email.Subjects); err != nil {
		return err
	}
	return nil
}

func setFromEnv(cfg *Config) {
	cfg.Mongo.URL = os.Getenv("MONGODB_URL")
	cfg.Mongo.Username = os.Getenv("MONGODB_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGODB_PASSWORD")
	cfg.Mongo.Database = os.Getenv("MONGODB_DATABASE")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")
	cfg.HTTP.Host = os.Getenv("HTTP_HOST")

	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
}

func SetDefault() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.maxHeaderBytes", defaultHTTPMaxHeaderBytes)
	viper.SetDefault("http.writeTimeout", defaultHTTPWriteTimeout)
	viper.SetDefault("http.readTimeout", defaultHTTPReadTimeout)
	viper.SetDefault("http.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("http.refreshTokenTTL", defaultRefreshTokenTTL)
}
