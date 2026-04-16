package osaccount

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/Olian04/go-me/pkg/identity/model"
	"github.com/Olian04/go-me/pkg/identity/provider"
)

// Name is the CLI --source value for this provider.
const Name = "osaccount"

type osaccount struct{}

// New returns the osaccount provider.
func New() provider.Provider {
	return osaccount{}
}

func (osaccount) Name() string { return Name }

func (osaccount) Run(ctx context.Context) provider.Result {
	_ = ctx
	u, err := user.Current()
	if err != nil {
		return provider.Result{
			Envelope: model.SourceEnvelope{
				Name:       Name,
				Status:     model.StatusError,
				DurationMs: 0,
				Data:       model.OsAccountData{},
				Warnings:   []string{fmt.Sprintf("user.Current: %v", err)},
			},
		}
	}

	data := model.OsAccountData{
		Username: u.Username,
		UID:      u.Uid,
		GID:      u.Gid,
		HomeDir:  u.HomeDir,
		Shell:    shellForUser(u),
	}

	status := model.StatusOK
	if data.Username == "" || data.UID == "" {
		status = model.StatusPartial
	}

	if runtime.GOOS != "windows" {
		gids, err := u.GroupIds()
		if err == nil {
			data.GroupIDs = append([]string(nil), gids...)
			var names []string
			for _, gid := range gids {
				lookup, err := user.LookupGroupId(gid)
				if err != nil {
					names = append(names, gid)
					continue
				}
				names = append(names, fmt.Sprintf("%s(%s)", gid, lookup.Name))
			}
			data.Groups = names
		} else {
			data.Groups = nil
			if status == model.StatusOK {
				status = model.StatusPartial
			}
		}
	}

	subject := model.Subject{
		Username: data.Username,
		UID:      data.UID,
		GID:      data.GID,
		HomeDir:  data.HomeDir,
		Shell:    data.Shell,
	}

	return provider.Result{
		Envelope: model.SourceEnvelope{
			Name:   Name,
			Status: status,
			Data:   data,
		},
		SubjectPatch: &subject,
	}
}

func shellForUser(u *user.User) string {
	if s := os.Getenv("SHELL"); s != "" && runtime.GOOS != "windows" {
		return s
	}
	return ""
}

// NormalizeShell trims shell path for display.
func NormalizeShell(s string) string {
	return strings.TrimSpace(s)
}
