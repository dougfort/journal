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

	var readResults []journal.ReadItem
	for item := range NewReader(&buffer) {
		if item.ItemType == journal.Error {
			t.Fatalf("error reading results: %v", item.Item)
		}
		readResults = append(readResults, item)
	}

	if len(readResults) != 1 {
		t.Fatalf("expected 1 result; found %d", len(readResults))
	}

	if readResults[0].ItemType != journal.Create {
		t.Fatalf("unexpected item type %v", readResults[0].ItemType)
	}
}
