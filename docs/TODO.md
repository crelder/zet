# TODOs 

## Bug fixes and Refactorings

1. Better error handling for getZettel(id, zettels) in index.go
2. Folder `VIEWS/unlinked` shows sometimes less than 10 zettel. It should be always 10 (if there are 10).
3. Folder `VIEWS/unlinked` should not contain ids that already exist int the index.
4. Add test case that in `VIEWS/unlinked` are only these files, that have no predecessor link
5. Add test case that checks that circular dependencies are detected in `VIEWS/index`.
6. Add missing testcase for `VIEWS/unlinked`: shouldn't list `220122a - 191212b.txt`
7. Delete AddFolgezettel, because it is implicitly also tested via validate and views tests? Does order matter? It should be in the way it is in the filename.

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