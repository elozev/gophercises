package main

import (
	"fmt"
	"task/pkg/db"
	"time"
)

type ListCmd struct {
	ShowAll        bool   `name:"show-all" help:"Show all tasks (inc. completed)"`
	CompletedSince string `name:"completed-since" help:"Show all tasks since given duration. (1s, 5m, 3h, etc.)"`
}

func (ls *ListCmd) Run(ctx *Context) error {
	showAll := ls.ShowAll
	completedSince := ls.CompletedSince

	showCompletedSince := completedSince != ""

	var earlier time.Time

	if showCompletedSince {
		now := time.Now()

		sinceDuration, err := time.ParseDuration(completedSince)
		if err != nil {
			fmt.Printf("Failed to parse %s as a duration! Use proper formatting such as 1s, 3m, 4h, etc.", completedSince)
			panic(err)
		}

		earlier = now.Add(-1 * sinceDuration)

		fmt.Printf("Showing completed tasks since %s\n", earlier.Format(time.DateTime))
	}

	tasks, err := db.AllTasks()
	if err != nil {
		fmt.Println("Failed to get all tasks! Exiting...")
		return err
	}

	for _, task := range tasks {
		if showCompletedSince && task.Timestamp != nil && task.Timestamp.After(earlier) {
			fmt.Println(task.ToString())
		} else if !showCompletedSince && (showAll || !task.Completed) {
			fmt.Println(task.ToString())
		}
	}

	return err
}
