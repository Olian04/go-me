package main

import (
	"context"
	"fmt"

	"github.com/Olian04/go-me/pkg/aggregate"
	"github.com/Olian04/go-me/pkg/gnu"
	"github.com/urfave/cli/v3"
)

func runWhoami(ctx context.Context, cmd *cli.Command) error {
	payload, err := aggregate.Aggregate(ctx, aggregate.Options{
		Timeout:    defaultTimeout,
		Sources:    nil,
		BestEffort: true,
		Strict:     false,
	})
	if err != nil {
		return mapAggErr(err)
	}
	u := gnu.FormatWhoami(payload)
	if u == "" {
		return cli.Exit("cannot determine effective user", exitStrict)
	}
	_, _ = fmt.Fprintln(cmd.Writer, u)
	return nil
}
