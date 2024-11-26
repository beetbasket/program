package program

import (
	"context"
	"github.com/trymoose/debug"
	"go.uber.org/fx"
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

func execSubCommand(subcommand string, ctx context.Context, args []string) {
	sc, ok := subcommands[subcommand]
	if ok {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		pr := fx.New(
			append(
				sc.Module,
				fx.Supply(OsArgs(args)),
				fx.Provide(fx.Annotate(
					fx.Supply(&Context{ctx}),
					fx.As(new(context.Context)),
				)),
				fx.Invoke(sc.Main),
			)...,
		)

		go func() {
			defer cancel()
			select {
			case <-pr.Done():
			case <-pr.Wait():
			}
		}()

		pr.Run()
	} else if subcommand != "" {
		execSubCommand("", ctx, append([]string{subcommand}, args...))
	} else {
		panic("no main command set")
	}
}
