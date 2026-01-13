package local

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/plinx2/grepo/example/entity"
	"github.com/plinx2/grepo/example/port"
)

type RepoUser struct {
	path  string
	users map[string]*entity.User
}

var _ port.RepoUser = (*RepoUser)(nil)

func NewRepoUser(dir string) *RepoUser {
	path := filepath.Join(dir, "users.json")

	b, _ := os.ReadFile(path)

	m := make(map[string]*entity.User)
	json.Unmarshal(b, &m)

	return &RepoUser{
		path:  path,
		users: m,
	}
}

func (r *RepoUser) GetUser(ctx context.Context, id string) (*entity.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (r *RepoUser) SaveUser(ctx context.Context, user *entity.User) error {
	r.users[user.ID] = user
	b, err := json.MarshalIndent(r.users, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(r.path, b, 0o644); err != nil {
		return err
	}
	return nil
}
