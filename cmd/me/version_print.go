package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Olian04/go-me/cmd/me/render"
	"github.com/Olian04/go-me/cmd/me/version"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

// printVersionForCLI writes build metadata in the same output modes as identity (--text/--compact/--json/--yaml).
// Called from cli.VersionPrinter before the root Action runs.
func printVersionForCLI(cmd *cli.Command) {
	r := cmd.Root()
	flags := parseRootFlags(r)
	textOut, compactOut, jsonOut, yamlOut, err := resolveRootOutput(flags)
	if err != nil {
		_, _ = fmt.Fprintln(r.ErrWriter, err.Error())
		os.Exit(exitUsage)
	}

	info := version.Get()
	w := r.Writer

	switch {
	case compactOut:
		_, _ = fmt.Fprintf(w, "%s/%s/%s\n", info.Version, info.Revision, info.BuildTime)
	case jsonOut:
		b, err := json.Marshal(info)
		if err != nil {
			_, _ = fmt.Fprintln(r.ErrWriter, err.Error())
			os.Exit(exitInternal)
		}
		_, _ = fmt.Fprintln(w, string(b))
	case yamlOut:
		b, err := yaml.Marshal(info)
		if err != nil {
			_, _ = fmt.Fprintln(r.ErrWriter, err.Error())
			os.Exit(exitInternal)
		}
		_, _ = fmt.Fprint(w, string(b))
	case textOut:
		_, _ = fmt.Fprintln(w, render.VersionText(info))
	default:
		_, _ = fmt.Fprintln(r.ErrWriter, "no output mode")
		os.Exit(exitUsage)
	}
}
