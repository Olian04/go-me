//go:build windows

package sysinfo

import (
	"os/exec"
	"strings"
)

// NameAndVersion returns a human-friendly OS name and version string when detectable.
func NameAndVersion() (name, version string) {
	return windowsCaption()
}

func windowsCaption() (name, version string) {
	out, err := exec.Command("cmd", "/c", "ver").CombinedOutput()
	if err != nil {
		return "", ""
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return "", ""
	}
	if i := strings.Index(s, "Version"); i >= 0 {
		rest := strings.TrimSpace(s[i+len("Version"):])
		rest = strings.Trim(rest, "[]")
		return "Windows", rest
	}
	return "Windows", s
}
