package main

import (
	"fmt"
	"github.com/crelder/zet/internal/adapter/driven"
	"github.com/crelder/zet/internal/adapter/driving"
	"github.com/crelder/zet/internal/core/as"
	"log"
	"os"
)

func main() {
	if r := run(); r != nil {
		log.Fatal(r)
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
	// Since the CLI tool can only be called within a zettelkasten directory,
	// the current working directory is also the path to the zettelkasten directory.
	wd, err := os.Getwd()
	if err != nil {
		return cli.App{}, fmt.Errorf("could not read the current working directory")
	}

	// Wire app together
	r := repo.NewRepo(wd)
	viewer := as.NewViewer(r, r)
	importer := as.NewImporter(r, r)
	validator := as.NewValidator(r)

	return cli.NewApp(importer, viewer, validator), nil
}
