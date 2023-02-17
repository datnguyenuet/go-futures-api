package middleware

import (
	"go.uber.org/zap"
)

type Manager struct {
	logger *zap.SugaredLogger
}

func NewMiddlewareManager(logger *zap.SugaredLogger) *Manager {
	return &Manager{logger: logger}
}
