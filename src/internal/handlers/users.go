package handlers

import (
	"errors"
	"net/http"

	"github.com/avalonprod/gasstrem/src/internal/models"
	"github.com/avalonprod/gasstrem/src/internal/services"
	"github.com/gin-gonic/gin"
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

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type refreshTokenInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *Handlers) UsersSignUp(c *gin.Context) {
	var input UserSignUpInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.Users.UsersSignUp(c.Request.Context(), services.UserSignUpInput{
		Name:     input.Name,
		Surname:  input.Surname,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handlers) UsersSignIn(c *gin.Context) {
	var input UserSignInInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	res, err := h.services.Users.UsersSignIn(c.Request.Context(), services.UserSignInInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, models.ErrorUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})

}

func (h *Handlers) userRefreshToken(c *gin.Context) {
	var input refreshTokenInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Users.RefreshTokens(c.Request.Context(), input.Token)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
