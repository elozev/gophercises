package cyoa

import "log"

var Check = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
