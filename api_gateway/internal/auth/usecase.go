package auth

import (
	"context"
	"go-futures-api/internal/models"
)

type UseCase interface {
	FindOne(ctx context.Context, userId string) (*models.User, error)
}
