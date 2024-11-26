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

var subcommands = map[string]*Subcommand{
	"": {
		Main: func() {},
	},
}

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

func execSubCommand(subcommand string, ctx context.Context, args []string) {
	sc, ok := subcommands[subcommand]
	if ok {
		fx.New(
			append(
				sc.Module,
				fx.Supply(OsArgs(args)),
				fx.Invoke(sc.Main),
			)...,
		).Run()
	} else if subcommand != "" {
		execSubCommand("", ctx, append([]string{subcommand}, args...))
	} else {
		panic("no main command set")
	}
}
