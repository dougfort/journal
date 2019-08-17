package journal

import (
	"github.com/deciphernow/gm-control-api/api/objecttype"
)

// Writer maintains an append-only journal of gm-control-api transactions
type Writer interface {

	// Create
	Create(otype objecttype.ObjectType, object interface{}) error

	// Modify
	Modify(otype objecttype.ObjectType, object interface{}) error

	// Delete
	Delete(otype objecttype.ObjectType, key string) error

	// Version
	Version(semVer string) error
}

// ReadItemType is the type of item that can come from a reader
type ReadItemType uint32

const (
	_ ReadItemType = iota
	Error
	Create
	Modify
	Delete
	Version
)

// ReadItem is an individual entry parsed from a Journal
type ReadItem struct {
	ItemType ReadItemType
	Item     interface{}
}

// Reader is a channel that returns parsed items from a Journal
type Reader <-chan ReadItem

type ErrorItem error

type CreateItem struct {
	OType  objecttype.ObjectType
	Object interface{}
}

type ModifyItem struct {
	OType  objecttype.ObjectType
	Object interface{}
}

type DeleteItem struct {
	OType objecttype.ObjectType
	Key   string
}

type VersionItem string
