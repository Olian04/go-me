package authproviders

import (
	"context"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Olian04/go-me/pkg/identity/model"
	"github.com/Olian04/go-me/pkg/identity/provider"
)

// Name is the CLI --source value for this provider.
const Name = "authproviders"

type authproviders struct{}

// New returns the authproviders provider.
func New() provider.Provider {
	return authproviders{}
}

func (authproviders) Name() string { return Name }

func (authproviders) Run(ctx context.Context) provider.Result {
	_ = ctx
	data := model.AuthProvidersData{Git: &model.GitAuthData{}}

	if name, email, ok := gitGlobalUser(ctx); ok {
		data.Git.UserName = name
		data.Git.UserEmail = email
	}

	data.Cloud = &model.CloudAuthData{}

	status := model.StatusPartial
	if data.Git.UserName != "" || data.Git.UserEmail != "" {
		status = model.StatusOK
	}

	return provider.Result{
		Envelope: model.SourceEnvelope{
			Name:   Name,
			Status: status,
			Data:   data,
		},
	}
}

func gitGlobalUser(ctx context.Context) (name, email string, ok bool) {
	path := gitConfigPath()
	if path != "" {
		// #nosec G304 -- path is ~/.gitconfig under os/user.Current().HomeDir, not user-controlled.
		if b, err := os.ReadFile(path); err == nil {
			name, email = parseGitConfig(string(b))
			if name != "" || email != "" {
				return name, email, true
			}
		}
	}
	return gitFromExec(ctx)
}

func gitConfigPath() string {
	home, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.Join(home.HomeDir, ".gitconfig")
}

func parseGitConfig(s string) (name, email string) {
	var inUser bool
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sec := strings.ToLower(strings.Trim(line, "[]"))
			inUser = sec == "user"
			continue
		}
		if !inUser {
			continue
		}
		if strings.HasPrefix(strings.ToLower(line), "name") {
			name = strings.Trim(strings.TrimPrefix(line, "name"), " \t=")
			name = strings.Trim(name, `"`)
		}
		if strings.HasPrefix(strings.ToLower(line), "email") {
			email = strings.Trim(strings.TrimPrefix(line, "email"), " \t=")
			email = strings.Trim(email, `"`)
		}
	}
	return name, email
}

func gitFromExec(ctx context.Context) (name, email string, ok bool) {
	path, err := exec.LookPath("git")
	if err != nil {
		return "", "", false
	}
	// #nosec G204 -- path is from exec.LookPath("git"), not user-controlled input.
	c := exec.CommandContext(ctx, path, "config", "--global", "user.name")
	out, err := c.Output()
	if err == nil {
		name = strings.TrimSpace(string(out))
	}
	// #nosec G204 -- path is from exec.LookPath("git"), not user-controlled input.
	c2 := exec.CommandContext(ctx, path, "config", "--global", "user.email")
	out2, err := c2.Output()
	if err == nil {
		email = strings.TrimSpace(string(out2))
	}
	if name != "" || email != "" {
		return name, email, true
	}
	return "", "", false
}
