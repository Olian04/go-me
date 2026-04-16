package envcontext

import (
	"context"
	"os"
	"strings"

	"github.com/Olian04/go-me/pkg/identity/model"
	"github.com/Olian04/go-me/pkg/identity/provider"
)

// Name is the CLI --source value for this provider.
const Name = "envcontext"

type envcontext struct{}

// New returns the envcontext provider.
func New() provider.Provider {
	return envcontext{}
}

func (envcontext) Name() string { return Name }

func (envcontext) Run(ctx context.Context) provider.Result {
	_ = ctx
	data := model.EnvContextData{}
	if v := os.Getenv("SUDO_USER"); v != "" {
		data.SudoUser = v
	}
	if v := os.Getenv("SUDO_UID"); v != "" {
		data.SudoUID = v
	}
	if u := os.Getenv("SSH_USER"); u != "" {
		data.SSHUser = u
	} else if os.Getenv("SSH_CONNECTION") != "" || os.Getenv("SSH_CLIENT") != "" {
		if u := os.Getenv("USER"); u != "" {
			data.SSHUser = u
		}
	}

	if ci := detectCI(); ci != nil {
		data.CI = ci
	}

	status := model.StatusOK
	if data.SudoUser == "" && data.SudoUID == "" && data.SSHUser == "" && data.CI == nil {
		status = model.StatusPartial
	}

	return provider.Result{
		Envelope: model.SourceEnvelope{
			Name:   Name,
			Status: status,
			Data:   data,
		},
	}
}

func detectCI() *model.EnvCIData {
	if os.Getenv("CI") == "" && os.Getenv("GITHUB_ACTIONS") == "" && os.Getenv("GITLAB_CI") == "" {
		return nil
	}
	ci := &model.EnvCIData{IsCI: true}
	switch {
	case os.Getenv("GITHUB_ACTIONS") != "":
		ci.Provider = "github"
		if a := os.Getenv("GITHUB_ACTOR"); a != "" {
			ci.Actor = a
		}
	case os.Getenv("GITLAB_CI") != "":
		ci.Provider = "gitlab"
		if a := os.Getenv("GITLAB_USER_LOGIN"); a != "" {
			ci.Actor = a
		}
	case os.Getenv("CI") != "":
		ci.Provider = strings.ToLower(os.Getenv("CI"))
	}
	return ci
}
