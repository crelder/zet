package cli

import (
	"fmt"
	"github.com/crelder/zet"
	"github.com/crelder/zet/pkg/export"
	"github.com/crelder/zet/pkg/index"
	"os"
)

const version = "0.3.0"

type App struct {
	importer  zet.Importer
	indexer   index.Indexer
	exporter  export.Exporter
	validator zet.Validator
	initiator zet.Initiator
}

func NewApp(importer zet.Importer, exporter export.Exporter, indexer index.Indexer, validator zet.Validator, initiator zet.Initiator) App {
	return App{
		importer:  importer,
		indexer:   indexer,
		exporter:  exporter,
		validator: validator,
		initiator: initiator,
	}
}

// Start orchestrates the cli call and calls the respective services in the application.
func (cli App) Start() error {
	if len(os.Args) <= 1 {
		fmt.Printf("Within your zettelkasten folder:\n")
		printUsage()
		return nil
	}

	subcmd := os.Args[1]
	if subcmd == "help" {
		fmt.Printf("Within your zettelkasten folder:\n")
		printUsage()
		return nil
	}

	if subcmd == "version" {
		fmt.Printf("zet version %v", version)
		return nil
	}

	if subcmd == "init" {
		if len(os.Args) == 3 {
			if os.Args[2] == "example" {
				fmt.Printf("Downloading an example zettelkasten...\n\n")
				err := cli.initiator.InitExample()
				if err != nil {
					return err
				}
				fmt.Printf("Created an example zettelkasten.\n\nYou can start a tutorial by opening 'zettelkasten/zettel/220527z'.")
				return nil
			}
		}
		err := cli.initiator.Init()
		if err != nil {
			return err
		}
		fmt.Printf("Created an empty zettelkasten, see 'zettelkasten' directory")
		return nil
	}

	// From here on, the user must be in his zettelkasten to execute the following commands.
	if !isCalledFromZetDir() {
		fmt.Println("It seems you are not in your zettelkasten directory.")
		fmt.Println("Navigate to your zettelkasten directory to execute zet commands or see 'zet help'.")
		return nil
	}

	switch subcmd {
	case "import":
		if len(os.Args) < 3 {
			return fmt.Errorf("no path provided. Please provide a path to the folder, where the textfiles lie, which you want to import")
		}

		n, err2 := cli.importer.Import(os.Args[2])
		if err2 != nil {
			if n == 0 {
				return fmt.Errorf("error importing: %v", err2)
			}
			if n > 0 {
				return fmt.Errorf("only %v zettel got imported, because of error %q.\n\nPlease check manually which zettel are still missing", n, err2)
			}
		}
		fmt.Printf("Imported %d zettel into your zettel folder", n)
		return nil
	case "export":
		if len(os.Args) > 2 {
			return fmt.Errorf("command 'zet export' does not need any parameters")
		}
		err := cli.exporter.Export()
		if err != nil {
			return fmt.Errorf("Could not create views: %v\n", err)
		}
		return nil
	case "index":
		if len(os.Args) > 2 {
			return fmt.Errorf("command 'zet index' does not need any parameters")
		}
		err := cli.indexer.Create()
		if err != nil {
			return fmt.Errorf("Could not create index: %v\n", err)
		}
		return nil
	case "validate":
		if len(os.Args) > 2 {
			return fmt.Errorf("command 'zet validate' does not need any parameters")
		}
		inconsistencies, err := cli.validator.Val()
		if err != nil {
			fmt.Printf("An error ocurred while validating your zettelkasten: %v", err)
		}
		if len(inconsistencies) > 0 {
			for _, err := range inconsistencies {
				fmt.Printf("%v\n", err)
			}
			return nil
		}
		fmt.Printf(`Your zettelkasten seems to be okay:
All ids are unique.
All links point to an existing zettel.
All ids in the index point to an existing zettel.
All bibkeys have a corresponding reference.`)

		return nil
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		printUsage()
		return nil
	}
}

const usage = `Usage: zet <command> [<args>]
      
These are common zet commands:
   export		   Generate folder 'EXPORT', which contains files with aggregated data 
   index           Generate folder 'INDEX', which contains thematic access points into your zettelkasten
   init            Creates an empty zettelkasten
   init example    Downloads an example zettelkasten which is a tutorial
   validate        Check your zettelkasten's consistency

All Zet commands operate read-only on the three elements of the zettelkasten:
  * index.txt        (contains manually created starting points into your zettelkasten)
  * folder 'zettel'  (contains all zettel as a .txt, .png or .pdf file)
  * references.bib   (contains information on sources - needed especially for scientific writing)`

func printUsage() {
	fmt.Printf(usage)
}

const (
	ZettelFolder   = "zettel"
	ReferencesFile = "references.bib"
	IndexFile      = "index.txt"
)

// isCalledFromZetDir checks if a folder "zettel" and two files references.bib and index.txt exist
// If one of these are non-existent it assumes that the user is not in his zettelkasten directory.
func isCalledFromZetDir() bool {
	_, err := os.Stat(ZettelFolder)
	_, err2 := os.Stat(ReferencesFile)
	_, err3 := os.Stat(IndexFile)
	if os.IsNotExist(err) || os.IsNotExist(err2) || os.IsNotExist(err3) {
		return false
	}
	return true
}
