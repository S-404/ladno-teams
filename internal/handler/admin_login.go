package handler

import (
	"html"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminLoginPage(c *gin.Context) {
	if user, ok := h.authenticateFromRequest(c); ok && user.IsAdmin {
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	renderHTML(c, adminLoginPage(adminLoginErrorMessage(c.Query("error"))))
}

func (h *Handler) AdminLoginSubmit(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	if login == "" || password == "" {
		renderHTML(c, adminLoginPage("Укажите логин и пароль"))
		return
	}

	user, apiErr := h.services.User.GetByCredentials(login, password)
	if apiErr != nil {
		renderHTML(c, adminLoginPage("Неверный логин или пароль"))
		return
	}

	if user.IsBlocked {
		renderHTML(c, adminLoginPage("Пользователь заблокирован"))
		return
	}

	if !user.IsAdmin {
		renderHTML(c, adminLoginPage("Доступ только для администраторов"))
		return
	}

	tokens, apiErr := h.services.Auth.Auth(*user)
	if apiErr != nil {
		renderHTML(c, adminLoginPage("Не удалось создать сессию"))
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
		return "Доступ только для администраторов"
	case "blocked":
		return "Пользователь заблокирован"
	default:
		return ""
	}
}

func adminLoginPage(errorMessage string) string {
	errorBlock := ""
	if errorMessage != "" {
		errorBlock = `<div class="error">` + html.EscapeString(errorMessage) + `</div>`
	}

	return `<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Admin Login</title>
	<style>
		body { font-family: sans-serif; margin: 0; background: #f5f5f5; display: flex; min-height: 100vh; align-items: center; justify-content: center; }
		.card { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); width: 100%; max-width: 360px; }
		h1 { margin-top: 0; font-size: 1.5rem; }
		label { display: block; margin-bottom: 0.25rem; font-weight: 600; }
		input { width: 100%; padding: 0.5rem; margin-bottom: 1rem; box-sizing: border-box; }
		button { width: 100%; padding: 0.6rem; background: #333; color: white; border: none; border-radius: 4px; cursor: pointer; font-size: 1rem; }
		button:hover { background: #111; }
		.error { background: #fee; color: #900; padding: 0.75rem; border-radius: 4px; margin-bottom: 1rem; }
	</style>
</head>
<body>
	<div class="card">
		<h1>Admin Login</h1>
		` + errorBlock + `
		<form method="POST" action="/admin/login">
			<label for="login">Login</label>
			<input id="login" name="login" type="text" autocomplete="username" required>
			<label for="password">Password</label>
			<input id="password" name="password" type="password" autocomplete="current-password" required>
			<button type="submit">Войти</button>
		</form>
	</div>
</body>
</html>`
}
