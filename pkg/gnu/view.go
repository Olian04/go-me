package gnu

import (
	"os/user"

	"github.com/Olian04/go-me/pkg/identity/model"
)

// IDView is data needed for GNU-compatible id output.
type IDView struct {
	EUID   string
	EGID   string
	EUser  string
	EGroup string
	RUID   string
	RGID   string
	RUser  string
	RGroup string
	// GroupIDs lists all group ids (numeric strings) when known.
	GroupIDs []string
}

// IDOptions is the v1 supported GNU id flag subset.
type IDOptions struct {
	User   bool
	Group  bool
	Groups bool
	Name   bool
	Real   bool
}

// BuildIDView derives a view from the canonical payload.
func BuildIDView(p *model.Payload) IDView {
	var v IDView
	if p == nil {
		return v
	}
	v.EUID = p.Subject.UID
	v.EGID = p.Subject.GID
	v.EUser = p.Subject.Username

	for _, s := range p.Sources {
		if s.Name == "osaccount" {
			if d, ok := s.Data.(model.OsAccountData); ok {
				if v.EUID == "" {
					v.EUID = d.UID
				}
				if v.EGID == "" {
					v.EGID = d.GID
				}
				if v.EUser == "" {
					v.EUser = d.Username
				}
				v.GroupIDs = append([]string(nil), d.GroupIDs...)
			}
		}
	}

	if v.EGroup == "" && v.EGID != "" {
		if g, err := user.LookupGroupId(v.EGID); err == nil {
			v.EGroup = g.Name
		}
	}
	if v.EUser == "" && v.EUID != "" {
		if u, err := user.LookupId(v.EUID); err == nil {
			v.EUser = u.Username
		}
	}

	v.RUID = v.EUID
	v.RGID = v.EGID
	v.RUser = v.EUser
	v.RGroup = v.EGroup

	for _, s := range p.Sources {
		if s.Name != "envcontext" {
			continue
		}
		d, ok := s.Data.(model.EnvContextData)
		if !ok {
			continue
		}
		if d.SudoUID != "" && d.SudoUID != v.EUID {
			v.RUID = d.SudoUID
			if d.SudoUser != "" {
				v.RUser = d.SudoUser
			} else if u, err := user.LookupId(d.SudoUID); err == nil {
				v.RUser = u.Username
			}
		}
	}

	if v.RGroup == "" && v.RGID != "" {
		if g, err := user.LookupGroupId(v.RGID); err == nil {
			v.RGroup = g.Name
		}
	}

	return v
}

// FormatWhoami returns the effective username (GNU whoami).
func FormatWhoami(p *model.Payload) string {
	v := BuildIDView(p)
	if v.EUser != "" {
		return v.EUser
	}
	return ""
}

// FormatID renders GNU id-compatible text for the v1 flag subset.
func FormatID(p *model.Payload, opt IDOptions) string {
	v := BuildIDView(p)
	any := opt.User || opt.Group || opt.Groups
	if !any {
		return formatIDDefault(v)
	}
	var parts []string
	if opt.User {
		parts = append(parts, formatUser(v, opt))
	}
	if opt.Group {
		parts = append(parts, formatGroup(v, opt))
	}
	if opt.Groups {
		parts = append(parts, formatGroups(v, opt))
	}
	out := ""
	for i, s := range parts {
		if i > 0 {
			out += " "
		}
		out += s
	}
	return out
}

func formatIDDefault(v IDView) string {
	u := formatIDPair("uid", v.EUID, v.EUser)
	g := formatIDPair("gid", v.EGID, v.EGroup)
	gs := formatGroupsList(v)
	return u + " " + g + " groups=" + gs
}

func formatIDPair(kind, id, name string) string {
	if name != "" {
		return kind + "=" + id + "(" + name + ")"
	}
	return kind + "=" + id
}

func formatGroupsList(v IDView) string {
	if len(v.GroupIDs) == 0 {
		return ""
	}
	var b []string
	for _, gid := range v.GroupIDs {
		n := ""
		if g, err := user.LookupGroupId(gid); err == nil {
			n = g.Name
		}
		if n != "" {
			b = append(b, gid+"("+n+")")
		} else {
			b = append(b, gid)
		}
	}
	return joinSpace(b)
}

func joinSpace(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += " "
		}
		out += s
	}
	return out
}

func formatUser(v IDView, opt IDOptions) string {
	id := pickUID(v, opt)
	name := pickUName(v, opt)
	if opt.Name {
		if name != "" {
			return name
		}
		return id
	}
	return id
}

func formatGroup(v IDView, opt IDOptions) string {
	id := pickGID(v, opt)
	name := pickGName(v, opt)
	if opt.Name {
		if name != "" {
			return name
		}
		return id
	}
	return id
}

func formatGroups(v IDView, opt IDOptions) string {
	ids := v.GroupIDs
	if opt.Name {
		var names []string
		for _, gid := range ids {
			if g, err := user.LookupGroupId(gid); err == nil {
				names = append(names, g.Name)
			} else {
				names = append(names, gid)
			}
		}
		return joinSpace(names)
	}
	return joinSpace(ids)
}

func pickUID(v IDView, opt IDOptions) string {
	if opt.Real {
		return v.RUID
	}
	return v.EUID
}

func pickGID(v IDView, opt IDOptions) string {
	if opt.Real {
		return v.RGID
	}
	return v.EGID
}

func pickUName(v IDView, opt IDOptions) string {
	if opt.Real {
		return v.RUser
	}
	return v.EUser
}

func pickGName(v IDView, opt IDOptions) string {
	if opt.Real {
		return v.RGroup
	}
	return v.EGroup
}
