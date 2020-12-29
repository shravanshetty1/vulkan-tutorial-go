package main

import (
	"log"
	"vulkan-tutorial-go/2-instance/app"
)

func main() {

	a := app.New()

	err := a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
