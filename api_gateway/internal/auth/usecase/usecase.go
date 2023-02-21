package usecase

import (
	"context"
	"go-futures-api/internal/models"
	authService "go-futures-api/proto/auth"
	"go.uber.org/zap"
)

type authUseCase struct {
	logger      *zap.SugaredLogger
	authService authService.UserServiceClient
}

func NewAuthUseCase(logger *zap.SugaredLogger, authService authService.UserServiceClient) *authUseCase {
	return &authUseCase{
		logger:      logger,
		authService: authService,
	}
}

func (a *authUseCase) FindOne(ctx context.Context, userId string) (*models.User, error) {
	userRes, err := a.authService.FindOne(ctx, &authService.UserById{Id: userId})
	if err != nil {
		return nil, err
	}

	if userRes.Id == "" {
		return nil, nil
	}

	user, err := models.UserFromProto(userRes)
	if err != nil {
		return nil, err
	}
	return user, nil
}
