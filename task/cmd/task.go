package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/boltdb/bolt"
)

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

const (
	BUCKET_NAME = "tasks"
)

type Context struct {
	Debug bool
	DB    *bolt.DB
}

type AddCmd struct {
	Task []string `arg:"" name:"task" help:"Task to add."`
}

func (add *AddCmd) Run(ctx *Context) error {
	db := ctx.DB
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_NAME))
		id, _ := b.NextSequence()
		nextId := int(id)

		err := b.Put(itob(nextId), []byte(strings.Join(add.Task, " ")))
		return err
	})
	return err
}

type ListCmd struct{}

func (ls *ListCmd) Run(ctx *Context) error {
	db := ctx.DB
	err := db.View(func(tx *bolt.Tx) error {
		fmt.Println("Tasks in bucket: " + BUCKET_NAME)
		b := tx.Bucket([]byte(BUCKET_NAME))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("%d - %s\n", binary.BigEndian.Uint64(k), v)
		}
		return nil
	})
	return err
}

type DoCmd struct {
	TaskId uint32 `arg:"" name:"task index" help:"Task to complete."`
}

func (do *DoCmd) Run(ctx *Context) error {
	fmt.Println("Completing task: ", do.TaskId)
	return nil
}

var Cli struct {
	Debug bool `help:"Enable debug mode."`

	Add  AddCmd  `cmd:"" help:"Add a todo."`
	Do   DoCmd   `cmd:"" help:"Complete a todo."`
	List ListCmd `cmd:"" help:"List all todos."`
}

func initDb() *bolt.DB {
	db, err := bolt.Open("tasks.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		if err != nil {
			log.Fatalf("create bucket: %s\n", err)
		}
		return nil
	})

	return db
}

func main() {
	db := initDb()
	defer db.Close()

	fmt.Println("task is a CLI to manage your TODOs")
	ctx := kong.Parse(&Cli)
	err := ctx.Run(&Context{Debug: Cli.Debug, DB: db})

	ctx.FatalIfErrorf(err)
}
