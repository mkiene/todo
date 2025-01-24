package main

import (
	"os"
)

func handle_input() {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "add":
			err := create_task()
			if err != nil {
				// log.Fatal(err)
			}
		}
	}
}
