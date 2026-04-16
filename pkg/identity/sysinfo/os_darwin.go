//go:build darwin

package sysinfo

import (
	"os/exec"
	"strings"
)

// NameAndVersion returns a human-friendly OS name and version string when detectable.
func NameAndVersion() (name, version string) {
	return darwinSwVers()
}

func darwinSwVers() (name, version string) {
	verOut, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return "macOS", ""
	}
	version = strings.TrimSpace(string(verOut))
	nameOut, err := exec.Command("sw_vers", "-productName").Output()
	if err != nil {
		return "macOS", version
	}
	n := strings.TrimSpace(string(nameOut))
	if n == "" {
		n = "macOS"
	}
	return n, version
}
