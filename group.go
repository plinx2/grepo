package grepo

import "context"

type BeforeHook func(ctx context.Context, desc Descriptor, i any) (context.Context, error)
type AfterHook func(ctx context.Context, desc Descriptor, i any, o any)
type ErrorHook func(ctx context.Context, desc Descriptor, i any, e error)

type GroupHook struct {
	before []BeforeHook
	after  []AfterHook
	error  []ErrorHook
}

func NewGroupHook() *GroupHook {
	return &GroupHook{}
}

func (h *GroupHook) AddBefore(hook BeforeHook) *GroupHook {
	h.before = append(h.before, hook)
	return h
}

func (h *GroupHook) AddAfter(hook AfterHook) *GroupHook {
	h.after = append(h.after, hook)
	return h
}

func (h *GroupHook) AddError(hook ErrorHook) *GroupHook {
	h.error = append(h.error, hook)
	return h
}

type Group struct {
	name string
	hook *GroupHook
}

func NewGroup(name string) *Group {
	return &Group{
		name: name,
		hook: NewGroupHook(),
	}
}

func (g *Group) Name() string {
	return g.name
}

func (g *Group) MarshalJSON() ([]byte, error) {
	return []byte(`"` + g.name + `"`), nil
}

func doHookBefore(ctx context.Context, desc Descriptor, input any, hooks []BeforeHook) (context.Context, error) {
	for _, hook := range hooks {
		var err error
		ctx, err = hook(ctx, desc, input)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

func doHookAfter(ctx context.Context, desc Descriptor, input any, output any, hooks []AfterHook) {
	for _, hook := range hooks {
		hook(ctx, desc, input, output)
	}
}

func doHookError(ctx context.Context, desc Descriptor, input any, err error, hooks []ErrorHook) {
	for _, hook := range hooks {
		hook(ctx, desc, input, err)
	}
}

func hookBefore(ctx context.Context, desc Descriptor, input any, groups []*Group) (context.Context, error) {
	for _, g := range groups {
		var err error
		for _, beforeHook := range g.hook.before {
			ctx, err = beforeHook(ctx, desc, input)
			if err != nil {
				return ctx, err
			}
		}
	}
	if ctx == nil {
		return nil, ErrNotFound
	}
	return ctx, nil
}

func hookAfter(ctx context.Context, desc Descriptor, input any, output any, groups []*Group) {
	for i := len(groups) - 1; i >= 0; i-- {
		g := groups[i]
		for _, afterHook := range g.hook.after {
			afterHook(ctx, desc, input, output)
		}
	}
}

func hookError(ctx context.Context, desc Descriptor, input any, err error, groups []*Group) {
	for i := len(groups) - 1; i >= 0; i-- {
		g := groups[i]
		for _, errorHook := range g.hook.error {
			errorHook(ctx, desc, input, err)
		}
	}
}
