package handler

import (
	"fmt"
	"html"
	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/handler/validation"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) HomePage(c *gin.Context) {
	if user, ok := h.authenticateFromRequest(c); ok {
		renderHTML(c, authHomePage(user))
		return
	}
	renderHTML(c, authGuestHomePage())
}

func (h *Handler) AuthLoginPage(c *gin.Context) {
	if _, ok := h.authenticateFromRequest(c); ok {
		c.Redirect(http.StatusFound, "/")
		return
	}
	renderHTML(c, authLoginPage(authFormErrorMessage(c.Query("error")), ""))
}

func (h *Handler) AuthLoginSubmit(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	if errMsg := validateAuthForm(login, password); errMsg != "" {
		renderHTML(c, authLoginPage(errMsg, login))
		return
	}

	user, apiErr := h.services.User.GetByCredentials(login, password)
	if apiErr != nil {
		renderHTML(c, authLoginPage("Неверный логин или пароль", login))
		return
	}

	if user.IsBlocked {
		renderHTML(c, authLoginPage("Пользователь заблокирован", login))
		return
	}

	tokens, apiErr := h.services.Auth.Auth(*user)
	if apiErr != nil {
		renderHTML(c, authLoginPage("Не удалось создать сессию", login))
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
		renderHTML(c, authRegisterInviteRequiredPage())
		return
	}

	inviteGuid, err := uuid.Parse(guidStr)
	if err != nil {
		renderHTML(c, authInviteRegisterPage(guidStr, "", "Некорректная ссылка приглашения", false))
		return
	}

	invite, apiErr := h.services.Invite.GetByGuid(inviteGuid)
	if apiErr != nil {
		renderHTML(c, authInviteRegisterPage(guidStr, "", "Приглашение не найдено", false))
		return
	}

	isExpired := time.Now().UTC().After(invite.ExpiredAt)
	teamName := ""
	if team, teamErr := h.services.Team.GetByGuid(invite.TeamGuid); teamErr == nil {
		teamName = team.Name
	}

	renderHTML(c, authInviteRegisterPage(guidStr, teamName, authFormErrorMessage(c.Query("error")), isExpired))
}

func (h *Handler) AuthRegisterSubmit(c *gin.Context) {
	guidStr := c.PostForm("guid")
	login := c.PostForm("login")
	password := c.PostForm("password")
	passwordConfirm := c.PostForm("password_confirm")
	name := c.PostForm("name")

	inviteGuid, err := uuid.Parse(guidStr)
	if err != nil {
		renderHTML(c, authInviteRegisterPage("", "", "Некорректная ссылка приглашения", false))
		return
	}

	if errMsg := validateAuthForm(login, password); errMsg != "" {
		renderHTML(c, authInviteRegisterPage(guidStr, "", errMsg, false))
		return
	}
	if password != passwordConfirm {
		renderHTML(c, authInviteRegisterPage(guidStr, "", "Пароли не совпадают", false))
		return
	}
	if name == "" {
		renderHTML(c, authInviteRegisterPage(guidStr, "", "Укажите имя", false))
		return
	}

	user, apiErr := h.services.Invite.RegisterByInvite(dto.InviteRegisterRequestDto{
		Guid:     inviteGuid,
		Login:    login,
		Password: password,
		Name:     name,
	})
	if apiErr != nil {
		renderHTML(c, authInviteRegisterPage(guidStr, "", apiErr.Message, false))
		return
	}

	tokens, apiErr := h.services.Auth.Auth(*user)
	if apiErr != nil {
		renderHTML(c, authInviteRegisterPage(guidStr, "", "Аккаунт создан, но не удалось войти. Попробуйте войти вручную.", false))
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
		return "Пользователь заблокирован"
	default:
		return ""
	}
}

func validateAuthForm(login, password string) string {
	if login == "" || password == "" {
		return "Укажите логин и пароль"
	}
	if !validation.LoginRegex.MatchString(login) {
		return "Логин: 3–32 символа, только латиница, цифры, _ и -"
	}
	if !validation.PasswordRegex.MatchString(password) {
		return "Пароль: от 6 до 128 символов"
	}
	return ""
}

