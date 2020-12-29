package app

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

const width = 800
const height = 600

type app struct {
	instance vk.Instance
	config   AppConfig
}

type AppConfig struct {
	EnableValidationLayers bool
	ValidationLayers       []string
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
func (a *app) initVulkan(win *glfw.Window) error {

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
func (a *app) mainLoop(win *glfw.Window) {
	for !win.ShouldClose() {
		glfw.PollEvents()
	}
}
func (a *app) cleanup(win *glfw.Window) {
	vk.DestroyInstance(a.instance, nil)
	win.Destroy()
	glfw.Terminate()
}

func (a *app) createInstance(win *glfw.Window) error {

	requiredExtensions := win.GetRequiredInstanceExtensions()
	if a.config.EnableValidationLayers {
		requiredExtensions = append(requiredExtensions, "VK_EXT_debug_report\x00")
	}
	err := checkExtensionSupport(requiredExtensions)
	if err != nil {
		return err
	}

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
		EnabledExtensionCount:   uint32(len(requiredExtensions)),
		PpEnabledExtensionNames: requiredExtensions,
	}

	if a.config.EnableValidationLayers {
		err := checkValidationLayerSupport(a.config.ValidationLayers)
		if err != nil {
			return err
		}
		instanceCreateInfo.PpEnabledLayerNames = a.config.ValidationLayers
		instanceCreateInfo.EnabledLayerCount = uint32(len(a.config.ValidationLayers))
	}

	var instance vk.Instance
	res := vk.CreateInstance(&instanceCreateInfo, nil, &instance)
	if res != vk.Success {
		return fmt.Errorf("failed to create instance")
	}

	a.instance = instance

	return nil
}

func checkExtensionSupport(requiredExtensions []string) error {
	var extensionCount uint32
	vk.EnumerateInstanceExtensionProperties("", &extensionCount, nil)
	extensionProperties := make([]vk.ExtensionProperties, extensionCount)
	vk.EnumerateInstanceExtensionProperties("", &extensionCount, extensionProperties)

	supportedExtensions := make(map[string]bool)
	for _, extensionProperty := range extensionProperties {
		extensionProperty.Deref()
		supportedExtensions[vk.ToString(extensionProperty.ExtensionName[:])] = true
		extensionProperty.Free()
	}

	for _, requiredExtension := range requiredExtensions {
		requiredExtension = strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		}, requiredExtension)
		if !supportedExtensions[requiredExtension] {
			return fmt.Errorf(requiredExtension + " - is not a supported extension")
		}
	}

	return nil
}

func checkValidationLayerSupport(requiredLayers []string) error {
	var layerCount uint32
	vk.EnumerateInstanceLayerProperties(&layerCount, nil)
	layerProperties := make([]vk.LayerProperties, layerCount)
	vk.EnumerateInstanceLayerProperties(&layerCount, layerProperties)

	supportedLayers := make(map[string]bool)
	for _, layerProperty := range layerProperties {
		layerProperty.Deref()
		supportedLayers[vk.ToString(layerProperty.LayerName[:])] = true
		layerProperty.Free()
	}

	for _, requiredLayer := range requiredLayers {
		requiredLayer = strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		}, requiredLayer)
		if !supportedLayers[requiredLayer] {
			return fmt.Errorf(requiredLayer + " - is not a supported layer")
		}
	}

	return nil
}
