// Package aggregate orchestrates identity providers and merges a canonical payload.
package aggregate

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Olian04/go-me/pkg/identity/authproviders"
	"github.com/Olian04/go-me/pkg/identity/envcontext"
	"github.com/Olian04/go-me/pkg/identity/model"
	"github.com/Olian04/go-me/pkg/identity/network"
	"github.com/Olian04/go-me/pkg/identity/osaccount"
	"github.com/Olian04/go-me/pkg/identity/provider"
	"github.com/Olian04/go-me/pkg/identity/sysinfo"
)

// DefaultSources is the provider order when --source is omitted.
var DefaultSources = []string{
	osaccount.Name,
	envcontext.Name,
	network.Name,
	sysinfo.Name,
}

// SourceFlagUsage is the full --source flag description for CLI help and packaging docs.
const SourceFlagUsage = "Identity providers to run, in order (repeatable or comma-separated). " +
	"When omitted, runs: osaccount, envcontext, network, sysinfo. " +
	"When any --source is present, it replaces that full list entirely (not additive). " +
	"Valid names: osaccount, envcontext, network, sysinfo, authproviders."

// KnownSource reports whether name is a valid provider id.
func KnownSource(name string) bool {
	name = strings.TrimSpace(strings.ToLower(name))
	switch name {
	case osaccount.Name, envcontext.Name, network.Name, sysinfo.Name, authproviders.Name:
		return true
	default:
		return false
	}
}

// ParseSourceFlags flattens repeatable and comma-separated --source values.
func ParseSourceFlags(parts []string) []string {
	var out []string
	for _, p := range parts {
		for _, s := range strings.Split(p, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				out = append(out, strings.ToLower(s))
			}
		}
	}
	return out
}

// Options controls aggregate behavior.
type Options struct {
	Timeout    time.Duration
	Sources    []string // empty = DefaultSources
	BestEffort bool
	Strict     bool
}

// Aggregate runs selected providers and returns the canonical payload.
func Aggregate(ctx context.Context, opts Options) (*model.Payload, error) {
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	start := time.Now()
	sources := opts.Sources
	if len(sources) == 0 {
		sources = append([]string(nil), DefaultSources...)
	}

	var errs []model.ErrorEntry
	var unknown []string
	resolved := make([]string, 0, len(sources))
	for _, s := range sources {
		if KnownSource(s) {
			resolved = append(resolved, s)
			continue
		}
		unknown = append(unknown, s)
		errs = append(errs, model.ErrorEntry{
			Source:  s,
			Code:    "unknown_source",
			Message: fmt.Sprintf("unknown identity source %q", s),
		})
	}

	if len(unknown) > 0 && opts.Strict {
		return nil, StrictUnknownSourceError{Names: unknown}
	}

	registry := map[string]provider.Provider{
		osaccount.Name:     osaccount.New(),
		envcontext.Name:    envcontext.New(),
		network.Name:       network.New(),
		sysinfo.Name:       sysinfo.New(),
		authproviders.Name: authproviders.New(),
	}

	payload := &model.Payload{
		Errors: errs,
		Meta: model.Meta{
			Timestamp:  model.NowRFC3339(),
			BestEffort: opts.BestEffort || !opts.Strict,
		},
	}

	if host, err := os.Hostname(); err == nil {
		payload.Meta.Hostname = host
	}

	var subject model.Subject

	for _, name := range resolved {
		p, ok := registry[name]
		if !ok {
			continue
		}
		pstart := time.Now()
		pctx := ctx
		res := p.Run(pctx)
		dur := time.Since(pstart).Milliseconds()
		res.Envelope.DurationMs = dur
		payload.Sources = append(payload.Sources, res.Envelope)

		if res.SubjectPatch != nil {
			subject = mergeSubject(subject, *res.SubjectPatch)
		}

		if res.Envelope.Status == model.StatusError && opts.Strict {
			return nil, StrictProviderError{Name: name, Message: firstWarning(res.Envelope.Warnings)}
		}
		if res.Envelope.Status == model.StatusError && opts.BestEffort {
			payload.Errors = append(payload.Errors, model.ErrorEntry{
				Source:  name,
				Code:    "provider_error",
				Message: firstWarning(res.Envelope.Warnings),
			})
		}
	}

	payload.Subject = subject
	if payload.Meta.Hostname == "" {
		for _, s := range payload.Sources {
			if s.Name == network.Name {
				if d, ok := s.Data.(model.NetworkData); ok && d.Hostname != "" {
					payload.Meta.Hostname = d.Hostname
					break
				}
			}
		}
	}

	payload.Meta.DurationMs = time.Since(start).Milliseconds()
	payload.Meta.BestEffort = opts.BestEffort || !opts.Strict

	return payload, nil
}

func firstWarning(w []string) string {
	if len(w) == 0 {
		return "provider failed"
	}
	return w[0]
}

func mergeSubject(base, patch model.Subject) model.Subject {
	out := base
	if patch.Username != "" {
		out.Username = patch.Username
	}
	if patch.DisplayName != "" {
		out.DisplayName = patch.DisplayName
	}
	if patch.UID != "" {
		out.UID = patch.UID
	}
	if patch.GID != "" {
		out.GID = patch.GID
	}
	if patch.HomeDir != "" {
		out.HomeDir = patch.HomeDir
	}
	if patch.Shell != "" {
		out.Shell = patch.Shell
	}
	return out
}
