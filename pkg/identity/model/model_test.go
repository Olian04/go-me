package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestJSONGoldenFixture(t *testing.T) {
	path := filepath.Join("testdata", "golden.json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var p Payload
	if err := json.Unmarshal(b, &p); err != nil {
		t.Fatal(err)
	}
	if p.Subject.Username != "alice" {
		t.Fatalf("username: %q", p.Subject.Username)
	}
	out, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	var p2 Payload
	if err := json.Unmarshal(out, &p2); err != nil {
		t.Fatal(err)
	}
	var sawStats bool
	for _, src := range p2.Sources {
		if src.Name != "sysinfo" {
			continue
		}
		b, err := json.Marshal(src.Data)
		if err != nil {
			t.Fatal(err)
		}
		var d SysInfoData
		if err := json.Unmarshal(b, &d); err != nil {
			t.Fatal(err)
		}
		if d.Platform != "linux" || d.Arch != "amd64" {
			t.Fatalf("sysinfo: %+v", d)
		}
		sawStats = true
		break
	}
	if !sawStats {
		t.Fatal("expected sysinfo source in golden")
	}
}

func TestYAMLMarshal(t *testing.T) {
	p := Payload{
		Subject: Subject{Username: "bob"},
		Meta:    Meta{BestEffort: true},
	}
	b, err := yaml.Marshal(&p)
	if err != nil {
		t.Fatal(err)
	}
	var p2 Payload
	if err := yaml.Unmarshal(b, &p2); err != nil {
		t.Fatal(err)
	}
	if p2.Subject.Username != "bob" {
		t.Fatal()
	}
}
