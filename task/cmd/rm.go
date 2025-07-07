package main

import (
	"fmt"
	"task/pkg/db"
)

type RmCmd struct {
	TaskId int `arg:"" name:"task index" help:"Task to remove."`
}

func (rm *RmCmd) Run(ctx *Context) error {
	taskId := rm.TaskId

	fmt.Printf("Removing task n.%d from the list\n", taskId)

	err := db.DeleteTask(taskId)

	return err
}
