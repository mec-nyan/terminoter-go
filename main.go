package main

import (
	"log"

	tea "charm.land/bubbletea/v2"

	"github.com/mec-nyan/terminoter-go/internal/args"
	"github.com/mec-nyan/terminoter-go/internal/model"
	"github.com/mec-nyan/terminoter-go/internal/setup"
)

func main() {
	// App setup:
	//
	// Get command line arguments, if any:
	opts := args.ParseArgs()

	// Initial setup (save location, configuration options, etc).
	err := setup.Setup(&opts)
	if err != nil {
		log.Fatalf("initialisation error: %v", err)
	}

	app := tea.NewProgram(model.InitialModel(opts))

	if _, err := app.Run(); err != nil {
		log.Fatalf("Oops! Something went terribly wrong...: %v", err)
	}
}
