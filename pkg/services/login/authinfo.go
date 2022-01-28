package login

import (
	"context"

	"github.com/grafana/grafana/pkg/models"
)

type AuthInfoService interface {
	LookupAndUpdate(ctx context.Context, query *models.GetUserByAuthInfoQuery) (*models.User, error)
	UpdateAuthInfo(ctx context.Context, cmd *models.UpdateAuthInfoCommand) error
	SetAuthInfo(ctx context.Context, cmd *models.SetAuthInfoCommand) error
}
