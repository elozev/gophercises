package main

import (
	"fmt"
	"os"
	"path/filepath"
	"task/pkg/db"

	"github.com/alecthomas/kong"
	"github.com/mitchellh/go-homedir"
)

const (
	BUCKET_NAME = "tasks"
)

type Context struct {
	Debug bool
}

var Cli struct {
	Debug bool `help:"Enable debug mode."`

	Add  AddCmd  `cmd:"" help:"Add a todo."`
	Do   DoCmd   `cmd:"" help:"Complete a todo."`
	List ListCmd `cmd:"" help:"List all todos."`
	Rm   RmCmd   `cmd:"" help:"Remove a todo."`
}

func main() {
	home, _ := homedir.Dir()

	dbPath := filepath.Join(home, "tasks.db")
	must(db.Init(dbPath))

	fmt.Println("task is a cli to manage your todos")
	ctx := kong.Parse(&Cli)
	err := ctx.Run(&Context{Debug: Cli.Debug})

	ctx.FatalIfErrorf(err)
}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
