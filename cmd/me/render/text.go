package render

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/Olian04/go-me/cmd/me/version"
	"github.com/Olian04/go-me/pkg/identity/model"
)

// Text renders the default human identity summary per docs/design/cli-api.md.
// It includes subject, system (OS / architecture / platform), host, env/network/auth hints, and warnings—no run timing or strict-mode metadata.
func Text(p *model.Payload) string {
	if p == nil {
		return ""
	}
	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	s := p.Subject

	_, _ = fmt.Fprintf(tw, "Username\t%s\n", orUnknown(s.Username))
	if s.DisplayName != "" {
		_, _ = fmt.Fprintf(tw, "Display name\t%s\n", s.DisplayName)
	}
	_, _ = fmt.Fprintf(tw, "UID\t%s\n", orUnknown(s.UID))
	_, _ = fmt.Fprintf(tw, "GID\t%s\n", orUnknown(s.GID))
	_, _ = fmt.Fprintf(tw, "Home\t%s\n", orUnknown(s.HomeDir))
	_, _ = fmt.Fprintf(tw, "Shell\t%s\n", orUnknown(s.Shell))
	_, _ = fmt.Fprintf(tw, "Hostname\t%s\n", orUnknown(p.Meta.Hostname))

	for _, row := range sysInfoRows(p) {
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", row.label, row.value)
	}
	for _, row := range envContextRows(p) {
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", row.label, row.value)
	}
	for _, row := range networkExtraRows(p) {
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", row.label, row.value)
	}
	for _, row := range authProviderRows(p) {
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", row.label, row.value)
	}
	_ = tw.Flush()

	out := strings.TrimRight(b.String(), "\n")
	if len(p.Errors) == 0 {
		return out
	}

	var w strings.Builder
	w.WriteString(out)
	w.WriteString("\n\nWarnings:\n")
	for _, e := range p.Errors {
		_, _ = fmt.Fprintf(&w, "  - [%s] %s: %s\n", e.Source, e.Code, e.Message)
	}
	return strings.TrimRight(w.String(), "\n")
}

func orUnknown(s string) string {
	if strings.TrimSpace(s) == "" {
		return "<unknown>"
	}
	return s
}

type labelValue struct {
	label, value string
}

func sysInfoRows(p *model.Payload) []labelValue {
	for _, src := range p.Sources {
		if src.Name != "sysinfo" {
			continue
		}
		d, ok := src.Data.(model.SysInfoData)
		if !ok {
			continue
		}
		var rows []labelValue
		if d.OSName != "" {
			rows = append(rows, labelValue{"OS", d.OSName})
		} else if d.Platform != "" {
			rows = append(rows, labelValue{"OS", d.Platform})
		}
		if d.OSVersion != "" {
			rows = append(rows, labelValue{"OS version", d.OSVersion})
		}
		if d.Arch != "" {
			rows = append(rows, labelValue{"Architecture", d.Arch})
		}
		if d.OSName != "" && d.Platform != "" {
			rows = append(rows, labelValue{"Platform", d.Platform})
		}
		return rows
	}
	return nil
}

func envContextRows(p *model.Payload) []labelValue {
	var rows []labelValue
	for _, src := range p.Sources {
		if src.Name != "envcontext" {
			continue
		}
		d, ok := src.Data.(model.EnvContextData)
		if !ok {
			continue
		}
		if d.SudoUser != "" {
			rows = append(rows, labelValue{"Sudo user", d.SudoUser})
		}
		if d.SudoUID != "" {
			rows = append(rows, labelValue{"Sudo UID", d.SudoUID})
		}
		if d.SSHUser != "" {
			rows = append(rows, labelValue{"SSH user", d.SSHUser})
		}
		if d.CI != nil && d.CI.IsCI {
			if d.CI.Actor != "" {
				rows = append(rows, labelValue{"CI actor", d.CI.Actor})
			}
			if d.CI.Provider != "" {
				rows = append(rows, labelValue{"CI provider", d.CI.Provider})
			}
			if d.CI.Actor == "" && d.CI.Provider == "" {
				rows = append(rows, labelValue{"CI", "active"})
			}
		}
		return rows
	}
	return rows
}

func networkExtraRows(p *model.Payload) []labelValue {
	for _, src := range p.Sources {
		if src.Name != "network" {
			continue
		}
		d, ok := src.Data.(model.NetworkData)
		if !ok {
			continue
		}
		var rows []labelValue
		if d.FQDN != "" {
			rows = append(rows, labelValue{"FQDN", d.FQDN})
		}
		if d.Domain != "" {
			rows = append(rows, labelValue{"Domain", d.Domain})
		}
		if d.Workgroup != "" {
			rows = append(rows, labelValue{"Workgroup", d.Workgroup})
		}
		if len(rows) > 0 {
			return rows
		}
	}
	return nil
}

func authProviderRows(p *model.Payload) []labelValue {
	for _, src := range p.Sources {
		if src.Name != "authproviders" {
			continue
		}
		d, ok := src.Data.(model.AuthProvidersData)
		if !ok {
			continue
		}
		var rows []labelValue
		if d.Git != nil {
			if d.Git.UserName != "" {
				rows = append(rows, labelValue{"Git user", d.Git.UserName})
			}
			if d.Git.UserEmail != "" {
				rows = append(rows, labelValue{"Git email", d.Git.UserEmail})
			}
		}
		if d.Cloud != nil {
			if d.Cloud.AWS != nil {
				if d.Cloud.AWS.ARN != "" {
					rows = append(rows, labelValue{"AWS ARN", d.Cloud.AWS.ARN})
				}
				if d.Cloud.AWS.AccountID != "" {
					rows = append(rows, labelValue{"AWS account", d.Cloud.AWS.AccountID})
				}
			}
			if d.Cloud.GCP != nil {
				if d.Cloud.GCP.Account != "" {
					rows = append(rows, labelValue{"GCP account", d.Cloud.GCP.Account})
				}
				if d.Cloud.GCP.Project != "" {
					rows = append(rows, labelValue{"GCP project", d.Cloud.GCP.Project})
				}
			}
			if d.Cloud.Azure != nil {
				if d.Cloud.Azure.User != "" {
					rows = append(rows, labelValue{"Azure user", d.Cloud.Azure.User})
				}
				if d.Cloud.Azure.TenantID != "" {
					rows = append(rows, labelValue{"Azure tenant", d.Cloud.Azure.TenantID})
				}
				if d.Cloud.Azure.SubscriptionID != "" {
					rows = append(rows, labelValue{"Azure subscription", d.Cloud.Azure.SubscriptionID})
				}
			}
		}
		return rows
	}
	return nil
}

// VersionText formats resolved build metadata for human-readable --version output.
func VersionText(i version.Info) string {
	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(tw, "Version\t%s\n", i.Version)
	_, _ = fmt.Fprintf(tw, "Revision\t%s\n", i.Revision)
	_, _ = fmt.Fprintf(tw, "Built\t%s\n", i.BuildTime)
	_ = tw.Flush()
	return strings.TrimRight(b.String(), "\n")
}
