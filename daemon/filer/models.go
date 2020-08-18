package filer

import "time"

type File struct {
	Key        string
	Value      string
	Desc       string
	CreateTime time.Time
	FileMD5    string
	FileSize   int64
}

type FileKey []byte
type FileValue []byte
type FileDesc []byte

type FileMeta struct {
	FileMD5  string
	FileSize int64
}
