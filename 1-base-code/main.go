package main

import (
	"log"

	"github.com/vulkan-go/vulkan"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func NewApp() *App {
	app := &App{}
	return app
}

type App struct {
}

func (a *App) Run() error {
	var err error

	win, err := a.initWindow()
	if err != nil {
		return err
	}

	err = a.initVulkan()
	if err != nil {
		return err
	}

	a.mainLoop(win)
	a.cleanup(win)

	return nil
}

const width = 800
const height = 600

func (a *App) initWindow() (*glfw.Window, error) {
	err := glfw.Init()
	if err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	win, err := glfw.CreateWindow(width, height, "Vulkan", nil, nil)
	if err != nil {
		return nil, err
	}

	return win, nil

}
func (a *App) initVulkan() error {

	procAddr := glfw.GetVulkanGetInstanceProcAddress()
	if procAddr == nil {
		panic("GetInstanceProcAddress is nil")
	}
	vulkan.SetGetInstanceProcAddr(procAddr)

	err := vulkan.Init()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
func (a *App) mainLoop(win *glfw.Window) {
	for !win.ShouldClose() {
		glfw.PollEvents()
	}
}
func (a *App) cleanup(win *glfw.Window) {
	win.Destroy()
	glfw.Terminate()
}

func main() {
	app := NewApp()

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
