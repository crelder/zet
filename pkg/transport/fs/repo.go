package fs

import (
	"errors"
	"fmt"
	"github.com/crelder/zet"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Repo allows access to the content of your zettelkasten.
// Repo satisfies the zet.Repo, view.Persister and imports.Reader interface.
// path represents the path to the directory, where your zettelkasten lies.
type Repo struct {
	parser zet.Parser
	path   string
}

type zettelFile struct {
	zettel   zet.Zettel
	filename string
	path     string
}

func New(path string, p zet.Parser) Repo {
	return Repo{
		parser: p,
		path:   path,
	}
}

// GetZettel returns all zettel of your zettelkasten.
// zet.Zettel are ordered by id (descending).
//
// The second parameter []string contains all inconsistencies that occureed during parsing the zettel filename.
// Only if not an error happened, it will return zettel and inconsistencies.
func (r Repo) GetZettel() ([]zet.Zettel, error) {
	var zettels []zet.Zettel
	files, err := r.getFiles()
	if err != nil {
		return nil, err
	}
	for _, zf := range files {
		zettels = append(zettels, zf.zettel)
	}

	// Sort zettel by Id to make sure, that following results e.g. building the index,
	// always returning the same results.
	sort.Slice(zettels, func(i, j int) bool {
		return zettels[i].Id < zettels[j].Id

	})

	zettels, err = addFolgezettel(zettels)
	if err != nil {
		return nil, err
	}

	return zettels, nil
}

// addFolgezettel calculates the Folgezettel for each Zettel.
// This is needed, because in the filename, only the predecessor zettel are provided.
// The predecessors are used to calculate the Folgezettel for each Zettel.
// In case of a detected double id it will return an error.
func addFolgezettel(zettels []zet.Zettel) ([]zet.Zettel, error) {
	if len(zettels) == 0 {
		return nil, nil
	}

	var result []zet.Zettel
	folgezettelIds := getFolgezettelIds(zettels)
	for _, zettel := range zettels {
		zettel.Folgezettel = folgezettelIds[zettel.Id]
		result = append(result, zettel)
	}

	// Sort zettel by Id and sort the automatically added Folgezettel ids to make sure,
	// that following operations e.g. building the index always returns a reproducible result.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id

	})

	for _, z := range result {
		sort.Slice(z.Folgezettel, func(i, j int) bool {
			return z.Folgezettel[i] < z.Folgezettel[j]
		})
	}

	return result, nil
}

// getFolgezettelIds returns a map that has the id of a zettel and all follow up zettel ids (Folgezettel).
func getFolgezettelIds(zettels []zet.Zettel) map[string][]string {
	zetMap := make(map[string][]string)
	for _, zettel := range zettels {
		for _, predecessor := range zettel.Predecessor { // Normally this should be just one, but just in case...
			zetMap[predecessor] = append(zetMap[predecessor], zettel.Id)
		}
	}
	return zetMap
}

// getFiles returns raw data about zettel.
// When errors occur while parsing the filenames, getFiles returns these specific errors as the second parameter.
// All other errors are returned as the last parameter.
// Invisible files will be skipped - this works only for unix systems, since invisibility is determined
// by a dot at the beginning of the name.
//
// In case of double ids, an error get returned.
func (r Repo) getFiles() ([]zettelFile, error) {
	// Read all the zettel
	var zettelPath = r.path + "/zettel"
	dirEntries, err := os.ReadDir(zettelPath)
	if err != nil {
		return nil, fmt.Errorf("fs: %v", err)
	}
	var zettelFiles []zettelFile
	var z zet.Zettel
	for _, file := range dirEntries {
		if notValid(file) {
			continue
		}
		var parseErr error
		z, parseErr = r.parser.Filename(file.Name())
		if parseErr != nil {
			return nil, parseErr
		}
		zettelFiles = append(zettelFiles, zettelFile{
			zettel:   z,
			filename: file.Name(),
			path:     zettelPath,
		})
	}

	return zettelFiles, nil
}

// notValid checks if the filename is valid. Only visible files are valid.
// Therefore, the check currently only works on unix systems.
func notValid(file os.DirEntry) bool {
	return strings.HasPrefix(file.Name(), ".")
}

