package main

import (
	"context"
	"fmt"

	"github.com/Olian04/go-me/pkg/aggregate"
	"github.com/Olian04/go-me/pkg/gnu"
	"github.com/urfave/cli/v3"
)

func runID(ctx context.Context, cmd *cli.Command) error {
	opt := gnu.IDOptions{
		User:   cmd.Bool("user"),
		Group:  cmd.Bool("group"),
		Groups: cmd.Bool("groups"),
		Name:   cmd.Bool("name"),
		Real:   cmd.Bool("real"),
	}

	payload, err := aggregate.Aggregate(ctx, aggregate.Options{
		Timeout:    defaultTimeout,
		Sources:    nil,
		BestEffort: true,
		Strict:     false,
	})
	if err != nil {
		return mapAggErr(err)
	}
	out := gnu.FormatID(payload, opt)
	_, _ = fmt.Fprintln(cmd.Writer, out)
	return nil
}
