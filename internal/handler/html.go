package handler

import (
	"log"
	"net/http"

	"ladno-teams/internal/views"

	"github.com/gin-gonic/gin"
)

func renderView(viewsEngine *views.Engine, c *gin.Context, name string, data any) {
	html, err := viewsEngine.Render(name, data)
	if err != nil {
		log.Printf("render view %s: %v", name, err)
		c.String(http.StatusInternalServerError, "template error")
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func renderAdminPage(viewsEngine *views.Engine, c *gin.Context, page string, data any) {
	html, err := viewsEngine.RenderLayout("layouts/admin", page, data)
	if err != nil {
		log.Printf("render admin page %s: %v", page, err)
		c.String(http.StatusInternalServerError, "template error")
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func renderAuthPage(viewsEngine *views.Engine, c *gin.Context, title, page string, data any) {
	html, err := viewsEngine.RenderAuth(title, page, data)
	if err != nil {
		log.Printf("render auth page %s: %v", page, err)
		c.String(http.StatusInternalServerError, "template error")
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func renderAdminLoginPage(viewsEngine *views.Engine, c *gin.Context, data views.AuthPage) {
	html, err := viewsEngine.RenderAdminLogin(data)
	if err != nil {
		log.Printf("render admin login: %v", err)
		c.String(http.StatusInternalServerError, "template error")
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func renderHTMLCloseModal(viewsEngine *views.Engine, c *gin.Context, page string, data any) {
	c.Header("HX-Trigger", `{"closeModal":{"target":"body"}}`)
	html, err := viewsEngine.RenderCloseModal(page, data)
	if err != nil {
		log.Printf("render close modal %s: %v", page, err)
		c.String(http.StatusInternalServerError, "template error")
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (h *Handler) renderView(c *gin.Context, name string, data any) {
	renderView(h.views, c, name, data)
}

func (h *Handler) renderAdminPage(c *gin.Context, page string, data any) {
	renderAdminPage(h.views, c, page, data)
}

func (h *Handler) renderAuthPage(c *gin.Context, title, page string, data any) {
	renderAuthPage(h.views, c, title, page, data)
}

func (h *Handler) renderAdminLoginPage(c *gin.Context, data views.AuthPage) {
	renderAdminLoginPage(h.views, c, data)
}

func (h *Handler) renderHTMLCloseModal(c *gin.Context, page string, data any) {
	renderHTMLCloseModal(h.views, c, page, data)
}
