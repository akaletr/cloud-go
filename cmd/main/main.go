package main

import (
	"log"

	"cmd/main/main.go/internal/app"
)

func main() {
	myApp, err := app.New()
	if err != nil {
		log.Fatalln(err)
	}

	defer log.Fatalln(myApp.Stop())
	log.Fatalln(myApp.Start())
}
