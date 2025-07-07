package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

type Task struct {
	Id        int        `json:"id"`
	Text      string     `json:"text"`
	Completed bool       `json:"completed"`
	Timestamp *time.Time `json:"timestamp"`
}

func (t *Task) EncodeJson() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Task) ToString() string {
	c := " "
	timestamp := " "
	if t.Completed {
		c = "x"
		timestamp = t.Timestamp.Format(time.DateTime)
	}

	return fmt.Sprintf("%d. [%s] %s %s", t.Id, c, t.Text, timestamp)
}

func Init(dbPath string) error {
	var err error

	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func CreateTask(task string) (int, error) {
	var id int

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		key := itob(int(id64))
		id = int(id64)

		newTask := Task{
			Id:   id,
			Text: task,
		}

		json, err := newTask.EncodeJson()
		if err != nil {
			return err
		}

		return b.Put(key, json)
	})

	if err != nil {
		return -1, err
	}
	return id, nil
}

func AllTasks() ([]Task, error) {
	var tasks []Task

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task Task
			err := json.Unmarshal(v, &task)
			if err != nil {
				panic(err)
			}

			tasks = append(tasks, task)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

func DoTask(key int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		byteTaskId := itob(int(key))
		taskRes := b.Get(byteTaskId)
		if len(taskRes) < 1 {
			fmt.Printf("Task with id %d not found\n", key)
			return nil
		}

		var task Task
		err := json.Unmarshal(taskRes, &task)
		if err != nil {
			return err
		}

		task.Completed = !task.Completed

		if task.Completed {
			now := time.Now()
			task.Timestamp = &now
		} else {
			task.Timestamp = nil
		}

		jsonTask, err := task.EncodeJson()
		if err != nil {
			return err
		}

		err = b.Put(byteTaskId, jsonTask)
		return err
	})

	return err
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
