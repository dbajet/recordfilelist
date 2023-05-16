package writer

import (
	"path"
)

type ReferencedDirectory struct {
	Source       string
	SubDirectory string
}

func (ref ReferencedDirectory) Disk() string {
	return path.Base(ref.Source)
}

type ReferencedFile struct {
	Disk      string
	Directory string
	FileName  string
	FileSize  int64
}
type Writer interface {
	WriteInformation(file ReferencedFile)
	NewDirectory(directory ReferencedDirectory)
	EndDirectory(directory ReferencedDirectory)
	InitiateStore(disk string)
	CloseStore()
}
