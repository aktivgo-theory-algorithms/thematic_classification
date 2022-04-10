package main

import (
	"DatasetCreater/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}
