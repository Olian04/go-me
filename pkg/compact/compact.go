// Package compact implements the v1 slash-separated fingerprint format.
package compact

import (
	"strings"

	"github.com/Olian04/go-me/pkg/identity/model"
)

// Slots returns the six fixed compact slots for the payload.
func Slots(p *model.Payload) [6]string {
	var s [6]string
	if p == nil {
		return s
	}
	s[0] = platformScope(p)
	s[1] = hostScope(p)
	s[2] = accountScope(p)
	s[3] = principalScope(p)
	s[4] = groupScope(p)
	s[5] = contextScope(p)
	for i := range s {
		s[i] = normalizeSegment(s[i])
	}
	return s
}

// String joins slots with slashes (five separators).
func String(p *model.Payload) string {
	s := Slots(p)
	return strings.Join(s[:], "/")
}

// StrictRequiredMissing returns true when strict compact policy should fail (no primary id).
func StrictRequiredMissing(p *model.Payload) bool {
	if p == nil {
		return true
	}
	s := Slots(p)
	// Require account or principal slot for a minimal identity fingerprint.
	return s[2] == "" && s[3] == ""
}

func platformScope(p *model.Payload) string {
	base := sysInfoPlatform(p)
	if base == "" {
		return ""
	}
	// Optional hint from network workgroup (Windows-ish class).
	for _, src := range p.Sources {
		if src.Name == "network" {
			if d, ok := src.Data.(model.NetworkData); ok && d.Workgroup != "" {
				return base + ":" + d.Workgroup
			}
		}
	}
	return base
}

func sysInfoPlatform(p *model.Payload) string {
	for _, src := range p.Sources {
		if src.Name != "sysinfo" {
			continue
		}
		d, ok := src.Data.(model.SysInfoData)
		if !ok {
			continue
		}
		if d.Platform != "" {
			return d.Platform
		}
	}
	return ""
}

func hostScope(p *model.Payload) string {
	if p.Meta.Hostname != "" {
		if strings.Contains(p.Meta.Hostname, ".") {
			return p.Meta.Hostname
		}
	}
	for _, src := range p.Sources {
		if src.Name == "network" {
			if d, ok := src.Data.(model.NetworkData); ok {
				if d.FQDN != "" {
					return d.FQDN
				}
				if d.Hostname != "" {
					return d.Hostname
				}
			}
		}
	}
	if p.Meta.Hostname != "" {
		return p.Meta.Hostname
	}
	return ""
}

func accountScope(p *model.Payload) string {
	if p.Subject.Username != "" {
		return p.Subject.Username
	}
	for _, src := range p.Sources {
		if src.Name == "osaccount" {
			if d, ok := src.Data.(model.OsAccountData); ok && d.Username != "" {
				return d.Username
			}
		}
	}
	return ""
}

func principalScope(p *model.Payload) string {
	if p.Subject.UID != "" {
		return p.Subject.UID
	}
	for _, src := range p.Sources {
		if src.Name == "osaccount" {
			if d, ok := src.Data.(model.OsAccountData); ok && d.UID != "" {
				return d.UID
			}
		}
	}
	return ""
}

func groupScope(p *model.Payload) string {
	if p.Subject.GID != "" {
		return p.Subject.GID
	}
	for _, src := range p.Sources {
		if src.Name == "osaccount" {
			if d, ok := src.Data.(model.OsAccountData); ok && d.GID != "" {
				return d.GID
			}
		}
	}
	return ""
}

func contextScope(p *model.Payload) string {
	for _, src := range p.Sources {
		if src.Name != "envcontext" {
			continue
		}
		d, ok := src.Data.(model.EnvContextData)
		if !ok {
			continue
		}
		if d.SudoUser != "" {
			return "sudo:" + d.SudoUser
		}
		if d.SSHUser != "" {
			return "ssh:" + d.SSHUser
		}
		if d.CI != nil && d.CI.IsCI {
			a := d.CI.Actor
			if a == "" {
				a = d.CI.Provider
			}
			return "ci:" + a
		}
	}
	return ""
}

// normalizeSegment trims and escapes slashes. It does not change letter case: usernames,
// hostnames, ARNs, and other slot values are not universally case-insensitive across OSes.
func normalizeSegment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "/", "%2f")
	return s
}

// EscapeSegment is exported for tests documenting the slash escape contract.
func EscapeSegment(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "/", "%2f")
}

// ErrStrictCompact is returned when strict mode cannot emit compact output.
type ErrStrictCompact struct {
	Msg string
}

func (e ErrStrictCompact) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "compact output unavailable in strict mode (missing required slots)"
}

// FormatOrStrict returns the compact string or a strict-mode error when slots are insufficient.
func FormatOrStrict(p *model.Payload, strict bool) (string, error) {
	if strict && StrictRequiredMissing(p) {
		return "", ErrStrictCompact{Msg: "missing account or principal scope"}
	}
	return String(p), nil
}
