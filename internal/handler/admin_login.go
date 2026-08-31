package handler

import (
	"ladno-teams/internal/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminLoginPage(c *gin.Context) {
	if user, ok := h.authenticateFromRequest(c); ok && user.IsAdmin {
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	h.renderAdminLoginPage(c, views.AuthPage{
		ErrorMessage: adminLoginErrorMessage(c.Query("error")),
	})
}

func (h *Handler) AdminLoginSubmit(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	if login == "" || password == "" {
		h.renderAdminLoginPage(c, views.AuthPage{ErrorMessage: "Login and password are required"})
		return
	}

	user, apiErr := h.services.User.GetByCredentials(login, password)
	if apiErr != nil {
		h.renderAdminLoginPage(c, views.AuthPage{ErrorMessage: "Invalid login or password"})
		return
	}

	if user.IsBlocked {
		h.renderAdminLoginPage(c, views.AuthPage{ErrorMessage: "User is blocked"})
		return
	}

	if !user.IsAdmin {
		h.renderAdminLoginPage(c, views.AuthPage{ErrorMessage: "Admin access only"})
		return
	}

	tokens, apiErr := h.services.Auth.Auth(*user)
	if apiErr != nil {
		h.renderAdminLoginPage(c, views.AuthPage{ErrorMessage: "Failed to create session"})
		return
	}

	h.setAuthCookie(c, tokens.RefreshToken)
	c.Redirect(http.StatusFound, "/admin")
}

func (h *Handler) AdminLogout(c *gin.Context) {
	if token, err := c.Cookie("token"); err == nil {
		_ = h.services.Session.DestroyByToken(token)
	}
	h.clearAuthCookie(c)
	c.Redirect(http.StatusFound, adminLoginPath)
}

func adminLoginErrorMessage(code string) string {
	switch code {
	case "admin_required":
		return "Admin access only"
	case "blocked":
		return "User is blocked"
	default:
		return ""
	}
}
