// Package provider defines the identity source plugin interface.
package provider

import (
	"context"

	"github.com/Olian04/go-me/pkg/identity/model"
)

// Provider resolves one named identity source.
type Provider interface {
	Name() string
	Run(ctx context.Context) Result
}

// Result is the outcome of a single provider run.
type Result struct {
	Envelope model.SourceEnvelope
	// SubjectPatch merges into the canonical subject when non-nil (typically osaccount only).
	SubjectPatch *model.Subject
}
