package program

import (
	"context"
	"github.com/google/uuid"
	"github.com/trymoose/debug"
	"log/slog"
	"os"
)

type SubcommandFn func(ctx context.Context, args []string) error

var DefaultSubcommand = uuid.NewString()

var subcommands = map[string]SubcommandFn{
	DefaultSubcommand: func(ctx context.Context, args []string) error { return nil },
}

func Subcommand(name string, command SubcommandFn) {
	subcommands[name] = command
}

func Main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if debug.Debug {
		slog.Debug("running in debug mode")
	}

	if err := runSubCommand(ctx); err != nil {
		slog.Error("exited with error", slog.Any("error", err))
	}
}

func runSubCommand(ctx context.Context) error {
	if len(os.Args) > 1 {
		return execSubCommand(os.Args[1], ctx, os.Args[2:])
	}
	return execSubCommand(DefaultSubcommand, ctx, os.Args[1:])
}

func execSubCommand(subcommand string, ctx context.Context, args []string) error {
	fn, ok := subcommands[subcommand]
	if ok {
		return fn(ctx, args)
	}
	return execSubCommand(DefaultSubcommand, ctx, append([]string{subcommand}, args...))
}