// GetIndex returns the index of your zettelkasten.
// One index entry is a keyword with several ids in the form of:
//        Evolution: 170311a
// One index entry can also have several ids:
//        Technology: 220112d, 190314f
//
// An Index is used as access point into a line of thought (=zettel chain) regarding this keyword.
// ParsingErrors are returned with the second parameter []error.
// All other errors via the last parameter.
func (r Repo) GetIndex() (zet.Index, error) {
	var index map[string][]string
	f, err := os.ReadFile(r.path + "/index.txt")
	if err != nil {
		return nil, fmt.Errorf("fs: %v", err)
	}

	index, parseErr := r.parser.Index(string(f))
	if parseErr != nil {
		return nil, parseErr[0]
	}
	return index, nil

}

func (r Repo) GetBibkeys() ([]string, error) {
	f, err := os.ReadFile(r.path + "/references.bib")
	if err != nil {
		return nil, fmt.Errorf("fs: %v", err)
	}

	return r.parser.Reference(string(f)), nil
}

// CreateSyml creates links into your zettelkasten.
// Prefix specifies the folder where the links get places,
// it is type of the view (e.g. keywords, literature sources, context).
// m contains as keys the keyword and as []string one or more ids.
func (r Repo) CreateSyml(prefix string, m map[string][]string) error {
	for link, ids := range m {
		for _, id := range ids {
			fp, err := r.getFilePath(id)
			_, filename := filepath.Split(fp)
			if err != nil {
				return err
			}

			oldname := fp
			var newname string
			if prefix != "" {
				newname = r.path + "/VIEWS/" + prefix + "/" + link + "/" + filename
			}
			if prefix == "" {
				newname = r.path + "/VIEWS/" + link + "/" + filename
			}

			err2 := persist(oldname, newname)
			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}

// CreateInfo persists some statistics in form of a txt file about a topic like e.g. keywords, context or literature.
func (r Repo) CreateInfo(filename string, m map[string][]string) error {
	err := existsOrMake(r.path + "/VIEWS/stats")
	if err != nil {
		return err
	}
	var stats []byte
	for word, ids := range m {
		stats = append(stats, []byte(word+";"+strconv.Itoa(len(ids))+"\n")...)
	}
	err2 := os.WriteFile(r.path+"/VIEWS/stats/"+filename+".csv", stats, fs.ModePerm)
	if err2 != nil {
		return err2
	}
	return nil
}

func persist(oldname, newname string) error {
	dir, _ := filepath.Split(newname)
	err := existsOrMake(dir)
	if err != nil {
		return err
	}

	err2 := os.Link(oldname, newname)
	if err2 != nil {
		return fmt.Errorf("fs: could not create symlink: %v\n", err2)
	}

	return nil
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

func existsOrMake(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("fs: %v", err)
	}

	return nil
}

// CreateFolgezettelStruct creates a tree like link structure, so called "Folgezettel" in the fs.
// This represents the physical representation how Niklas Luhmann arranged his Zettel in his
// wooden zettelkasten boxes. This is used for creating chains of thoughts.
// Topic is e.g. "Evolution" and the map contains all links[linkname]targetId
func (r Repo) CreateFolgezettelStruct(links map[string]string) error {
	for linkName, targetId := range links {

		fp, err := r.getFilePath(targetId)
		if err != nil {
			return err
		}

		oldname := fp
		newname := path.Join(r.path, linkName)
		dir, _ := filepath.Split(newname)
		err = existsOrMake(dir)
		if err != nil {
			return err
		}

		err = os.Link(oldname, newname)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetContents reads the files in the path and extracts from allowed files the text content.
func (r Repo) GetContents(path string) ([]string, error) {
	var contents []string

	filepaths, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("repo GetContents: %v for path: %v", err, path)
	}
	filepaths = filterAllowed(filepaths)

	for _, fp := range filepaths {
		dat, _ := os.ReadFile(path + "/" + fp.Name())
		contents = append(contents, string(dat))
	}
	return contents, nil
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

// Save creates text files with a valid filename and the content.
// The parameter expects map[filename]content.
func (r Repo) Save(zfs map[string]string) (int, error) {
	impPath := r.path + "/zettel"

	if exists(impPath) == false {
		return 0, errors.New("target persisting folder 'zettel/' should exist, but doesn't")
	}

	// Write files to Import path
	counter := 0
	for filename, content := range zfs {
		filename := impPath + "/" + filename
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			return counter, fmt.Errorf("could not write file %q: %v", filename, err)
		}
		counter++
	}
	return len(zfs), nil
}
