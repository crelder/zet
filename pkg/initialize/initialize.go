package initialize

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

// Initiator satisfies the zet.Initiator interface.
type Initiator struct {
	path string
}

func New(path string) Initiator {
	return Initiator{
		path: path,
	}
}

// Init creates an empty zettelkasten containing an index.txt, references.bib and a folder 'zettel/' for storing the zettel.
func (i Initiator) Init() error {
	var initPath = i.path + "/zettelkasten"
	if exists(initPath) {
		return fmt.Errorf("target path already exists: %v", initPath)
	}

	err := makeDir(i.path + "/zettelkasten")
	if err != nil {
		return err
	}
	err2 := os.WriteFile(i.path+"/zettelkasten/references.bib", []byte{}, fs.ModePerm)
	if err2 != nil {
		return err2
	}
	err3 := os.WriteFile(i.path+"/zettelkasten/index.txt", []byte{}, fs.ModePerm)
	if err3 != nil {
		return err3
	}
	err4 := os.Mkdir(i.path+"/zettelkasten/zettel/", fs.ModePerm)
	if err4 != nil {
		return err4
	}
	return nil
}

// InitExample creates a zettelkasten example for learning and demonstration purposes.
func (i Initiator) InitExample() error {
	var initPath = i.path + "/tmp_download"
	if exists(initPath) {
		return fmt.Errorf("target path already exists: %v", initPath)
	}

	resp, err := http.Get("https://github.com/crelder/zet_example/archive/refs/heads/main.zip") // Downloads a zip file
	if err != nil {
		return fmt.Errorf("could not download example: %v. Are you connected to the internet?", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("could not download example, code: %v", resp.StatusCode)
	}

	// Create target dir
	err2 := os.MkdirAll(initPath, fs.ModePerm)
	if err2 != nil {
		return err2
	}
	// Create the file
	out, err3 := os.Create(initPath + "/example.zip")
	if err3 != nil {
		return err3
	}
	_, err4 := io.Copy(out, resp.Body)
	if err4 != nil {
		return err4
	}

	err5 := filterExtract(initPath+"/example.zip", i.path)
	if err5 != nil {
		return err5
	}

	os.Rename(i.path+"/zet_example-main", i.path+"/zettelkasten")

	// We successfully extracted the zip file and therefore can delete it now, including the tmp folder.
	err6 := os.Remove(initPath + "/example.zip")
	if err6 != nil {
		return err6
	}
	err7 := os.Remove(initPath)
	if err7 != nil {
		return err7
	}

	return nil
}

// filterExtract takes as parameter a zip source filename and a destination path and extracts the ZIP archive.
func filterExtract(zipFilename, destPath string) error {
	// Open the source filename for reading
	zipReader, err := zip.OpenReader(zipFilename)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// For each file in the archive
	for _, archiveReader := range zipReader.File {

		// Open the file in the archive
		archiveFile, err := archiveReader.Open()
		if err != nil {
			return err
		}
		defer archiveFile.Close()

		// Prepare to write the file
		finalPath := filepath.Join(destPath, archiveReader.Name)

		// Check if the file to extract is just a directory
		if archiveReader.FileInfo().IsDir() {
			err = os.MkdirAll(finalPath, 0755)
			if err != nil {
				return err
			}
			// Continue to the next file in the archive
			continue
		}

		const maxLength = 150 // Max zip file filename length
		if len(archiveReader.Name) >= maxLength {
			return errors.New("Too long filename: " + archiveReader.Name)
		}

		// Create all needed directories
		if os.MkdirAll(filepath.Dir(finalPath), 0755) != nil {
			return err
		}

		// Prepare to write the destination file
		destinationFile, err := os.OpenFile(finalPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, archiveReader.Mode())
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		// Write the destination file
		if _, err = io.Copy(destinationFile, archiveFile); err != nil {
			return err
		}
	}

	return nil
}
func makeDir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return nil
}

// exists returns false only if the directory doesn't exist.
// In all other cases it returns a save true.
func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}
		return true
	}
	return true
}
