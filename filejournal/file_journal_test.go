package filejournal

import (
	"bytes"
	"testing"

	"github.com/deciphernow/gm-control-api/api"
	"github.com/deciphernow/gm-control-api/api/objecttype"

	"github.com/dougfort/journal"
)

func TestFileJournal(t *testing.T) {
	var err error
	var writer journal.Writer
	var buffer bytes.Buffer

	writer, err = NewWriter(&buffer)
	if err != nil {
		t.Fatalf("NewWriter failed: %s", err)
	}

	c1 := api.Cluster{}
	err = writer.Create(objecttype.Cluster, c1)
	if err != nil {
		t.Fatalf("writer.Create( failed: %s", err)
	}
}
