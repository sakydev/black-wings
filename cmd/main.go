package main

import (
	"black-wings/cmd/account"
	"black-wings/cmd/search"
	db "black-wings/database"
	"black-wings/internal"
	"black-wings/internal/repositories"
	"black-wings/internal/types"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

const AppName = "black-wings"

var settings map[string]string

func main() {
	ctx := context.Background()
	injector := do.DefaultInjector
	database := db.GetDatabase()
	setup(ctx, database, injector)

	defer database.Close()

	var format string

	coreCommand := &cobra.Command{
		Use:   AppName,
		Short: "Search everywhere",
		Long:  `A command line utility to search everywhere.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Use '%s --help' for more information.\n", AppName)
		},
	}

	coreCommand.PersistentFlags().StringVarP(&format, "format", "f", "json", "Data format to use (default: json)")

	commandConfiguration := types.CommandConfiguration{
		Injector: injector,
		Database: database,
		Context:  ctx,
		Settings: settings,
		Format:   format,
	}
	searchCommand := search.NewSearchCommand(commandConfiguration)
	accountCommands := account.NewAccountCommands(commandConfiguration)

	coreCommand.AddCommand(searchCommand)
	coreCommand.AddCommand(accountCommands)

	err := coreCommand.Execute()
	if err != nil {
		fmt.Println(color.RedString("Error: %v\n", err))
		os.Exit(1)
	}
}

func setup(ctx context.Context, database *sql.DB, injector *do.Injector) {
	var err error
	settings, err = repositories.SettingImpl{}.List(ctx, database)
	if err != nil {
		fmt.Println(color.RedString("error getting settings: %v\n", err))
		os.Exit(1)
	}

	internal.WireDependencies(injector)
}
