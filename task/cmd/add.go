package main

import (
	"strings"
	"task/pkg/db"
)

type AddCmd struct {
	Task []string `arg:"" name:"task" help:"Task to add."`
}

func (add *AddCmd) Run(ctx *Context) error {
	_, err := db.CreateTask(strings.Join(add.Task, " "))
	return err
}
