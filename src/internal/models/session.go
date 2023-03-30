package models

import "time"

type Session struct {
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
	ExpiresTime  time.Time `json:"expiresTime" bson:"expiresTime"`
}
