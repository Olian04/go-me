package compact

import (
	"strings"
	"testing"

	"github.com/Olian04/go-me/pkg/identity/model"
)

func TestSixSlotsFiveSlashes(t *testing.T) {
	p := &model.Payload{
		Meta: model.Meta{Hostname: "h"},
		Sources: []model.SourceEnvelope{
			{Name: "sysinfo", Data: model.SysInfoData{Platform: "linux"}},
		},
		Subject: model.Subject{
			Username: "u",
			UID:      "1",
			GID:      "2",
		},
	}
	s := String(p)
	if strings.Count(s, "/") != 5 {
		t.Fatalf("want 5 slashes, got %q", s)
	}
	parts := strings.Split(s, "/")
	if len(parts) != 6 {
		t.Fatalf("want 6 parts, got %d: %q", len(parts), s)
	}
}

func TestEmptySlotsPreserveSeparators(t *testing.T) {
	p := &model.Payload{
		Sources: []model.SourceEnvelope{
			{Name: "sysinfo", Data: model.SysInfoData{Platform: "darwin"}},
		},
		Subject: model.Subject{},
	}
	s := String(p)
	if strings.Count(s, "/") != 5 {
		t.Fatalf("got %q", s)
	}
}

func TestStrictMissingAccountAndPrincipal(t *testing.T) {
	p := &model.Payload{
		Sources: []model.SourceEnvelope{
			{Name: "sysinfo", Data: model.SysInfoData{Platform: "linux"}},
		},
	}
	if !StrictRequiredMissing(p) {
		t.Fatal("expected strict missing")
	}
	_, err := FormatOrStrict(p, true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEscapeSlash(t *testing.T) {
	if got := EscapeSegment("a/b"); got != "a%2fb" {
		t.Fatalf("got %q", got)
	}
}
