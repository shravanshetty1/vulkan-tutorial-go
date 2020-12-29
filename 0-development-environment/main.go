package main

import (
	"fmt"
	"github.com/vulkan-go/vulkan"
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	procAddr := glfw.GetVulkanGetInstanceProcAddress()
	if procAddr == nil {
		panic("GetInstanceProcAddress is nil")
	}
	vulkan.SetGetInstanceProcAddr(procAddr)

	err := glfw.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = vulkan.Init()
	if err != nil {
		log.Fatal(err)
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	win, err := glfw.CreateWindow(800, 600, "Vulkan Window", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	var extensionCount uint32 = 0
	vulkan.EnumerateInstanceExtensionProperties("", &extensionCount, nil)
	fmt.Println(fmt.Sprint(extensionCount)+" extensions supported")

	for !win.ShouldClose() {
		glfw.PollEvents()
	}

	win.Destroy()
	glfw.Terminate()
}