func authPageLayout(title, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s — Ladno Teams</title>
	<style>
		* { box-sizing: border-box; }
		body { font-family: sans-serif; margin: 0; background: #f0f2f5; color: #1a1a1a; min-height: 100vh; }
		.header { background: #1e3a5f; color: #fff; padding: 1rem 1.5rem; display: flex; align-items: center; gap: 1.5rem; flex-wrap: wrap; }
		.header a { color: #fff; text-decoration: none; opacity: 0.9; }
		.header a:hover { opacity: 1; text-decoration: underline; }
		.header .brand { font-weight: 700; font-size: 1.1rem; margin-right: auto; }
		.main { display: flex; min-height: calc(100vh - 56px); align-items: center; justify-content: center; padding: 2rem 1rem; }
		.card { background: #fff; padding: 2rem; border-radius: 10px; box-shadow: 0 2px 12px rgba(0,0,0,0.08); width: 100%%; max-width: 400px; }
		.card-wide { max-width: 440px; }
		h1 { margin: 0 0 0.25rem; font-size: 1.5rem; }
		.subtitle { color: #666; margin: 0 0 1.5rem; font-size: 0.95rem; }
		label { display: block; margin-bottom: 0.25rem; font-weight: 600; font-size: 0.9rem; }
		input { width: 100%%; padding: 0.55rem 0.65rem; margin-bottom: 1rem; border: 1px solid #ccc; border-radius: 6px; font-size: 1rem; }
		input:focus { outline: none; border-color: #1e3a5f; box-shadow: 0 0 0 2px rgba(30,58,95,0.15); }
		button, .btn { display: inline-block; padding: 0.6rem 1.2rem; background: #1e3a5f; color: #fff; border: none; border-radius: 6px; cursor: pointer; font-size: 1rem; text-decoration: none; text-align: center; }
		button:hover, .btn:hover { background: #152a45; }
		button.full { width: 100%%; }
		.error { background: #fef2f2; color: #991b1b; padding: 0.75rem; border-radius: 6px; margin-bottom: 1rem; font-size: 0.9rem; }
		.info { background: #eff6ff; color: #1e40af; padding: 0.75rem; border-radius: 6px; margin-bottom: 1rem; font-size: 0.9rem; }
		.footer-link { margin-top: 1.25rem; text-align: center; font-size: 0.9rem; color: #555; }
		.footer-link a { color: #1e3a5f; }
		.hint { font-size: 0.8rem; color: #888; margin: -0.75rem 0 1rem; }
		.home-card { max-width: 520px; text-align: center; }
		.home-actions { display: flex; gap: 0.75rem; justify-content: center; flex-wrap: wrap; margin-top: 1.5rem; }
		.btn-outline { background: transparent; color: #1e3a5f; border: 1px solid #1e3a5f; }
		.btn-outline:hover { background: #f0f4f8; }
	</style>
</head>
<body>
	<header class="header">
		<a href="/" class="brand">Ladno Teams</a>
		<a href="/login">Вход</a>
	</header>
	<main class="main">%s</main>
</body>
</html>`, html.EscapeString(title), body)
}

func authLoginPage(errorMessage, loginValue string) string {
	errorBlock := authErrorBlock(errorMessage)
	return authPageLayout("Вход", fmt.Sprintf(`<div class="card">
		<h1>Вход</h1>
		<p class="subtitle">Войдите в аккаунт Ladno Teams</p>
		%s
		<form method="POST" action="/login">
			<label for="login">Логин</label>
			<input id="login" name="login" type="text" autocomplete="username" required value="%s">
			<label for="password">Пароль</label>
			<input id="password" name="password" type="password" autocomplete="current-password" required>
			<button type="submit" class="full">Войти</button>
		</form>
		<p class="footer-link">Регистрация доступна только по ссылке-приглашению от команды.</p>
	</div>`, errorBlock, html.EscapeString(loginValue)))
}

func authRegisterInviteRequiredPage() string {
	return authPageLayout("Регистрация", `<div class="card">
		<h1>Регистрация</h1>
		<p class="subtitle">Создать аккаунт можно только по приглашению</p>
		<div class="info">Попросите лидера команды отправить вам ссылку с guid инвайта.</div>
		<p class="footer-link"><a href="/login">Войти</a></p>
	</div>`)
}

func authInviteRegisterPage(guid, teamName, errorMessage string, isExpired bool) string {
	errorBlock := authErrorBlock(errorMessage)
	infoBlock := ""
	if teamName != "" {
		infoBlock = fmt.Sprintf(`<div class="info">Приглашение в команду <strong>%s</strong></div>`, html.EscapeString(teamName))
	}
	if isExpired {
		return authPageLayout("Регистрация", fmt.Sprintf(`<div class="card card-wide">
			<h1>Регистрация</h1>
			%s
			<div class="error">Срок действия приглашения истёк</div>
			<p class="footer-link"><a href="/login">Войти</a></p>
		</div>`, infoBlock))
	}

	return authPageLayout("Регистрация", fmt.Sprintf(`<div class="card card-wide">
		<h1>Регистрация</h1>
		<p class="subtitle">Создайте аккаунт и присоединитесь к команде</p>
		%s
		%s
		<form method="POST" action="/register">
			<input type="hidden" name="guid" value="%s">
			<label for="name">Имя</label>
			<input id="name" name="name" type="text" required maxlength="255" autocomplete="name">
			<label for="login">Логин</label>
			<input id="login" name="login" type="text" autocomplete="username" required>
			<p class="hint">3–32 символа: латиница, цифры, _ и -</p>
			<label for="password">Пароль</label>
			<input id="password" name="password" type="password" autocomplete="new-password" required>
			<p class="hint">От 6 до 128 символов</p>
			<label for="password_confirm">Подтверждение пароля</label>
			<input id="password_confirm" name="password_confirm" type="password" autocomplete="new-password" required>
			<button type="submit" class="full">Зарегистрироваться</button>
		</form>
		<p class="footer-link">Уже есть аккаунт? <a href="/login">Войти</a></p>
	</div>`, infoBlock, errorBlock, html.EscapeString(guid)))
}

func authGuestHomePage() string {
	body := `<div class="card home-card">
		<h1>Ladno Teams</h1>
		<p class="subtitle">Обмен workspace с командами</p>
		<div class="home-actions">
			<a href="/login" class="btn">Войти</a>
		</div>
		<p class="footer-link">Новые пользователи регистрируются по ссылке-приглашению.</p>
	</div>`
	return authPageLayout("Главная", body)
}

func authHomePage(user *dto.UserDto) string {
	body := fmt.Sprintf(`<div class="card home-card">
		<h1>Добро пожаловать</h1>
		<p class="subtitle">Вы вошли как <strong>%s</strong></p>
		<div class="home-actions">
			<form method="POST" action="/logout" style="margin:0;">
				<button type="submit" class="btn btn-outline">Выйти</button>
			</form>
		</div>
		<p class="footer-link" style="margin-top:2rem;">Для работы с API используйте <code>POST /api/auth/login</code> или cookie <code>token</code> для refresh.</p>
	</div>`, html.EscapeString(user.Login))
	return authPageLayout("Главная", body)
}

func authErrorBlock(message string) string {
	if message == "" {
		return ""
	}
	return `<div class="error">` + html.EscapeString(message) + `</div>`
}
