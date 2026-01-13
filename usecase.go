package grepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/plinx2/grepo/refl"
)

type UseCaseHook[I any, O any] struct {
	before []func(ctx context.Context, i I) (context.Context, error)
	after  []func(ctx context.Context, i I, o *O)
	error  []func(ctx context.Context, i I, e error)
}

func NewUseCaseHook[I any, O any]() *UseCaseHook[I, O] {
	return &UseCaseHook[I, O]{}
}

func (h *UseCaseHook[I, O]) AddBefore(hook func(ctx context.Context, i I) (context.Context, error)) *UseCaseHook[I, O] {
	h.before = append(h.before, hook)
	return h
}

func (h *UseCaseHook[I, O]) AddAfter(hook func(ctx context.Context, i I, o *O)) *UseCaseHook[I, O] {
	h.after = append(h.after, hook)
	return h
}

func (h *UseCaseHook[I, O]) AddError(hook func(ctx context.Context, i I, e error)) *UseCaseHook[I, O] {
	h.error = append(h.error, hook)
	return h
}

type Executor[I any, O any] interface {
	Execute(ctx context.Context, input I) (*O, error)
}

type ExecutorFunc[I any, O any] func(context.Context, I) (*O, error)

func (fn ExecutorFunc[I, O]) Execute(ctx context.Context, input I) (*O, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return fn(ctx, input)
}

type Descriptor interface {
	Operation() string
	Input() any
	Output() any
	Groups() []string
}

type Interactor[I any, O any] struct {
	uc     Executor[I, O]
	op     string
	hook   *UseCaseHook[I, O]
	groups []*Group
}

func newInteractor[I any, O any](uc Executor[I, O]) *Interactor[I, O] {
	return &Interactor[I, O]{
		uc:   uc,
		hook: NewUseCaseHook[I, O](),
	}
}

func (i *Interactor[I, O]) Execute(ctx context.Context, input I) (*O, error) {
	var err error
	for _, beforeHook := range i.hook.before {
		ctx, err = beforeHook(ctx, input)
		if err != nil {
			return nil, err
		}
		if ctx == nil {
			return nil, errors.New("context is nil in before hook")
		}
	}

	output, err := i.uc.Execute(ctx, input)
	if err != nil {
		for _, errorHook := range i.hook.error {
			errorHook(ctx, input, err)
		}
		return nil, err
	}

	for _, afterHook := range i.hook.after {
		afterHook(ctx, input, output)
	}
	return output, nil
}

func (i *Interactor[I, O]) Operation() string {
	if i.op == "" {
		rt := reflect.TypeOf(i.uc)
		for rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
		}
		return rt.Name()
	}
	return i.op
}

func (i *Interactor[I, O]) Input() any {
	input := new(I)
	return *input
}

func (i *Interactor[I, O]) Output() any {
	output := new(O)
	return *output
}

func (i *Interactor[I, O]) Groups() []string {
	g := make([]string, 0, len(i.groups))
	for _, group := range i.groups {
		g = append(g, group.name)
	}
	return g
}

func (i *Interactor[I, O]) MarshalJSON() ([]byte, error) {
	b := strings.Builder{}

	b.WriteString("{")

	operation := i.Operation()
	b.WriteString(fmt.Sprintf("%q: %q", "Operation", operation))

	inut := i.Input()
	inputSpec := refl.TypeOf(inut)
	inputJSON, _ := json.Marshal(inputSpec)
	b.WriteString(",")
	b.WriteString(fmt.Sprintf("%q: %s", "Input", inputJSON))

	output := i.Output()
	outputSpec := refl.TypeOf(output)
	outputJSON, _ := json.Marshal(outputSpec)
	b.WriteString(",")
	b.WriteString(fmt.Sprintf("%q: %s", "Output", outputJSON))

	groups := i.Groups()
	groupsJSON, _ := json.Marshal(groups)
	b.WriteString(",")
	b.WriteString(fmt.Sprintf("%q: %s", "Groups", groupsJSON))

	b.WriteString("}")

	return []byte(b.String()), nil
}

type UseCaseBuilder[I any, O any] struct {
	uc *Interactor[I, O]
}

func NewUseCaseBuilder[I any, O any](uc Executor[I, O]) *UseCaseBuilder[I, O] {
	return &UseCaseBuilder[I, O]{
		uc: newInteractor(uc),
	}
}

func (b *UseCaseBuilder[I, O]) WithOperation(name string) *UseCaseBuilder[I, O] {
	b.uc.op = name
	return b
}

func (b *UseCaseBuilder[I, O]) WithHook(hook *UseCaseHook[I, O]) *UseCaseBuilder[I, O] {
	b.uc.hook = hook
	return b
}

func (b *UseCaseBuilder[I, O]) WithGroup(group *Group) *UseCaseBuilder[I, O] {
	b.uc.groups = append(b.uc.groups, group)
	return b
}

func (b *UseCaseBuilder[I, O]) Build() *Interactor[I, O] {
	return b.uc
}
