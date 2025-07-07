package main

import (
	"fmt"
	"task/pkg/db"
)

type DoCmd struct {
	TaskId int `arg:"" name:"task index" help:"Task to complete."`
}

func (do *DoCmd) Run(ctx *Context) error {
	fmt.Println("Completing task: ", do.TaskId)
	err := db.DoTask(do.TaskId)
	return err
}
