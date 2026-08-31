package handler

import (
	"net/http"

	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IAuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RefreshAccessToken(c *gin.Context)
	InviteRegister(c *gin.Context)
	AcceptInvite(c *gin.Context)
}

type AuthHandler struct {
	BaseHandler
	services *service.Service
}

func NewAuthHandler(services *service.Service, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		BaseHandler: *NewBaseHandler(validate),
		services:    services,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.AuthRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	user, apiErr := h.services.User.GetByCredentials(req.Login, req.Password)
	if apiErr != nil {
		exception.HttpResponseException(c, exception.BadRequest("wrong credentials"))
		return
	}

	if user.IsBlocked {
		exception.HttpResponseException(c, exception.Forbidden("user is blocked"))
		return
	}

	generatedToken, apiErr := h.services.Auth.Auth(*user)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	cookie := h.services.Auth.AuthCookie(generatedToken.RefreshToken)
	c.SetCookie(
		cookie.Name,
		cookie.Value,
		cookie.MaxAge,
		cookie.Path,
		cookie.Domain,
		cookie.Secure,
		cookie.HttpOnly,
	)

	c.JSON(http.StatusOK, dto.AuthResponseDto{
		AccessToken: generatedToken.AccessToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	if apiErr := h.services.Session.DestroyByToken(token); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, dto.EmptyResponse{})
}

func (h *AuthHandler) RefreshAccessToken(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	generatedToken, apiErr := h.services.Auth.RefreshAccessToken(token)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	cookie := h.services.Auth.AuthCookie(generatedToken.RefreshToken)
	c.SetCookie(
		cookie.Name,
		cookie.Value,
		cookie.MaxAge,
		cookie.Path,
		cookie.Domain,
		cookie.Secure,
		cookie.HttpOnly,
	)

	c.JSON(http.StatusOK, dto.AuthResponseDto{
		AccessToken: generatedToken.AccessToken,
	})
}

func (h *AuthHandler) InviteRegister(c *gin.Context) {
	var req dto.InviteRegisterRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	user, apiErr := h.services.Invite.RegisterByInvite(req)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	generatedToken, apiErr := h.services.Auth.Auth(*user)
	if apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	cookie := h.services.Auth.AuthCookie(generatedToken.RefreshToken)
	c.SetCookie(
		cookie.Name,
		cookie.Value,
		cookie.MaxAge,
		cookie.Path,
		cookie.Domain,
		cookie.Secure,
		cookie.HttpOnly,
	)

	c.JSON(http.StatusCreated, dto.AuthResponseDto{
		AccessToken: generatedToken.AccessToken,
	})
}

func (h *AuthHandler) AcceptInvite(c *gin.Context) {
	var req dto.AcceptInviteRequestDto
	if err := h.ValidateRequestBody(c, &req); err != nil {
		exception.HttpResponseException(c, exception.RequestValidationError(err.Error()))
		return
	}

	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	if apiErr := h.services.Invite.Accept(req.Guid, user.Guid); apiErr != nil {
		exception.HttpResponseException(c, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.EmptyResponse{})
}
