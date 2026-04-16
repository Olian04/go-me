package main

import (
	"errors"

	"github.com/urfave/cli/v3"
)

// errMutuallyExclusiveOutputModes is returned when more than one of text/compact/json/yaml is set.
var errMutuallyExclusiveOutputModes = errors.New("only one of --text, --compact, --json, or --yaml may be set")

type rootFlags struct {
	Text    bool
	Compact bool
	JSON    bool
	YAML    bool
	NoColor bool
	Strict  bool
}

func parseRootFlags(r *cli.Command) rootFlags {
	return rootFlags{
		Text:    r.Bool("text"),
		Compact: r.Bool("compact"),
		JSON:    r.Bool("json"),
		YAML:    r.Bool("yaml"),
		NoColor: r.Bool("no-color"),
		Strict:  r.Bool("strict"),
	}
}

// resolveRootOutput derives which output mode is active. With no output flags, text is the default.
func resolveRootOutput(f rootFlags) (textOut, compactOut, jsonOut, yamlOut bool, err error) {
	outModes := 0
	if f.Text {
		outModes++
	}
	if f.Compact {
		outModes++
	}
	if f.JSON {
		outModes++
	}
	if f.YAML {
		outModes++
	}
	if outModes > 1 {
		return false, false, false, false, errMutuallyExclusiveOutputModes
	}
	textOut = outModes == 0 || f.Text
	compactOut = f.Compact
	jsonOut = f.JSON
	yamlOut = f.YAML
	return textOut, compactOut, jsonOut, yamlOut, nil
}
