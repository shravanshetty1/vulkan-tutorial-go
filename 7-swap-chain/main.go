package main

import (
	"log"
	"os"
	"vulkan-tutorial-go/7-swap-chain/app"
)

func main() {
	profile := os.Getenv("PROFILE")

	var enableValidationLayers bool
	if profile != "prod" {
		enableValidationLayers = true
	}

	a := app.New(app.AppConfig{EnableValidationLayers: enableValidationLayers, ValidationLayers: []string{
		"VK_LAYER_KHRONOS_validation\x00",
	}})

	err := a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
