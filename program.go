package program

import (
	"context"
	"github.com/trymoose/debug"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	"os"
)

import _ "github.com/beetbasket/program/pkg/log"

type Subcommand struct {
	Name   string
	Module []fx.Option
	Main   any
}

var subcommands = map[string]*Subcommand{}

func Register(sc *Subcommand) {
	subcommands[sc.Name] = sc
}

func Main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if debug.Debug {
		slog.Debug("running in debug mode")
	}

	runSubCommand(ctx)
}

func runSubCommand(ctx context.Context) {
	if len(os.Args) > 1 {
		execSubCommand(os.Args[1], ctx, os.Args[2:])
	} else {
		execSubCommand("", ctx, os.Args[1:])
	}
}

type OsArgs []string
type Context struct{ context.Context }

var _ context.Context = Context{}

func execSubCommand(subcommand string, ctx context.Context, args []string) {
	sc, ok := subcommands[subcommand]
	if ok {
		var opts []fx.Option
		if sc.Main != nil {
			opts = append(opts, fx.Invoke(sc.Main))
		}

		fx.New(
			fx.Options(sc.Module...),
			fx.Provide(
				slog.Default,
				osArgs(args),
				fx.Annotate(
					provideContext,
					fx.As(new(context.Context)),
				),
			),
			fx.WithLogger(withSlogger),
			fx.Options(opts...),
		).Run()
	} else if subcommand != "" {
		execSubCommand("", ctx, append([]string{subcommand}, args...))
	} else {
		panic("no main command set")
	}
}

func osArgs(args []string) func() OsArgs {
	return func() OsArgs {
		return args
	}
}

func withSlogger(slogger *slog.Logger) fxevent.Logger {
	return &fxevent.SlogLogger{
		Logger: slogger,
	}
}

func provideContext(lf fx.Lifecycle) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	lf.Append(fx.StopHook(cancel))
	return ctx
}
