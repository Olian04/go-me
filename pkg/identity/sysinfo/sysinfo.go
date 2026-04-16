// Package sysinfo is an identity source for cheap host OS/runtime facts (GOOS, GOARCH, OS name/version).
package sysinfo

import (
	"context"
	"runtime"

	"github.com/Olian04/go-me/pkg/identity/model"
	"github.com/Olian04/go-me/pkg/identity/provider"
)

// Name is the CLI --source value for this provider.
const Name = "sysinfo"

type impl struct{}

// New returns the sysinfo provider.
func New() provider.Provider {
	return &impl{}
}

func (*impl) Name() string { return Name }

func (*impl) Run(ctx context.Context) provider.Result {
	osName, osVer := NameAndVersion()
	d := model.SysInfoData{
		Platform:  runtime.GOOS,
		Arch:      runtime.GOARCH,
		OSName:    osName,
		OSVersion: osVer,
	}
	return provider.Result{
		Envelope: model.SourceEnvelope{
			Name:   Name,
			Status: model.StatusOK,
			Data:   d,
		},
	}
}
