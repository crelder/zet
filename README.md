# zet

Zet is a command line tool for a personal knowledge management system using the zettelkasten method.

It was built to be simple, robust, and designed to last a lifetime.
## zet solves three problems

1. You want to find a thought that you know to exist in your knowledge management system
2. You want to find relevant thoughts around a particular topic that you don't know to exist in your knowledge management
   system
3. You want a tool that supports in developing lines of thoughts and train your ability to think in a structured manner

## How does `zet` work?

A folder (e.g. `zettelkasten`) holds everything you need for your personal knowledge management. It has three elements:

* a folder `zettel/` with zettel (= images or text files with the content of your thoughts)
* a `literature.bib` with information on literature sources
* an `index.txt` file (contains entry points into lines of thought about a particular topic)

### folder `zettel/`

Every zettel filename contains all metainformation.  
E.g.: `170211a - Algorithm, Efficiency - GopherCon, sedgewick2011 - 180204d.txt`

The filename consists of four parts, which are always in the same order:

1. The first part, the id `170211a` (= date).
2. The second part, `Algorithm, Efficiency`, contains keywords.
3. The third part, `GopherCon, sedgewick2011`, contains context `GopherCon` and literature sources `sedgewick2011`.
4. The fourth part, `180204d`, contains one or more "folgezettel". These are links to other zettel. These links form a
   tree (not a network).

Only the first part is mandatory. The following three parts are optional.

### index.txt

`index.txt` contains thematic entry points into your zettelkasten, for example:

```text
Complexity: 190119e
Interfaces: 220115p
```

### literature.bib

`literature.bib` contains the list of literature references. This file is not only beneficial for scientific writing.

### Import and Views functionality

You then navigate into your zettelkasten folder in the shell and can run the following zet commands:

```shell
Usage: zet <command> [<args>]

These are common zet commands:
  import <path> Copy text files in <path> to folder 'IMPORT/'
  validate      Check your zettelkasten's consistency
  views         Generate folder 'VIEWS/', which contains access points into your zettelkasten

All Zet commands operate read-only on the three elements of the zettelkasten:
  * folder 'zettel'  (contains all zettel as a .txt, .png or .pdf file)
  * index.txt        (contains manually created starting points into your zettelkasten)
  * literature.bib   (contains information on sources - needed especially for scientific writing)
```

`zet validate`: tells you if you have inconsistencies e.g., dead links, duplicated ids

`zet import <path-to-folder>`: imports all .txt files from a folder into the zettelkasten and parses the header in the
text file into a correct filename with meta information.

`zet views`: generates several access points into your zettelkasten.

To see how the `zet import` command works navigate to `internal/core/as/testdata/import/zettelkasten`. In order to
import the textfiles in `internal/core/as/testdata/import/new_zettel_files` run all tests (run `go test ./...`). You
will then see a folder "IMPORTS", with the newly imported zettel with a correct file name.

The header of the new zettel is used to create a filename. When a zettel is in the form of a text, it always contains
two to three lines with metadata. These first two to three lines form the "header" of a textfile.

To see how the `zet views` command works within your zet codebase navigate
to `internal/core/as/testdata/views/zettelkasten` and run all tests by executing `go test ./...`. You will then see a
folder `VIEWS/` with all the access points into the small example zettelkasten. You have here the following access points
into your zettelkasten via the following folders:

* `keywords/`: for every keyword you see all relevant zettel
* `context/`: for every context you see all relevant zettel
* `citations/`: for every citation (= literature source), you see all relevant zettel
* `explore/`: for every (zettel) id you see all other related zettel
* `index/`: for every index entry you see all relevant zettel

## Installation

Run `go install github.com/crelder/zet/cmd/cli@latest` (As soon as this project is a public github project...)

## Problems that the tool needs to tackle

These are limitations some existing personal knowledge management solutions have. This tool needs to solve these.

1. Lock-in of information
2. Long usage until it gets useful
3. f(shit)=shit, therefore protection against overflowing with useless or bad input
4. Your knowledge shouldn't be represented as a network because humans have difficulty reading networks. They are used
   to reading linear structures (e.g., books have a linear structure, Wikipedia doesn't)

## Code Structure

The project is structured following hexagonal architecture.

The application's entry point is in `cmd/cli/main.go`, where the app is wired together. The cli adapter will return
the app. It uses the driving ports of your core, which are located in `internal/core/port`, alongside your business
logic and application services.

Package `bl` contains all the business logic and models of this program. The here exposed functions help to fulfill the
three functions in the application service layer (package `as`):

* Creating all kinds of views on the zettelkasten
* Validating the zettelkasten (= checking for inconsistencies)
* Supporting the import of zettel in text format by automatically creating the filename based on the zettel header

The application layer uses the driven ports to access infrastructure functionality, the repo. This tool doesn't
deliberately use a database. You should be able to access your notes throughout your
life easily. Using just a folder is simpler and more robust (therefore, no information lock-in).

The two assumptions are:

a) that you will always have a filesystem with a folder, where you can store text and image files.

b) that you will always have a search functionality for searching through your filenames.

## Combining with other tools

This cli tool has its full potential when combined with other tools.

* Syncing across all your devices: e.g., icloud, nextcloud
* Version control: e.g. git
* Scanning paper zettel: e.g., Adobe Scan App
* Dictate thoughts: e.g., IA Writer with ios build-in speech recognition

These tools will most likely change over your lifetime. The invariant of this system are the three zettelkasten
elements (text and images with your thought, text files with an index, and literature).

## Further reading

I recommend this presentation on the questions how Niklas Luhmann worked with his
zettelkasten ([link](https://strengejacke.files.wordpress.com/2015/10/introduction-into-luhmanns-zettelkasten-thinking.pdf))

See current bugs, issues and planned feature [here](./TODO.md).

I build an earlier version; see [here](https://github.com/crelder/zettelkasten). I kept the functionality that I found
useful (the structure of the filenames, an importer functionality), but left the stuff that I didn't find helpful behind (a
gephi graph, a UI for seeing the folgezettel structure).

## Contributing

I develop this project in my spare time. Therefore, I am happy about any contributions to this project, be it new
features, bug fixes, fixing typos, improving documentation, suggestions on how to design the code better, you name it!

Just open a new issue ticket or a pull request. I promise to respond, but it might take a while: either because I am
busy with other stuff or because I am thinking about what you wrote.

## Acknowledgments
