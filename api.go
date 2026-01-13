package grepo

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

type APIOptions struct {
	fixedTime              *time.Time
	enableInputValidation  bool
	enableOutputValidation bool
	customFieldValidators  []FieldValidator
}

type APIOptionFunc func(*APIOptions)

func WithFixedTime(t time.Time) APIOptionFunc {
	return func(o *APIOptions) {
		o.fixedTime = &t
	}
}

func WithEnableInputValidation() APIOptionFunc {
	return func(o *APIOptions) {
		o.enableInputValidation = true
	}
}

func WithEnableOutputValidation() APIOptionFunc {
	return func(o *APIOptions) {
		o.enableOutputValidation = true
	}
}

func WithCustomFieldValidators(validators ...FieldValidator) APIOptionFunc {
	return func(o *APIOptions) {
		o.customFieldValidators = append(o.customFieldValidators, validators...)
	}
}

type API struct {
	m       map[string]Descriptor
	root    *Group
	options *APIOptions
}

func newAPI() *API {
	return &API{
		m:       make(map[string]Descriptor),
		root:    NewGroup("root"),
		options: &APIOptions{},
	}
}

func UseCase[I any, O any](api *API, op string) Executor[I, O] {
	return (ExecutorFunc[I, O](func(ctx context.Context, input I) (*O, error) {
		uc, ok := api.m[op]
		if !ok {
			return nil, ErrNotFound
		}

		interactor, ok := uc.(*Interactor[I, O])
		if !ok {
			return nil, ErrNotFound
		}

		if api.options.fixedTime != nil {
			ctx = withExecuteTime(ctx, *api.options.fixedTime)
		} else {
			ctx = withExecuteTime(ctx, time.Now())
		}

		groups := append([]*Group{api.root}, interactor.groups...)

		ctx, err := hookBefore(ctx, uc, input, groups)
		if err != nil {
			hookError(ctx, uc, input, err, groups)
			return nil, err
		}

		if api.options.enableInputValidation {
			if err := Validate(input, api.options.customFieldValidators...); err != nil {
				hookError(ctx, uc, input, err, groups)
				return nil, err
			}
		}

		output, err := interactor.Execute(ctx, input)
		if err != nil {
			hookError(ctx, uc, input, err, groups)
			return nil, err
		}

		if api.options.enableOutputValidation {
			if err := Validate(output, api.options.customFieldValidators...); err != nil {
				hookError(ctx, uc, input, err, groups)
				return nil, err
			}
		}

		hookAfter(ctx, uc, input, output, groups)
		return output, nil
	}))
}

func (a *API) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(a.m))
	for k := range a.m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	b := strings.Builder{}

	b.WriteString("{")

	for i, k := range keys {
		if i > 0 {
			b.WriteString(",")
		}
		d := a.m[k]

		ucJSON, _ := json.Marshal(d)
		b.WriteString(fmt.Sprintf("%q: %s", k, ucJSON))
	}

	b.WriteString("}")

	return []byte(b.String()), nil
}

type APIBuilder struct {
	api *API
}

func NewAPIBuilder() *APIBuilder {
	return &APIBuilder{
		api: newAPI(),
	}
}

func (b *APIBuilder) WithHook(hook *GroupHook) *APIBuilder {
	b.api.root.hook = hook
	return b
}

func (b *APIBuilder) WithUseCase(d Descriptor) *APIBuilder {
	op := d.Operation()
	b.api.m[op] = d
	return b
}

func (b *APIBuilder) WithOptions(opts ...APIOptionFunc) *APIBuilder {
	for _, opt := range opts {
		opt(b.api.options)
	}
	return b
}

func (b *APIBuilder) Build() *API {
	return b.api
}
