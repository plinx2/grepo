package port

import (
	"context"

	"github.com/plinx2/grepo/example/entity"
)

type RepoUser interface {
	GetUser(ctx context.Context, id string) (*entity.User, error)
	SaveUser(ctx context.Context, user *entity.User) error
}
