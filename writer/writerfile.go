package writer

import (
	"fmt"
	"log"
	"os"
)

// file writer
type FileWriter struct {
	file *os.File
}

func (fw *FileWriter) WriteInformation(file ReferencedFile) {
	fw.file.WriteString(fmt.Sprintf("%s (%1.1f Gb)\n", file.FileName, float64(file.FileSize)/1024/1024/1024))
}

func (fw *FileWriter) NewDirectory(directory ReferencedDirectory) {
	fw.file.WriteString(fmt.Sprintf("\n*** %s ***\n", directory.SubDirectory))
}

func (fw *FileWriter) EndDirectory(directory ReferencedDirectory) {
}

func (fw *FileWriter) InitiateStore(store string) {
	targetFile, err := os.Create(store)
	if err != nil {
		log.Fatal(err)
	}
	fw.file = targetFile
}

func (fw *FileWriter) CloseStore() {
	fw.file.Close()
}
