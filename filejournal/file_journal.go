package filejournal

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"

	"github.com/deciphernow/gm-control-api/api/objecttype"

	"github.com/dougfort/journal"
)

// | action-type | data size | data .... |

type FileWriter struct {
	Writer io.Writer
}

type actionType uint32

const (
	_ actionType = iota
	Create
	Modify
	Delete
	Version
)

func NewWriter(writer io.Writer) (*FileWriter, error) {
	return &FileWriter{Writer: writer}, nil
}

// Create
func (w *FileWriter) Create(otype objecttype.ObjectType, object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		return errors.Wrap(err, "Marshal")
	}
	return w.appendObject(Create, otype.ID(), data)
}

// Modify
func (w *FileWriter) Modify(otype objecttype.ObjectType, object interface{}) error {
	return errors.Errorf("not implemented")
}

// Delete
func (w *FileWriter) Delete(otype objecttype.ObjectType, key string) error {
	return errors.Errorf("not implemented")
}

// Version
func (w *FileWriter) Version(semVer string) error {
	return errors.Errorf("not implemented")
}

func (w *FileWriter) appendObject(
	action actionType,
	otype int64,
	data []byte,
) error {

}

func NewReader() (<-chan journal.ReadItem, error) {
	return nil, errors.Errorf("not implemented")
}
