package example

import (
	"fmt"
	"reflect"

	"github.com/plinx2/grepo"
	"github.com/plinx2/grepo/example/usecase"
	"github.com/plinx2/grepo/hooks"
	"github.com/plinx2/grepo/refl"
)

func NewAPI(
	findUser grepo.Executor[usecase.FindUsersInput, usecase.FindUsersOutput],
	getUser grepo.Executor[usecase.GetUserInput, usecase.GetUserOutput],
	saveUser grepo.Executor[usecase.SaveUserInput, usecase.SaveUserOutput],
) *grepo.API {
	rootHook := grepo.NewGroupHook()
	rootHook.AddBefore(hooks.HookBeforeSlog())
	rootHook.AddAfter(hooks.HookAfterSlog())
	rootHook.AddError(hooks.HookErrorSlog())

	return grepo.NewAPIBuilder().
		WithHook(rootHook).
		WithOptions(
			grepo.WithEnableInputValidation(),
			grepo.WithEnableOutputValidation(),
			grepo.WithCustomFieldValidators((grepo.FieldValidatorFunc(func(v reflect.Value, f *refl.Field) error {
				fmt.Println(f.Parent().Name, f.Field, f.Type.Name)
				return nil
			}))),
		).
		WithUseCase(
			grepo.NewUseCaseBuilder(findUser).
				Build(),
		).
		WithUseCase(
			grepo.NewUseCaseBuilder(getUser).
				Build(),
		).
		WithUseCase(
			grepo.NewUseCaseBuilder(saveUser).
				Build(),
		).
		Build()
}
