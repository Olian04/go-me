package main

import (
	"context"
	"errors"
	"time"

	"github.com/Olian04/go-me/cmd/me/version"
	"github.com/Olian04/go-me/pkg/aggregate"
	"github.com/urfave/cli/v3"
)

const defaultTimeout = 2 * time.Second

func init() {
	cli.VersionPrinter = printVersionForCLI
}

// exit codes per docs/design/cli-command-matrix.md
const (
	exitUsage    = 2
	exitStrict   = 3
	exitInternal = 4
)

func mapAggErr(err error) error {
	var sue aggregate.StrictUnknownSourceError
	if errors.As(err, &sue) {
		return cli.Exit(err.Error(), exitUsage)
	}
	var spe aggregate.StrictProviderError
	if errors.As(err, &spe) {
		return cli.Exit(err.Error(), exitStrict)
	}
	return cli.Exit(err.Error(), exitInternal)
}

func compatBefore(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	r := cmd.Root()
	if identityFlagsSet(r) {
		return ctx, cli.Exit("identity flags are not allowed with this command", exitUsage)
	}
	return ctx, nil
}

func identityFlagsSet(r *cli.Command) bool {
	names := []string{"text", "t", "compact", "c", "json", "yaml", "source", "strict", "timeout", "no-color"}
	for _, n := range names {
		if r.IsSet(n) {
			return true
		}
	}
	return false
}

func newApp() *cli.Command {
	return &cli.Command{
		Name:  "me",
		Usage: "show identity summary",
		Description: "Print identity for the current user and environment using identity providers.\n\n" +
			"With no --source, runs: osaccount (local user account), envcontext (sudo/SSH/CI hints), " +
			"network (hostname, FQDN, addresses), sysinfo (GOOS/GOARCH and OS name/version). " +
			"Use --source authproviders to include git user/email and cloud hints. " +
			"If you set --source, list only the providers you want; it replaces the full default list.",
		Version: version.Get().Version,
		// Avoid OsExiter inside Run so tests can inspect errors; main calls HandleExitCoder.
		ExitErrHandler: func(context.Context, *cli.Command, error) {},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "text", Aliases: []string{"t"}, Usage: "human-readable text output"},
			&cli.BoolFlag{Name: "compact", Aliases: []string{"c"}, Usage: "compact fingerprint output"},
			&cli.BoolFlag{Name: "json", Usage: "JSON output"},
			&cli.BoolFlag{Name: "yaml", Usage: "YAML output"},
			&cli.BoolFlag{Name: "no-color", Usage: "disable ANSI colors"},
			&cli.DurationFlag{Name: "timeout", Value: defaultTimeout, Usage: "provider deadline"},
			&cli.StringSliceFlag{Name: "source", Usage: aggregate.SourceFlagUsage},
			&cli.BoolFlag{Name: "strict", Usage: "strict validation; unknown sources fail"},
		},
		Action: runRoot,
		Commands: []*cli.Command{
			{
				Name:   "whoami",
				Usage:  "print effective user name (GNU whoami compatibility)",
				Before: compatBefore,
				Action: runWhoami,
			},
			{
				Name:  "id",
				Usage: "print user/group ids (GNU id compatibility subset)",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "user", Aliases: []string{"u"}, Usage: "print effective user id only"},
					&cli.BoolFlag{Name: "group", Aliases: []string{"g"}, Usage: "print effective group id only"},
					&cli.BoolFlag{Name: "groups", Aliases: []string{"G"}, Usage: "print supplementary group ids"},
					&cli.BoolFlag{Name: "name", Aliases: []string{"n"}, Usage: "print names instead of numeric ids"},
					&cli.BoolFlag{Name: "real", Aliases: []string{"r"}, Usage: "print real id instead of effective id"},
				},
				Before: compatBefore,
				Action: runID,
			},
		},
	}
}
