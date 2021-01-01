package app

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

const width = 800
const height = 600

type app struct {
	window               *glfw.Window
	physicalDevice       vk.PhysicalDevice
	instance             vk.Instance
	config               AppConfig
	debugMessenger       vk.DebugReportCallback
	logicalDevice        vk.Device
	windowSurface        vk.Surface
	swapChain            vk.Swapchain
	swapChainImages      []vk.Image
	swapChainExtent      vk.Extent2D
	swapChainImageFormat vk.Format
	swapChainImageViews  []vk.ImageView
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
	err = a.initWindow()
	if err != nil {
		return err
	}

	err = a.initVulkan()
	if err != nil {
		return err
	}

	a.mainLoop()
	a.cleanup()

	return nil
}

func (a *app) initWindow() error {
	err := glfw.Init()
	if err != nil {
		return err
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	win, err := glfw.CreateWindow(width, height, "Vulkan", nil, nil)
	if err != nil {
		return err
	}

	a.window = win

	return nil
}

func (a *app) mainLoop() {
	for !a.window.ShouldClose() {
		glfw.PollEvents()
	}
}
func (a *app) cleanup() {
	for i := range a.swapChainImageViews {
		vk.DestroyImageView(a.logicalDevice, a.swapChainImageViews[i], nil)
	}
	vk.DestroySwapchain(a.logicalDevice, a.swapChain, nil)
	vk.DestroyDevice(a.logicalDevice, nil)
	if a.config.EnableValidationLayers {
		vk.DestroyDebugReportCallback(a.instance, a.debugMessenger, nil)
	}
	vk.DestroySurface(a.instance, a.windowSurface, nil)
	vk.DestroyInstance(a.instance, nil)
	a.window.Destroy()
	glfw.Terminate()
}
