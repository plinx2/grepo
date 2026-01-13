package main

import (
	"github.com/plinx2/grepo"
	"github.com/plinx2/grepo/example/usecase"
	"github.com/spf13/cobra"
)

var findUsersInput usecase.FindUsersInput

var findUsersCmd = &cobra.Command{
	Use:   "find-users",
	Short: "Find users by criteria",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		api := ctx.Value("api").(*grepo.API)

		output, err := grepo.UseCase[usecase.FindUsersInput, usecase.FindUsersOutput](api, usecase.FindUsersOperation).Execute(ctx, findUsersInput)
		if err != nil {
			return err
		}
		return PrintJSON(output)
	},
}

func init() {
}
