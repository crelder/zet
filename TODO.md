## TODOs and Refactorings

1. `repo_test.go`: Bench test shows suboptimal process time for persisting. Introduce go routines to make this
   persisting in parallel.
2. `driven.go`:  refactor and remove `FileExists(link string) bool` from Interface, since it is only used for testing
   purposes
3. `CreateFolgezettelStruct(prefix string, links []bl.Symlink) error`: replace []bl.Symlink via map, so that it
   communicates via simple data structures across package borders (= fewer dependencies between packages)
4. `GetIndex() ([]bl.Index, error)`: Same as above: replace []bl.Index via a map
5. `type Zettel struct`: refactor and make all properties private
6. `type Zettel struct`: remove property `name`, which is the filename, from struct (the model shouldn't care about what
   it is called when persisted)
7. Completely rethink error treatment ("return or handle") from source to end in the application, especially in the
   repo, parsing Methods and the validator. E.g. `ParseFilename` should throw errors. E.g. add type repoErr (with
   filename and a errMsg), which collects parse errors.
8. Investigate bug for file: /Users/crelder/zettelkasten/VIEWS/keywords/Prozess/180917v - Werden, Prozess, Prozess,
   Stabilisierung - Daniel.png  
   `symlink /Users/c/zettelkasten/zettel/180917v - Werden, Prozess, Prozess, Stabilisierung - Daniel.png /Users/c/zettelkasten/VIEWS/keywords/Prozess/180917v - Werden, Prozess, Prozess, Stabilisierung - Daniel.png`:
   file exists -> Why does it already exists and throws a bug?  
   If the import or views command has at least one repo error, than tell the user that he should run the command zet
   validate.
9. `func parseContextFromFilename`: doesn't work when the filename has no keywords. Make this function more robust.
10. `index_test.go`: Add more testcases, like "no :", "two :", etc.
11. Write missing test cases for public methods in `zk.go`

## New features

### Import.go

1. import command should import the zettel directly into your zettelkasten/zettel folder. Then remove
   global `var zk bl.Zettelkasten`
1. TestCreateImports: implement a test case, which checks, if an error is returned, if IMPORT folder already exists

### Validate.go

1. validate if literature exists
1. validate if folgezettel form a tree structure (not a network)
1. validate if the filename has the correct format:
   `[ID] - [keywords, optional] - [citations OR context, optional] - [folgezettel, optional].[file type]`
1. validate if the header in the text files is correct (e.g. `[keywords]\n[date]\n[optional:context and literature]\n)
1. validate if every zettel has 0 or 1 predecessor

### View.go

1.Method `addSymlink`: introduce variable traveledIds []string, which checks, that we do not have traverse loops, which
would lead in an endless loop

1. After having implemented a sync Viewer.createViews() Method and noticing that it takes more than one minute, I
   measured via a Benchmark test. I then implemented the concurrency solution in the repo. But this could be optimized
   further by playing around with the channel capacity, etc. if needed.

1. If my concurrency calls were interlinked more deeply, I would use a context to track what is happening (since it can
   be hard in concurrency to understand what happened in the retrospective). But here it is just the parallel execution
   of one method basically and just collecting the errors returned. No interlink between the go routines.

### Repo

1. The getting the zettelkasten from the repo (parsing all filenames) is really fast. If I would notice a slow behavior
   I would a) measure it and b) use a cache in the repo.

### General

1. use brew or apt-get to install this app
1. add a separate folder with own works (e.g. papers, books, thesis, presentations, etc.)?
1. Use registry for magic strings (like "pathToZettel", etc.), something like:

```go
package main

import "sync"

type registry struct {
	m    sync.RWMutex
	data map[string]string
}

func (r *registry) setVal(k, v string) {
	r.m.Lock()
	defer r.m.Unlock()

	r.data[k] = v
}

func (r *registry) getVal(k string) (string, bool) {
	r.m.RLock()
	defer r.m.RUnlock()

	v, ok := r.data[k]
	return v, ok
}
```
