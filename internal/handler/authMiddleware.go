package handler

import (
	"fmt"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	adminLoginPath      = "/admin/login"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	header := c.GetHeader(AuthorizationHeader)
	if header == "" {
		exception.HttpResponseException(c, exception.AuthError("Empty Authorization Header"))
		return
	}

	user, ok := h.authenticateFromBearer(c)
	if !ok {
		return
	}

	c.Set(CtxUser, *user)
}

func (h *Handler) adminMiddleware(c *gin.Context) {
	user, ok := h.authenticateAdmin(c)
	if !ok {
		return
	}
	c.Set(CtxUser, *user)
}

func (h *Handler) authenticateAdmin(c *gin.Context) (*dto.UserDto, bool) {
	user, ok := h.authenticateFromRequest(c)
	if !ok {
		h.abortAdminUnauthed(c)
		return nil, false
	}

	if !user.IsAdmin {
		if h.isBrowserRequest(c) {
			c.Redirect(http.StatusFound, adminLoginPath+"?error=admin_required")
			c.Abort()
			return nil, false
		}
		exception.HttpResponseException(c, exception.Forbidden("admin access required"))
		return nil, false
	}

	return user, true
}

func (h *Handler) authenticateFromRequest(c *gin.Context) (*dto.UserDto, bool) {
	if user, ok := h.authenticateFromBearer(c); ok {
		return user, true
	}

	token, err := c.Cookie("token")
	if err != nil {
		return nil, false
	}

	parsedTokenPayload, err := h.services.Auth.ParseToken(token)
	if err != nil {
		return nil, false
	}

	if parsedTokenPayload.User.IsBlocked {
		if h.isBrowserRequest(c) {
			c.Redirect(http.StatusFound, adminLoginPath+"?error=blocked")
			c.Abort()
			return nil, false
		}
		exception.HttpResponseException(c, exception.Forbidden("user is blocked"))
		return nil, false
	}

	return &parsedTokenPayload.User, true
}

func (h *Handler) authenticateFromBearer(c *gin.Context) (*dto.UserDto, bool) {
	header := c.GetHeader(AuthorizationHeader)
	if header == "" {
		return nil, false
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" || len(headerParts[1]) == 0 {
		exception.HttpResponseException(c, exception.AuthError("Invalid Authorization Header"))
		return nil, false
	}

	parsedTokenPayload, err := h.services.Auth.ParseToken(headerParts[1])
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError(fmt.Sprintf("Invalid token: %s", err.Error())))
		return nil, false
	}

	if parsedTokenPayload.User.IsBlocked {
		exception.HttpResponseException(c, exception.Forbidden("user is blocked"))
		return nil, false
	}

	return &parsedTokenPayload.User, true
}

func (h *Handler) abortAdminUnauthed(c *gin.Context) {
	if c.GetHeader("HX-Request") != "" {
		c.Header("HX-Redirect", adminLoginPath)
		c.Status(http.StatusUnauthorized)
		c.Abort()
		return
	}

	if h.isBrowserRequest(c) {
		c.Redirect(http.StatusFound, adminLoginPath)
		c.Abort()
		return
	}

	exception.HttpResponseException(c, exception.AuthError("Unauthed"))
}

func (h *Handler) isBrowserRequest(c *gin.Context) bool {
	if c.GetHeader("HX-Request") != "" {
		return false
	}
	accept := c.GetHeader("Accept")
	return strings.Contains(accept, "text/html") || accept == "" || accept == "*/*"
}

func (h *Handler) setAuthCookie(c *gin.Context, refreshToken string) {
	cookie := h.services.Auth.AuthCookie(refreshToken)
	c.SetCookie(
		cookie.Name,
		cookie.Value,
		cookie.MaxAge,
		cookie.Path,
		cookie.Domain,
		cookie.Secure,
		cookie.HttpOnly,
	)
}

func (h *Handler) clearAuthCookie(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
}
