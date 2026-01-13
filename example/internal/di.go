package internal

import (
	"github.com/plinx2/grepo"
	"github.com/plinx2/grepo/example"
	"github.com/plinx2/grepo/example/internal/local"
	"github.com/plinx2/grepo/example/usecase"
)

func InitializeAPI() *grepo.API {
	repoUser := local.NewRepoUser(".")

	ucFindUser := usecase.NewFindUsers(repoUser)
	ucGetUser := usecase.NewGetUser(repoUser)
	ucSaveUser := usecase.NewSaveUser(repoUser)

	return example.NewAPI(ucFindUser, ucGetUser, ucSaveUser)
}
