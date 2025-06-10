package models

type File interface {
	GetContent() []byte
	GetPath() string
	GetChecksum() string
}
