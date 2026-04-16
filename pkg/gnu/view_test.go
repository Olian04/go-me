package gnu

import (
	"strings"
	"testing"

	"github.com/Olian04/go-me/pkg/identity/model"
)

func TestFormatIDDefaultShape(t *testing.T) {
	p := &model.Payload{
		Subject: model.Subject{
			Username: "alice",
			UID:      "1000",
			GID:      "1000",
		},
		Sources: []model.SourceEnvelope{
			{
				Name:   "osaccount",
				Status: model.StatusOK,
				Data: model.OsAccountData{
					Username: "alice",
					UID:      "1000",
					GID:      "1000",
					GroupIDs: []string{"1000", "27"},
				},
			},
		},
	}
	out := FormatID(p, IDOptions{})
	if !strings.Contains(out, "uid=") || !strings.Contains(out, "gid=") || !strings.Contains(out, "groups=") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestFormatIDUserNumeric(t *testing.T) {
	p := &model.Payload{
		Subject: model.Subject{UID: "1000", Username: "alice"},
	}
	out := FormatID(p, IDOptions{User: true})
	if out != "1000" {
		t.Fatalf("got %q", out)
	}
}
