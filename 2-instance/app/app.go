package app

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

const width = 800
const height = 600

type App struct {
	instance vk.Instance
}

func New() *App {
	app := &App{}
	return app
}

func (a *App) Run() error {
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
func (a *App) initVulkan(win *glfw.Window) error {

	procAddr := glfw.GetVulkanGetInstanceProcAddress()
	if procAddr == nil {
		return fmt.Errorf("GetInstanceProcAddress is nil")
	}
	vk.SetGetInstanceProcAddr(procAddr)

	err := vk.Init()
	if err != nil {
		return err
	}

	return a.createInstance(win)
}
func (a *App) mainLoop(win *glfw.Window) {
	for !win.ShouldClose() {
		glfw.PollEvents()
	}
}
func (a *App) cleanup(win *glfw.Window) {
	vk.DestroyInstance(a.instance, nil)
	win.Destroy()
	glfw.Terminate()
}

func (a *App) createInstance(win *glfw.Window) error {
	glfwExtensions := win.GetRequiredInstanceExtensions()

	//var extensionCount uint32
	//vk.EnumerateInstanceExtensionProperties("", &extensionCount, nil)
	//extensionProperties := make([]vk.ExtensionProperties, extensionCount)
	//vk.EnumerateInstanceExtensionProperties("", &extensionCount, extensionProperties)
	//
	//supportedExtensions := make(map[string]bool)
	//for _, extensionProperty := range extensionProperties {
	//	supportedExtensions[string(extensionProperty.ExtensionName[:])] = true
	//}
	//
	//for _, glfwExtension := range glfwExtensions {
	//	if !supportedExtensions[glfwExtension] {
	//		return fmt.Errorf("glfwExtension - " + glfwExtension + " - is not supported by vulkan")
	//	}
	//}

	applicationInfo := vk.ApplicationInfo{
		SType:              vk.StructureTypeApplicationInfo,
		PApplicationName:   "Hello Triangle",
		ApplicationVersion: vk.MakeVersion(1, 0, 0),
		PEngineName:        "No Engine",
		EngineVersion:      vk.MakeVersion(1, 0, 0),
		ApiVersion:         vk.MakeVersion(1, 0, 0),
	}

	instanceCreateInfo := vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo:        &applicationInfo,
		EnabledExtensionCount:   uint32(len(glfwExtensions)),
		PpEnabledExtensionNames: glfwExtensions,
	}

	res := vk.CreateInstance(&instanceCreateInfo, nil, &a.instance)
	if res != vk.Success {
		return fmt.Errorf("failed to create instance")
	}

	return nil
}
