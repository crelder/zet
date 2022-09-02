package main

import (
	"fmt"
	"github.com/crelder/zet/pkg/imports"
	"github.com/crelder/zet/pkg/initialize"
	"github.com/crelder/zet/pkg/parse"
	"github.com/crelder/zet/pkg/transport/cli"
	"github.com/crelder/zet/pkg/transport/fs"
	"github.com/crelder/zet/pkg/validate"
	"github.com/crelder/zet/pkg/view"
	"log"
	"os"
)

func main() {
	if r := run(); r != nil {
		log.Print(r)
	}
}

func run() error {
	app, err := createApp()
	if err != nil {
		return err
	}

	if r := app.Start(); r != nil {
		return r
	}
	return nil
}

func createApp() (cli.App, error) {
	// The cli application can only run within a zettelkasten directory.
	// Therefore, the current working directory is also the path to the zettelkasten directory.
	wd, err := os.Getwd()
	if err != nil {
		return cli.App{}, fmt.Errorf("could not read the current working directory: %v", err)
	}

	// Wire app together
	parser := parse.New()
	repo := fs.New(wd, parser)
	viewer := view.New(repo, repo)
	importer := imports.New(parser, repo, repo)
	validator := validate.New(repo)
	initiator := initialize.New(wd)

	return cli.NewApp(importer, viewer, validator, initiator), nil
}
