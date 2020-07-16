package filer

import "time"

type File struct {
	Key        FileKey
	Value      FileValue
	Desc       FileDesc
	CreateTime time.Time
}

type FileKey []byte
type FileValue []byte
type FileDesc []byte
