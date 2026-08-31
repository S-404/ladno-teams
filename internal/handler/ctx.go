package handler

import (
	"errors"
	"ladno-teams/internal/entity/dto"

	"github.com/gin-gonic/gin"
)

const (
	CtxUser = "user"
)

func GetCtxUser(c *gin.Context) (*dto.UserDto, error) {
	user, ok := c.Get(CtxUser)
	if !ok {
		return nil, errors.New("ctx user not found")
	}

	userData, ok := user.(dto.UserDto)
	if !ok {
		return nil, errors.New("ctx user is of invalid type")
	}

	return &userData, nil
}
