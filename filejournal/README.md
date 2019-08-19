# filejouranl

## synopsis
This is a flat file based system for backing up transactions to the
gm-control-api back end. 

## data format

### object format
|  int32    | int64       | int32          |  []byte   | int32     | []byte |
| item type | object type | timestamp size | timestamp | data size | data   | 

### semantic version format
| int32         | []byte  |
| versiion size | version |