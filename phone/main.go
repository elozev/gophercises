package main

import (
	"fmt"
	phonedb "phone/pkg/db"
	normalise "phone/pkg/phone"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "elozev"
	password = "1234"
	dbname   = "gophercises_test"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", host, port, user, password)
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	must(phonedb.Migrate("postgres", psqlInfo))
	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)

	err = db.Seed()
	must(err)

	phones, err := db.AllPhones()
	must(err)

	// number, err := getPhone(db, id)
	// must(err)

	// fmt.Printf("id=%d, number=%s\n", id, number)

	for _, p := range phones {
		fmt.Printf("Working on %+v\n", p)

		number := normalise.Clean(p.Number)
		if number != p.Number {
			fmt.Println("Updating or removing: ", number)
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				// delete this number
				must(db.DeletePhone(p.ID))
			} else {
				// update this number
				p.Number = number
				must(db.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No changes required")
		}

	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
