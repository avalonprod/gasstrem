package models

import "errors"

var (
	ErrorUserNotFound            = errors.New("user doesn't exists")
	ErrorVerificationCodeInvalid = errors.New("verification code is invalid")
	ErrUserAlreadyExists         = errors.New("user with such email already exists")
	ErrUserNotCompleteData       = errors.New("All input fields are required")
	ErrUserSendEmail             = errors.New("error while sending email")
)
