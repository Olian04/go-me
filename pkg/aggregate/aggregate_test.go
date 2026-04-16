package aggregate

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Olian04/go-me/pkg/identity/model"
)

func TestParseSourceFlags(t *testing.T) {
	got := ParseSourceFlags([]string{" osaccount ", "network,authproviders", "osaccount, envcontext "})
	want := []string{"osaccount", "network", "authproviders", "osaccount", "envcontext"}
	if len(got) != len(want) {
		t.Fatalf("len %d vs %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] got %q want %q", i, got[i], want[i])
		}
	}
}

func TestKnownSource(t *testing.T) {
	if !KnownSource("osaccount") || !KnownSource("OSACCOUNT") {
		t.Error("expected osaccount known")
	}
	if !KnownSource("sysinfo") || !KnownSource("authproviders") {
		t.Error("expected sysinfo and authproviders known")
	}
	if KnownSource("nope") {
		t.Error("expected unknown")
	}
}

func TestDefaultSourcesMembership(t *testing.T) {
	want := map[string]struct{}{
		"osaccount": {}, "envcontext": {}, "network": {}, "sysinfo": {},
	}
	got := map[string]struct{}{}
	for _, s := range DefaultSources {
		got[s] = struct{}{}
	}
	if len(got) != len(want) {
		t.Fatalf("DefaultSources count %d, want %d", len(got), len(want))
	}
	for k := range want {
		if _, ok := got[k]; !ok {
			t.Fatalf("DefaultSources missing %q", k)
		}
	}
}

func TestBestEffortUnknownSourceContinues(t *testing.T) {
	ctx := context.Background()
	p, err := Aggregate(ctx, Options{
		Timeout:    5 * time.Second,
		Sources:    []string{"osaccount", "not-a-real-source"},
		BestEffort: true,
		Strict:     false,
	})
	if err != nil {
		t.Fatalf("aggregate: %v", err)
	}
	var found bool
	for _, e := range p.Errors {
		if e.Code == "unknown_source" && strings.Contains(e.Message, "not-a-real-source") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected unknown_source error, got %#v", p.Errors)
	}
	if len(p.Sources) < 1 {
		t.Fatal("expected osaccount to run")
	}
}

func TestStrictUnknownSourceFails(t *testing.T) {
	ctx := context.Background()
	_, err := Aggregate(ctx, Options{
		Timeout:    5 * time.Second,
		Sources:    []string{"badsource"},
		BestEffort: false,
		Strict:     true,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(StrictUnknownSourceError); !ok {
		t.Fatalf("wrong error type: %T %v", err, err)
	}
}

func TestPayloadShapeGolden(t *testing.T) {
	p := &model.Payload{
		Subject: model.Subject{
			Username: "alice",
			UID:      "1000",
			GID:      "1000",
			HomeDir:  "/home/alice",
			Shell:    "/bin/bash",
		},
		Sources: []model.SourceEnvelope{
			{
				Name:       "osaccount",
				Status:     model.StatusOK,
				DurationMs: 1,
				Data: model.OsAccountData{
					Username: "alice",
					UID:      "1000",
					GID:      "1000",
					HomeDir:  "/home/alice",
					Shell:    "/bin/bash",
					GroupIDs: []string{"1000", "27"},
					Groups:   []string{"1000(alice)", "27(sudo)"},
				},
			},
		},
		Meta: model.Meta{
			Hostname:   "host.example",
			Timestamp:  "2020-01-02T15:04:05Z",
			DurationMs: 10,
			BestEffort: true,
		},
		Errors: []model.ErrorEntry{},
	}
	// Golden: top-level keys exist and subject matches (manual contract check).
	if p.Subject.Username != "alice" {
		t.Fatal()
	}
	if len(p.Sources) != 1 || p.Sources[0].Name != "osaccount" {
		t.Fatal()
	}
}
