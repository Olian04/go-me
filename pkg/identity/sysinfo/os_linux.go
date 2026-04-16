//go:build linux

package sysinfo

import (
	"bufio"
	"os"
	"strings"
)

// NameAndVersion returns a human-friendly OS name and version string when detectable.
func NameAndVersion() (name, version string) {
	return linuxFromOSRelease()
}

func linuxFromOSRelease() (name, version string) {
	b, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", ""
	}
	m := parseEnvFile(string(b))
	if v := m["PRETTY_NAME"]; v != "" {
		name = unquote(v)
	} else if v := m["NAME"]; v != "" {
		name = unquote(v)
	}
	if v := m["VERSION_ID"]; v != "" {
		version = unquote(v)
	} else if v := m["VERSION"]; v != "" {
		version = unquote(v)
	}
	return name, version
}

func parseEnvFile(s string) map[string]string {
	out := make(map[string]string)
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		i := strings.IndexByte(line, '=')
		if i <= 0 {
			continue
		}
		k := strings.TrimSpace(line[:i])
		v := strings.TrimSpace(line[i+1:])
		out[k] = v
	}
	return out
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.Trim(s, `"`)
	}
	return s
}
