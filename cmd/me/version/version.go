// Package version exposes build metadata resolved from ldflags and runtime/debug.
package version

import (
	"fmt"
	"runtime/debug"
	"strings"
)

var (
	// Version is the release tag or semantic version (ldflags).
	Version = "unknown"
	// Revision is the VCS revision (ldflags).
	Revision = "unknown"
	// BuildTime is the build timestamp in UTC (ldflags).
	BuildTime = "unknown"
)

// Info holds resolved build metadata.
type Info struct {
	Version   string
	Revision  string
	BuildTime string
}

// Get resolves Version, Revision, and BuildTime in priority order:
// 1. Non-empty ldflags values other than "unknown"
// 2. runtime/debug.ReadBuildInfo settings (vcs.revision, vcs.time, vcs.tag / module version)
// 3. Fallback "unknown"
func Get() Info {
	i := Info{
		Version:   pick(Version),
		Revision:  pick(Revision),
		BuildTime: pick(BuildTime),
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return finalize(i)
	}

	if i.Version == "" || i.Version == "unknown" {
		if v := moduleVersion(bi); v != "" {
			i.Version = v
		}
	}

	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			if i.Revision == "" || i.Revision == "unknown" {
				i.Revision = s.Value
			}
		case "vcs.time":
			if i.BuildTime == "" || i.BuildTime == "unknown" {
				i.BuildTime = s.Value
			}
		case "vcs.tag":
			if i.Version == "" || i.Version == "unknown" {
				i.Version = s.Value
			}
		}
	}

	return finalize(i)
}

func moduleVersion(bi *debug.BuildInfo) string {
	if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}
	return ""
}

func pick(ld string) string {
	s := strings.TrimSpace(ld)
	if s == "" || s == "unknown" {
		return ""
	}
	return s
}

func finalize(i Info) Info {
	if i.Version == "" {
		i.Version = "unknown"
	}
	if i.Revision == "" {
		i.Revision = "unknown"
	}
	if i.BuildTime == "" {
		i.BuildTime = "unknown"
	}
	return i
}

// String returns a single-line description (compact, key=value) for logs and APIs.
func String() string {
	i := Get()
	return i.String()
}

// String implements fmt.Stringer for Info.
func (i Info) String() string {
	return fmt.Sprintf("version=%s revision=%s build_time=%s", i.Version, i.Revision, i.BuildTime)
}
