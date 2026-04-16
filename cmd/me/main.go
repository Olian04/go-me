package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := newApp()
	ctx := context.Background()
	if err := app.Run(ctx, os.Args); err != nil {
		cli.HandleExitCoder(err) // usually OsExits; only plain errors reach the next lines
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
