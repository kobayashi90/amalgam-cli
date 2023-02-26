package main

import (
	"log"
	"os"
)

func main() {
	app := CmdApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
