package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"recordfilelist/writer"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/exp/slices"
)

func getFileList(writers []writer.Writer, directories []writer.ReferencedDirectory) []writer.ReferencedDirectory {
	invalidDirectories := []string{"/$RECYCLE.BIN/", "/System Volume Information/", "/lost+found/", "/.Trash-1000/"}
	invalidFiles := []string{".AmazonDriveConfig"}

	directory := directories[0]
	if !slices.Contains(invalidDirectories, directory.SubDirectory) {
		entries, err := os.ReadDir(fmt.Sprintf("%s%s", directory.Source, directory.SubDirectory))
		if err != nil {
			log.Fatal(err)
		}
		for _, aWriter := range writers {
			aWriter.NewDirectory(directory)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				directories = append(directories, writer.ReferencedDirectory{Source: directory.Source, SubDirectory: fmt.Sprintf("%s%s/", directory.SubDirectory, entry.Name())})
			} else if !slices.Contains(invalidFiles, entry.Name()) {
				fileInfo, _ := entry.Info()
				for _, aWriter := range writers {
					// fmt.Println(entry.Name(), errInfo)
					aWriter.WriteInformation(writer.ReferencedFile{
						Disk:      directory.Disk(),
						Directory: directory.SubDirectory,
						FileName:  strings.ToUpper(entry.Name()),
						FileSize:  fileInfo.Size(),
					})
				}
			}
		}
		for _, aWriter := range writers {
			aWriter.EndDirectory(directory)
		}
	}
	return directories[1:]
}
func main() {
	help := `Not enough or invalid arguments:
	--source directory from
	--target directory to
	--credentials file path to the DB credentials
	`
	var source, target, credentials string
	for idx := 1; idx < len(os.Args); idx += 2 {
		switch parameter := os.Args[idx]; parameter {
		case "--source":
			source = strings.TrimSuffix(os.Args[idx+1], "/")
		case "--target":
			target = strings.TrimSuffix(os.Args[idx+1], "/")
		case "--credentials":
			credentials = strings.TrimSuffix(os.Args[idx+1], "/")
		}
	}
	if source == "" || target == "" || credentials == "" {
		log.Fatal(help)
	}
	// file writer
	aFileWriter := writer.FileWriter{}
	aFileWriter.InitiateStore(fmt.Sprintf("%s/%s.txt", target, path.Base(source)))
	defer aFileWriter.CloseStore()
	// db writer
	aDbWriter := writer.DatabaseWriter{DbName: "forgogo", DbTable: "toremove", Credentials: credentials}
	aDbWriter.InitiateStore(path.Base(source))
	defer aDbWriter.CloseStore()
	// read directories without recursion
	directories := []writer.ReferencedDirectory{{Source: source, SubDirectory: "/"}}
	for {
		directories = getFileList([]writer.Writer{&aFileWriter, &aDbWriter}, directories)
		if len(directories) == 0 {
			break
		}
	}

}
