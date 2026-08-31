package service

import (
	"ladno-teams/internal/config"
)

type BaseService struct {
	cfg config.Config
}

func NewBaseService(cfg config.Config) *BaseService {
	return &BaseService{cfg: cfg}
}
