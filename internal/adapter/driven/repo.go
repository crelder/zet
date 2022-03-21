package repo

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/crelder/zet/internal/core/bl"
)

// Repo allows access to the content of your zettelkasten.
// path represents the path to the directory, where your zettelkasten lies.
type Repo struct {
	path string
}

type zettelFile struct {
	zettel   bl.Zettel
	filename string
	path     string
}

func NewRepo(path string) Repo {
	return Repo{
		path: path,
	}
}

// GetZettel returns all zettel of your zettelkasten.
func (r Repo) GetZettel() ([]bl.Zettel, error) {
	var zettels []bl.Zettel
	files, err := r.getFiles()
	if err != nil {
		return nil, err
	}
	for _, zf := range files {
		zettels = append(zettels, zf.zettel)
	}
	return zettels, nil
}

func (r Repo) getFiles() ([]zettelFile, error) {
	var zettelPath = r.path + "/zettel"
	dirEntries, err := os.ReadDir(zettelPath)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while reading the zettel from: %v", err)
	}

	var zettelFiles []zettelFile
	var z bl.Zettel
	for _, file := range dirEntries {
		z, _ = bl.ParseFilename(file.Name())
		zettelFiles = append(zettelFiles, zettelFile{
			zettel:   z,
			filename: file.Name(),
			path:     zettelPath,
		})
	}
	return zettelFiles, nil
}

// GetIndex returns the index of your zettelkasten.
// One index entry is a keyword and several ids which are
// access points into a line of thought regarding this keyword.
func (r Repo) GetIndex() ([]bl.Index, error) {
	var index []bl.Index
	f, err := os.ReadFile(r.path + "/index.txt")
	if err != nil {
		return []bl.Index{}, fmt.Errorf("an error occurred while reading the index: %v", err)
	}
	index = bl.ParseIndex(string(f))

	return index, nil
}

// CreateSyml creates links into your zettelkasten.
// Prefix specifies the folder where the links get places.
func (r Repo) CreateSyml(prefix string, m map[string][]string) error {
	for key, ids := range m {
		for _, id := range ids {
			fp, err := r.getFilePath(id)
			_, filename := filepath.Split(fp)
			if err != nil {
				return err
			}

			oldname := fp
			newname := r.path + "/VIEWS/" + prefix + "/" + key + "/" + filename

			persist(oldname, newname)
		}
	}
	return nil
}

func persist(oldname, newname string) {
	dir, _ := filepath.Split(newname)
	existsOrMake(dir)

	err := os.Symlink(oldname, newname)
	if err != nil {
		fmt.Printf("Could not create symlink: %q\n%v\n\n", newname, err)
	}
}

func (r Repo) getFilePath(id string) (string, error) {
	zfs, err := r.getFiles()
	if err != nil {
		return "", err
	}
	for _, zf := range zfs {
		if id == zf.zettel.Id {
			result := zf.path + "/" + zf.filename
			return result, nil
		}

	}
	return "", fmt.Errorf("id not found: %v", id)
}

func existsOrMake(dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		// TODO: Should return an error, which then gets handled!
		fmt.Printf("Got error: %v\n", err)
	}
}

func (r Repo) FileExists(link string) bool {
	_, err := os.Stat(r.path + "/" + link)
	return err == nil
}

// CreateFolgezettelStruct creates a folgezettel line within the repo.
func (r Repo) CreateFolgezettelStruct(topic string, links []bl.Symlink) error {
	for _, link := range links {
		fp, err := r.getFilePath(link.TargetId)
		if err != nil {
			return err
		}

		oldname := fp
		newname := r.path + "/VIEWS/index/" + topic + "/" + links[0].TargetId + "/" + link.LinkName
		dir, _ := filepath.Split(newname)
		existsOrMake(dir)

		err = os.Symlink(oldname, newname)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r Repo) GetContents(path string) []string {
	var contents []string

	filepaths, err := os.ReadDir(path)
	filepaths = filterAllowed(filepaths)
	if err != nil {
		fmt.Printf("Could not read files from directory: %v", path)
		os.Exit(0)
	}

	for _, fp := range filepaths {
		dat, _ := os.ReadFile(path + "/" + fp.Name())
		contents = append(contents, string(dat))
	}
	return contents
}

func exists(path string) bool {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return false
	}
	return true
}

func filterAllowed(filepaths []os.DirEntry) []os.DirEntry {
	var fps []os.DirEntry
	for _, fp := range filepaths {
		if isAllowed(fp.Name()) {
			fps = append(fps, fp)
		}
	}
	return fps
}

func isAllowed(fn string) bool {
	if fn[len(fn)-3:] == "txt" {
		return true
	}
	return false
}

// checkIsEmpty returns true, if the passed folder is either non-existent or empty.
// Otherwise, it returns false.
func checkIsEmpty(impPath string) bool {
	f, err := os.Open(impPath)
	if err != nil {
		return true
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	return err == io.EOF       // Either not empty or error, suits both cases
}

// SaveImports creates text files with a valid filename and the content.
// The parameter expects map[filename]content.
func (r Repo) SaveImports(zfs map[string]string) (int, error) {
	// Check if import folder is empty
	impPath := r.path + "/IMPORT"
	if checkIsEmpty(impPath) == false {
		return 0, fmt.Errorf("the folder where the files with the parsed filenmames get placed is not empty. Please delete within your zettelkasten the 'import' folder or empty it")
	}

	// If path '/IMPORT does not exist, create it
	if exists(impPath) == false {
		err := os.Mkdir(impPath, 0755)
		if err != nil {
			return 0, fmt.Errorf("could not create import folder at %q", impPath)
		}
	}

	// Write files to Import path
	counter := 0
	for filename, content := range zfs {
		filename := impPath + "/" + filename
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			return counter, fmt.Errorf("could not write file %v", filename)
		}
		counter++
	}
	return len(zfs), nil
}
