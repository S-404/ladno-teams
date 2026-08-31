package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BaseHandler struct {
	validate *validator.Validate
}

func NewBaseHandler(validate *validator.Validate) *BaseHandler {
	return &BaseHandler{validate: validate}
}

func (h *BaseHandler) ValidateRequestBody(c *gin.Context, requestModel any) error {
	if err := c.ShouldBindJSON(requestModel); err != nil {
		return err
	}

	if err := h.validate.Struct(requestModel); err != nil {
		return err
	}

	return nil
}

func (h *BaseHandler) parseUUID(c *gin.Context, param string) (uuid.UUID, bool) {
	value, err := uuid.Parse(c.Param(param))
	if err != nil {
		return uuid.Nil, false
	}
	return value, true
}

func (h *BaseHandler) parseQueryInt(c *gin.Context, param string, defaultValue int) int {
	value := c.Query(param)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
