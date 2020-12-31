package app

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

const width = 800
const height = 600

type app struct {
	physicalDevice vk.PhysicalDevice
	instance       vk.Instance
	config         AppConfig
	debugMessenger vk.DebugReportCallback
	logicalDevice  vk.Device
	windowSurface  vk.Surface
}

type AppConfig struct {
	EnableValidationLayers   bool
	ValidationLayers         []string
	RequiredDeviceExtensions []string
}

func New(config AppConfig) *app {
	app := &app{config: config}
	return app
}

func (a *app) Run() error {
	var err error

	win, err := a.initWindow()
	if err != nil {
		return err
	}

	err = a.initVulkan(win)
	if err != nil {
		return err
	}

	a.mainLoop(win)
	a.cleanup(win)

	return nil
}

func (a *app) initWindow() (*glfw.Window, error) {
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

func (a *app) mainLoop(win *glfw.Window) {
	for !win.ShouldClose() {
		glfw.PollEvents()
	}
}
func (a *app) cleanup(win *glfw.Window) {
	vk.DestroyDevice(a.logicalDevice, nil)
	if a.config.EnableValidationLayers {
		vk.DestroyDebugReportCallback(a.instance, a.debugMessenger, nil)
	}
	vk.DestroySurface(a.instance, a.windowSurface, nil)
	vk.DestroyInstance(a.instance, nil)
	win.Destroy()
	glfw.Terminate()
}
