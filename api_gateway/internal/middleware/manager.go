package middleware

import (
	authUseCase "go-futures-api/internal/auth"
	"go.uber.org/zap"
)

type Manager struct {
	logger      *zap.SugaredLogger
	authUseCase authUseCase.UseCase
}

func NewMiddlewareManager(logger *zap.SugaredLogger, authUseCase authUseCase.UseCase) *Manager {
	return &Manager{
		logger:      logger,
		authUseCase: authUseCase,
	}
}
