//go:build !linux && !darwin && !windows

package sysinfo

// NameAndVersion returns a human-friendly OS name and version string when detectable.
// Either value may be empty on unsupported platforms or when detection fails.
func NameAndVersion() (name, version string) {
	return "", ""
}
