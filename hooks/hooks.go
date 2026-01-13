package hooks

import (
	"context"
	"log/slog"

	"github.com/plinx2/grepo"
)

type HookSlogOptions struct {
	level slog.Level
	msg   string
}

type HookSlogOptionFunc func(*HookSlogOptions)

func WithLevel(level slog.Level) HookSlogOptionFunc {
	return func(o *HookSlogOptions) {
		o.level = level
	}
}

func WithMsg(msg string) HookSlogOptionFunc {
	return func(o *HookSlogOptions) {
		o.msg = msg
	}
}

func HookBeforeSlog(opts ...HookSlogOptionFunc) grepo.BeforeHook {
	options := &HookSlogOptions{
		level: slog.LevelInfo,
		msg:   "Starting operation",
	}
	for _, opt := range opts {
		opt(options)
	}
	return func(ctx context.Context, desc grepo.Descriptor, i any) (context.Context, error) {
		slog.Log(ctx, options.level, options.msg, "operation", desc.Operation(), "input", i)
		return ctx, nil
	}
}

func HookAfterSlog(opts ...HookSlogOptionFunc) grepo.AfterHook {
	options := &HookSlogOptions{
		level: slog.LevelInfo,
		msg:   "Finished operation",
	}
	for _, opt := range opts {
		opt(options)
	}
	return func(ctx context.Context, desc grepo.Descriptor, i any, o any) {
		slog.Log(ctx, options.level, options.msg, "operation", desc.Operation(), "output", o)
	}
}

func HookErrorSlog(opts ...HookSlogOptionFunc) grepo.ErrorHook {
	options := &HookSlogOptions{
		level: slog.LevelError,
		msg:   "Operation error",
	}
	for _, opt := range opts {
		opt(options)
	}
	return func(ctx context.Context, desc grepo.Descriptor, i any, e error) {
		slog.Log(ctx, options.level, options.msg, "operation", desc.Operation(), "input", i, "error", e)
	}
}
