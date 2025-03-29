package knowledge

import (
	"context"
	"support/internal/models"
)

type Repo interface {
	Get(ctx context.Context, intent string) (models.Knowledge, error)
	Exist(ctx context.Context, intent string) bool
}
