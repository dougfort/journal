package filejournal

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/deciphernow/gm-control-api/api"
	"github.com/deciphernow/gm-control-api/api/objecttype"

	"github.com/dougfort/journal"
)

// | action-type | data size | data .... |

type FileWriter struct {
	Writer io.Writer
}

func NewWriter(writer io.Writer) (*FileWriter, error) {
	return &FileWriter{Writer: writer}, nil
}

// Create
func (w *FileWriter) Create(otype objecttype.ObjectType, object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		return errors.Wrap(err, "Marshal")
	}
	return w.appendObject(journal.Create, otype.ID(), data)
}

// Modify
func (w *FileWriter) Modify(otype objecttype.ObjectType, object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		return errors.Wrap(err, "Marshal")
	}
	return w.appendObject(journal.Modify, otype.ID(), data)
}

// Delete
func (w *FileWriter) Delete(otype objecttype.ObjectType, key string) error {
	data, err := json.Marshal(object)
	if err != nil {
		return errors.Wrap(err, "Marshal")
	}
	return w.appendObject(journal.Modify, otype.ID(), data)
}

// Version
func (w *FileWriter) Version(semVer string) error {
	return errors.Errorf("not implemented")
}

func (w *FileWriter) appendObject(
	itemType journal.ReadItemType,
	otype int64,
	data []byte,
) error {
	var err error

	if err = binary.Write(w.Writer, binary.BigEndian, itemType); err != nil {
		return errors.Wrap(err, "binary.Write: itemType")
	}
	if err = binary.Write(w.Writer, binary.BigEndian, otype); err != nil {
		return errors.Wrap(err, "binary.Write: otype")
	}

	timestamp, err := time.Now().UTC().MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "MarshalBinary()")
	}
	if err = binary.Write(w.Writer, binary.BigEndian, uint32(len(timestamp))); err != nil {
		return errors.Wrap(err, "binary.Write: timestamp size")
	}
	if _, err = w.Writer.Write(timestamp); err != nil {
		return errors.Wrap(err, "binary.Write: timestamp")
	}

	if err = binary.Write(w.Writer, binary.BigEndian, uint32(len(data))); err != nil {
		return errors.Wrap(err, "binary.Write: data size")
	}
	_, err = w.Writer.Write(data)

	return err
}

func NewReader(reader io.Reader) journal.Reader {
	itemChan := make(chan journal.ReadItem)
	go func() {
		var err error
		var itemType uint32
		var otype int64
		var timeBuffer []byte
		var timestampSize uint32
		var dataBuffer []byte
		var dataSize uint32

		defer close(itemChan)
	READ_LOOP:
		for {
			if err = binary.Read(reader, binary.BigEndian, &itemType); err != nil {
				// this is the expected place to detect EOF
				if err != io.EOF {
					itemChan <- journal.ReadItem{
						ItemType: journal.Error,
						Item:     err,
					}
				}
				break READ_LOOP
			}
			if err = binary.Read(reader, binary.BigEndian, &otype); err != nil {
				itemChan <- journal.ReadItem{
					ItemType: journal.Error,
					Item:     err,
				}
				break READ_LOOP
			}
			if err = binary.Read(reader, binary.BigEndian, &timestampSize); err != nil {
				itemChan <- journal.ReadItem{
					ItemType: journal.Error,
					Item:     err,
				}
				break READ_LOOP
			}
			timeBuffer = make([]byte, timestampSize)
			if _, err = io.ReadFull(reader, timeBuffer); err != nil {
				itemChan <- journal.ReadItem{
					ItemType: journal.Error,
					Item:     err,
				}
				break READ_LOOP
			}

			if err = binary.Read(reader, binary.BigEndian, &dataSize); err != nil {
				itemChan <- journal.ReadItem{
					ItemType: journal.Error,
					Item:     err,
				}
				break READ_LOOP
			}
			dataBuffer = make([]byte, dataSize)
			if _, err = io.ReadFull(reader, dataBuffer); err != nil {
				itemChan <- journal.ReadItem{
					ItemType: journal.Error,
					Item:     err,
				}
				break READ_LOOP
			}

			item := constructItem(
				journal.ReadItemType(itemType),
				otype,
				timeBuffer,
				dataBuffer,
			)
			itemChan <- item
			if item.ItemType == journal.Error {
				break READ_LOOP
			}
		}
	}()

	return itemChan
}

func constructItem(
	itemType journal.ReadItemType,
	otype int64,
	timeBuffer []byte,
	dataBuffer []byte,
) journal.ReadItem {
	result := journal.ReadItem{ItemType: itemType}
	var err error

	if err = result.Timestamp.UnmarshalBinary(timeBuffer); err != nil {
		result.ItemType = journal.Error
		result.Item = errors.Wrap(err, "Timestamp.UnmarshalBinary")
		return result
	}

	switch itemType {
	case journal.Create, journal.Modify, journal.Delete:
		if result.Item, err = constructObject(otype, dataBuffer); err != nil {
			result.ItemType = journal.Error
			result.Item = errors.Wrap(err, "constructObject")
		}
	case journal.Version:
		result.Item = string(dataBuffer)
	default:
		result.ItemType = journal.Error
		result.Item = errors.Errorf("unknown itemType: %v", itemType)
	}

	return result
}

func constructObject(otype int64, data []byte) (interface{}, error) {
	var objectType objecttype.ObjectType
	var object interface{}
	var err error

	if objectType, err = objecttype.FromID(int(otype)); err != nil {
		return nil, errors.Wrapf(err, "objecttype.FromID(%v)", otype)
	}

	switch objectType {
	case objecttype.Zone:
		var zone api.Zone
		if err = json.Unmarshal(data, &zone); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal(%v)", objectType)
		}
		object = zone
	case objecttype.Proxy:
		var proxy api.Proxy
		if err = json.Unmarshal(data, &proxy); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal(%v)", objectType)
		}
		object = proxy
	case objecttype.Domain:
		var domain api.Domain
		if err = json.Unmarshal(data, &domain); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal(%v)", objectType)
		}
		object = domain
	case objecttype.Route:
		var route api.Route
		if err = json.Unmarshal(data, &route); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal(%v)", objectType)
		}
		object = route
	case objecttype.Cluster:
		var cluster api.Cluster
		if err = json.Unmarshal(data, &cluster); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal(%v)", objectType)
		}
		object = cluster
	case objecttype.SharedRules:
		var rules api.SharedRules
		if err = json.Unmarshal(data, &rules); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal(%v)", objectType)
		}
		object = rules
	case objecttype.Listener:
		var listener api.Listener
		if err = json.Unmarshal(data, &listener); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal(%v)", objectType)
		}
		object = listener
	default:
		return nil, errors.Errorf("unknown objecttype: %v", objectType)
	}

	return object, nil
}
