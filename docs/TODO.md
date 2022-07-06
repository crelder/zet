# TODOs 

## Bug fixes and Refactorings

1. Not skipping invisible files. `inconsistency: parse Filename: could not parse id from filename ".DS_Store"`
2. not detected that "B" is uppercase in bibkey: `170105a - Entrepreneurshipforschung, Kleinunternehmen, neuartige Unternehmen, Organisationslösungen - Bögenhold2016 - 170722.png`. Add testcase and fix.
3. Folder `VIEWS/unlinked` shows sometimes less than 10 zettel. It should be always 10 (if there are 10).
4. Folder `VIEWS/unlinked` should not contain ids that already exist int the index.
5. `func parseContextFromFilename` doesn't work when the filename has no keywords. Make this function more robust.
6. Add test case that in `VIEWS/unlinked` are only these files, that have no predecessor link
7. Add test case that checks that circular dependencies are detected in `VIEWS/index`.
8. Add missing testcase for `VIEWS/unlinked`: shouldn't list `220122a - 191212b.txt`
9. Write more testcases in `TestParseFilename` like failing cases, without ID, without keywords, etc. and then check the error.
10. `func getDeadLinks(zk zet.Zettelkasten) []string {`: use map to remove second for loop.
11. Changing this line doesn't return an error! `incons = append(incons, fmt.Errorf("inconsistency: %v", parseErr))`. But it should! Check testcases.

## New features for `zet validate`

1. Validate if the folgezettel structure forms a tree structure (not a network)
2. validate if the filename has the correct format:
   `[ID] - [keywords, optional] - [references OR context, optional] - [folgezettel, optional].[file type]`
3. validate if every zettel has [0,1] predecessor
4. Validator also gives the info: 19223 zettel, 120 indexes, 40 bibkeys
5. After importing, ask with prompt: do you want to create new views? If yes, run `zet views`
6. After creating views run validate ("there are inconsistencies. Run zet validate.")
7. Check similarity of an imported zettel. Does one already similar exist? If date is the same (except letter) and text
    content high similarity, error.

## Other

1. use brew or apt-get to install this app
1. Use registry for magic strings (like "pathToZettel", etc.)
2. Issue tracker in Github einführen. Survey: Why did you come here? Which feature do you seek? Validate, import, views, Gephi, Ui
   with folgezettel struct, Meta data in filename, other? Missing features? Downsites? New ideas for this program?
3. Add CI Build scripts in github, see Chapter 13:continous integration in "learning Test-Driven Development"