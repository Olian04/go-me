package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
)

// argv is a minimal argv slice for tests (urfave/cli v3 expects argv[0] like os.Args).
func argv(args ...string) []string {
	return append([]string{"me"}, args...)
}

func TestRootVersionFlagShowsVersion(t *testing.T) {
	cmd := newApp()
	var buf bytes.Buffer
	cmd.Writer = &buf
	cmd.ErrWriter = &buf
	err := cmd.Run(context.Background(), argv("--version"))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Version") || !strings.Contains(out, "Revision") {
		t.Fatalf("expected version text columns in output, got %q", out)
	}
}

func TestVersionWithJSONOutput(t *testing.T) {
	cmd := newApp()
	var buf bytes.Buffer
	cmd.Writer = &buf
	cmd.ErrWriter = &buf
	err := cmd.Run(context.Background(), argv("--version", "--json"))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(out, "{") || !strings.Contains(out, `"Version"`) {
		t.Fatalf("expected JSON build metadata, got %q", out)
	}
}

func TestHideVersionNotForcedWhenVersionSet(t *testing.T) {
	cmd := newApp()
	if cmd.HideVersion {
		t.Fatal("HideVersion should be false when Version is set")
	}
}

func TestVisibleFlagsAfterSetup(t *testing.T) {
	cmd := newApp()
	var buf bytes.Buffer
	cmd.Writer = &buf
	cmd.ErrWriter = &buf
	if err := cmd.Run(context.Background(), argv("--help")); err != nil {
		t.Fatal(err)
	}
	n := len(cmd.VisibleFlags())
	if n < 9 {
		t.Fatalf("expected at least help+version+our flags, got %d", n)
	}
}

func TestMutuallyExclusiveOutputModes(t *testing.T) {
	cmd := newApp()
	var buf bytes.Buffer
	cmd.Writer = &buf
	cmd.ErrWriter = &buf
	err := cmd.Run(context.Background(), argv("--json", "--yaml"))
	if err == nil {
		t.Fatal("expected error")
	}
	var ec cli.ExitCoder
	if !errors.As(err, &ec) {
		t.Fatalf("expected ExitCoder, got %T", err)
	}
	if ec.ExitCode() != exitUsage {
		t.Fatalf("exit code: got %d want %d", ec.ExitCode(), exitUsage)
	}
}
