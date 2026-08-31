package views

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed layouts pages partials ui
var files embed.FS

type Engine struct {
	templates *template.Template
}

func NewEngine() (*Engine, error) {
	tpl, err := template.New("").Funcs(templateFuncs()).ParseFS(
		files,
		"layouts/*.gohtml",
		"pages/admin/*.gohtml",
		"pages/auth/*.gohtml",
		"partials/admin/*.gohtml",
		"partials/admin/modals/*.gohtml",
		"ui/partials/*.gohtml",
	)
	if err != nil {
		return nil, fmt.Errorf("parse views: %w", err)
	}
	return &Engine{templates: tpl}, nil
}

func (e *Engine) Render(name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := e.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (e *Engine) RenderLayout(layout, page string, pageData any) (string, error) {
	content, err := e.Render(page, pageData)
	if err != nil {
		return "", err
	}
	return e.Render(layout, LayoutData{Content: template.HTML(content)})
}

func (e *Engine) RenderAuth(title, page string, pageData any) (string, error) {
	content, err := e.Render(page, pageData)
	if err != nil {
		return "", err
	}
	return e.Render("layouts/auth", LayoutData{Title: title, Content: template.HTML(content)})
}

func (e *Engine) RenderAdminLogin(pageData AuthPage) (string, error) {
	content, err := e.Render("pages/admin/login", pageData)
	if err != nil {
		return "", err
	}
	return e.Render("layouts/admin_login", LayoutData{Content: template.HTML(content)})
}

func (e *Engine) RenderCloseModal(page string, pageData any) (string, error) {
	content, err := e.Render(page, pageData)
	if err != nil {
		return "", err
	}
	return content + `<div id="modal-host" hx-swap-oob="innerHTML"></div>`, nil
}

func (e *Engine) RegisterAssets(router *gin.Engine) {
	styles, err := fs.Sub(files, "ui/styles")
	if err != nil {
		panic(err)
	}
	scripts, err := fs.Sub(files, "ui/scripts")
	if err != nil {
		panic(err)
	}

	router.StaticFS("/assets/styles", http.FS(styles))
	router.StaticFS("/assets/scripts", http.FS(scripts))
}
