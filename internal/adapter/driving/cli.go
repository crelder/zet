// Package cli implements an adapter between the cli and a port of the application logic (al).
package cli

import (
	"fmt"
	"github.com/crelder/zet/internal/core/port"
	"os"
)

const version = "0.0.1"

type App struct {
	importer  port.Importer
	viewer    port.Viewer
	validator port.Validator
}

func NewApp(importer port.Importer, viewer port.Viewer, validator port.Validator) App {
	return App{
		importer:  importer,
		viewer:    viewer,
		validator: validator,
	}
}

func (cli App) Start() error {
	checks()

	subcmd := os.Args[1]
	switch subcmd {
	case "import":
		if len(os.Args) != 3 {
			return fmt.Errorf("no path provided. Please provide a path to the folder, where the textfiles lie, which you want to import")
		}
		n, errs := cli.importer.CreateImports(os.Args[2])
		if errs != nil {
			for _, err := range errs {
				return fmt.Errorf("Error importing: %v\n", err)
			}
		}
		fmt.Printf("Imported %d zettel into '/IMPORT'", n)
		return nil
	case "views":
		if len(os.Args) > 2 {
			return fmt.Errorf("command 'zet view' does not need any parameters")
		}
		err := cli.viewer.CreateViews()
		if err != nil {
			return fmt.Errorf(
				"Could not create views: %v\n"+
					"Make sure your zettelkasten is consistent by checking with `zet validate`", err)
		}
		return nil
	case "validate":
		if len(os.Args) > 2 {
			return fmt.Errorf("command 'zet validate' does not need any parameters")
		}
		valErr := cli.validator.Val()
		if valErr != nil {
			fmt.Printf("There are some inconsistencies in your zettelkasten:\n")
			fmt.Printf("Error1 = Duplicate Ids, Error2 = Dead Links:\n")
			fmt.Printf("\n")
			for _, validatorErr := range valErr {
				fmt.Printf("Error: %v, for id: %v\n", validatorErr.ErrType, validatorErr.Id)
			}
			return nil
		}
		fmt.Printf("Your zettelkasten seems to be okay.")
		return nil
	case "version":
		fmt.Printf("zet version %v", version)
		return nil
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		printUsage()
		return nil
	}
}

func checks() {
	if !isCalledFromZetDir() {
		fmt.Println("It seems you are not in your zettelkasten directory.")
		fmt.Println("Navigate to your zettelkasten directory to execute zet commands.")
		os.Exit(0)
	}

	if len(os.Args) <= 1 {
		printUsage()
		os.Exit(0)
	}
}

const usage = `Usage: zet <command> [<args>]
		
These are common zet commands:
   import <path>	Copy text files in <path> to folder 'IMPORT'
   validate     	Check your zettelkasten's consistency
   views     	 	Generate folder 'VIEWS', which contains access points into you zettelkasten

All Zet commands operate read-only on the three elements of the zettelkasten:
  * folder 'zettel'  (contains all zettel as a .txt, .png or .pdf file)
  * index.txt        (contains manually created starting points into your zettelkasten)
  * literature.bib   (contains information on sources - needed especially for scientific writing)`

func printUsage() {
	fmt.Printf(usage)
}

const (
	ZettelFolder   = "zettel"
	LiteratureFile = "literature.bib"
	IndexFile      = "index.txt"
)

// isCalledFromZetDir checks if a folder "zettel" and two files literature.bib and index.txt exist
// If one of these are non-existent it assumes that the user is not in his zettelkasten directory.
func isCalledFromZetDir() bool {
	_, err := os.Stat(ZettelFolder)
	_, err2 := os.Stat(LiteratureFile)
	_, err3 := os.Stat(IndexFile)
	if os.IsNotExist(err) || os.IsNotExist(err2) || os.IsNotExist(err3) {
		return false
	}
	return true
}
