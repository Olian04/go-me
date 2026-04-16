package aggregate

import "strings"

// StrictUnknownSourceError is returned when --strict and an unknown --source is present.
type StrictUnknownSourceError struct {
	Names []string
}

func (e StrictUnknownSourceError) Error() string {
	return "unknown source(s): " + strings.Join(e.Names, ", ")
}

// StrictProviderError is returned when --strict and a provider reports StatusError.
type StrictProviderError struct {
	Name    string
	Message string
}

func (e StrictProviderError) Error() string {
	if e.Message != "" {
		return e.Name + ": " + e.Message
	}
	return e.Name + ": provider error"
}
