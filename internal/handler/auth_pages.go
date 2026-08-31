package handler

import (
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/handler/validation"
	"ladno-teams/internal/views"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) HomePage(c *gin.Context) {
	if user, ok := h.authenticateFromRequest(c); ok {
		h.renderAuthPage(c, "Home", "pages/auth/home_user", views.AuthHomePage{Login: user.Login})
		return
	}
	h.renderAuthPage(c, "Home", "pages/auth/home_guest", nil)
}

func (h *Handler) AuthLoginPage(c *gin.Context) {
	if _, ok := h.authenticateFromRequest(c); ok {
		c.Redirect(http.StatusFound, "/")
		return
	}
	h.renderAuthPage(c, "Sign in", "pages/auth/login", views.AuthPage{
		ErrorMessage: authFormErrorMessage(c.Query("error")),
	})
}

func (h *Handler) AuthLoginSubmit(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	if errMsg := validateAuthForm(login, password); errMsg != "" {
		h.renderAuthPage(c, "Sign in", "pages/auth/login", views.AuthPage{
			ErrorMessage: errMsg,
			Login:        login,
		})
		return
	}

	user, apiErr := h.services.User.GetByCredentials(login, password)
	if apiErr != nil {
		h.renderAuthPage(c, "Sign in", "pages/auth/login", views.AuthPage{
			ErrorMessage: "Invalid login or password",
			Login:        login,
		})
		return
	}

	if user.IsBlocked {
		h.renderAuthPage(c, "Sign in", "pages/auth/login", views.AuthPage{
			ErrorMessage: "User is blocked",
			Login:        login,
		})
		return
	}

	tokens, apiErr := h.services.Auth.Auth(*user)
	if apiErr != nil {
		h.renderAuthPage(c, "Sign in", "pages/auth/login", views.AuthPage{
			ErrorMessage: "Failed to create session",
			Login:        login,
		})
		return
	}

	h.setAuthCookie(c, tokens.RefreshToken)
	c.Redirect(http.StatusFound, "/")
}

func (h *Handler) AuthRegisterPage(c *gin.Context) {
	if _, ok := h.authenticateFromRequest(c); ok {
		c.Redirect(http.StatusFound, "/")
		return
	}

	guidStr := c.Query("guid")
	if guidStr == "" {
		guidStr = c.Query("invite")
	}

	if guidStr == "" {
		h.renderAuthPage(c, "Register", "pages/auth/register_required", nil)
		return
	}

	inviteGuid, err := uuid.Parse(guidStr)
	if err != nil {
		h.renderAuthPage(c, "Register", "pages/auth/register_invite", views.AuthRegisterPage{
			GUID:         guidStr,
			ErrorMessage: "Invalid invite link",
		})
		return
	}

	invite, apiErr := h.services.Invite.GetByGuid(inviteGuid)
	if apiErr != nil {
		h.renderAuthPage(c, "Register", "pages/auth/register_invite", views.AuthRegisterPage{
			GUID:         guidStr,
			ErrorMessage: "Invite not found",
		})
		return
	}

	isExpired := time.Now().UTC().After(invite.ExpiredAt)
	teamName := ""
	if team, teamErr := h.services.Team.GetByGuid(invite.TeamGuid); teamErr == nil {
		teamName = team.Name
	}

	h.renderAuthPage(c, "Register", "pages/auth/register_invite", views.AuthRegisterPage{
		GUID:         guidStr,
		TeamName:     teamName,
		ErrorMessage: authFormErrorMessage(c.Query("error")),
		IsExpired:    isExpired,
	})
}

func (h *Handler) AuthRegisterSubmit(c *gin.Context) {
	guidStr := c.PostForm("guid")
	login := c.PostForm("login")
	password := c.PostForm("password")
	passwordConfirm := c.PostForm("password_confirm")
	name := c.PostForm("name")

	inviteGuid, err := uuid.Parse(guidStr)
	if err != nil {
		h.renderAuthPage(c, "Register", "pages/auth/register_invite", views.AuthRegisterPage{
			ErrorMessage: "Invalid invite link",
		})
		return
	}

	if errMsg := validateAuthForm(login, password); errMsg != "" {
		h.renderAuthPage(c, "Register", "pages/auth/register_invite", views.AuthRegisterPage{
			GUID:         guidStr,
			ErrorMessage: errMsg,
		})
		return
	}
	if password != passwordConfirm {
		h.renderAuthPage(c, "Register", "pages/auth/register_invite", views.AuthRegisterPage{
			GUID:         guidStr,
			ErrorMessage: "Passwords do not match",
		})
		return
	}
	if name == "" {
		h.renderAuthPage(c, "Register", "pages/auth/register_invite", views.AuthRegisterPage{
			GUID:         guidStr,
			ErrorMessage: "Name is required",
		})
		return
	}

	user, apiErr := h.services.Invite.RegisterByInvite(dto.InviteRegisterRequestDto{
		Guid:     inviteGuid,
		Login:    login,
		Password: password,
		Name:     name,
	})
	if apiErr != nil {
		h.renderAuthPage(c, "Register", "pages/auth/register_invite", views.AuthRegisterPage{
			GUID:         guidStr,
			ErrorMessage: apiErr.Message,
		})
		return
	}

	tokens, apiErr := h.services.Auth.Auth(*user)
	if apiErr != nil {
		h.renderAuthPage(c, "Register", "pages/auth/register_invite", views.AuthRegisterPage{
			GUID:         guidStr,
			ErrorMessage: "Account created, but sign-in failed. Try signing in manually.",
		})
		return
	}

	h.setAuthCookie(c, tokens.RefreshToken)
	c.Redirect(http.StatusFound, "/")
}

func (h *Handler) AuthInviteRegisterPage(c *gin.Context) {
	guid := c.Query("guid")
	if guid == "" {
		guid = c.Query("invite")
	}
	c.Redirect(http.StatusFound, "/register?guid="+guid)
}

func (h *Handler) AuthInviteRegisterSubmit(c *gin.Context) {
	h.AuthRegisterSubmit(c)
}

func (h *Handler) AuthLogout(c *gin.Context) {
	if token, err := c.Cookie("token"); err == nil {
		_ = h.services.Session.DestroyByToken(token)
	}
	h.clearAuthCookie(c)
	c.Redirect(http.StatusFound, "/")
}

func authFormErrorMessage(code string) string {
	switch code {
	case "blocked":
		return "User is blocked"
	default:
		return ""
	}
}

func validateAuthForm(login, password string) string {
	if login == "" || password == "" {
		return "Login and password are required"
	}
	if !validation.LoginRegex.MatchString(login) {
		return "Login: 3–32 characters, letters, digits, _ and -"
	}
	if !validation.PasswordRegex.MatchString(password) {
		return "Password: 6 to 128 characters"
	}
	return ""
}
