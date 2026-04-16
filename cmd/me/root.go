package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Olian04/go-me/cmd/me/render"
	"github.com/Olian04/go-me/pkg/aggregate"
	"github.com/Olian04/go-me/pkg/compact"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func runRoot(ctx context.Context, cmd *cli.Command) error {
	r := cmd.Root()
	flags := parseRootFlags(r)
	textOut, compactOut, jsonOut, yamlOut, err := resolveRootOutput(flags)
	if err != nil {
		return cli.Exit(err.Error(), exitUsage)
	}

	strict := flags.Strict
	timeout := r.Duration("timeout")
	if !r.IsSet("timeout") {
		timeout = defaultTimeout
	}

	opts := aggregate.Options{
		Timeout:    timeout,
		Sources:    aggregate.ParseSourceFlags(r.StringSlice("source")),
		BestEffort: !strict,
		Strict:     strict,
	}

	payload, err := aggregate.Aggregate(ctx, opts)
	if err != nil {
		return mapAggErr(err)
	}

	switch {
	case compactOut:
		s, err := compact.FormatOrStrict(payload, strict)
		if err != nil {
			var ce compact.ErrStrictCompact
			if errors.As(err, &ce) {
				return cli.Exit(err.Error(), exitStrict)
			}
			return cli.Exit(err.Error(), exitInternal)
		}
		_, _ = fmt.Fprintln(r.Writer, s)
		return nil
	case jsonOut:
		b, err := json.Marshal(payload)
		if err != nil {
			return cli.Exit(err.Error(), exitInternal)
		}
		_, _ = fmt.Fprintln(r.Writer, string(b))
		return nil
	case yamlOut:
		b, err := yaml.Marshal(payload)
		if err != nil {
			return cli.Exit(err.Error(), exitInternal)
		}
		_, _ = fmt.Fprint(r.Writer, string(b))
		return nil
	case textOut:
		_ = flags.NoColor // reserved for future ANSI styling
		_, _ = fmt.Fprintln(r.Writer, render.Text(payload))
		return nil
	default:
		return cli.Exit("no output mode", exitUsage)
	}
}
