package app

import (
	"fmt"
	"log"
	"strings"
	"unicode"
	"unsafe"

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

	err = a.createInstance(win)
	if err != nil {
		return err
	}

	if a.config.EnableValidationLayers {
		err = a.setupDebugMessenger()
		if err != nil {
			return err
		}
	}

	err = a.createWindowSurface(win)
	if err != nil {
		return err
	}

	err = a.pickPhysicalDevice()
	if err != nil {
		return err
	}

	err = a.createLogicalDevice()
	if err != nil {
		return err
	}

	return nil
}

func (a *app) createWindowSurface(win *glfw.Window) error {

	surfaceAddr, err := win.CreateWindowSurface(a.instance, nil)
	if err != nil {
		return err
	}

	a.windowSurface = vk.SurfaceFromPointer(surfaceAddr)

	return nil
}

func (a *app) createLogicalDevice() error {
	indices := findQueueFamilies(a.physicalDevice, a.windowSurface)

	uniqueQueueFamily := map[uint32]bool{
		*indices.graphicsFamily: true,
		*indices.presentFamily:  true,
	}

	var queueCreateInfos []vk.DeviceQueueCreateInfo
	for queueFamilyindex := range uniqueQueueFamily {
		queueCreateInfos = append(queueCreateInfos, vk.DeviceQueueCreateInfo{
			SType:            vk.StructureTypeDeviceQueueCreateInfo,
			QueueFamilyIndex: queueFamilyindex,
			QueueCount:       1,
			PQueuePriorities: []float32{1},
		})
	}

	//deviceFeatures := []vk.PhysicalDeviceFeatures{}

	deviceCreateInfo := vk.DeviceCreateInfo{
		SType:                vk.StructureTypeDeviceCreateInfo,
		QueueCreateInfoCount: uint32(len(queueCreateInfos)),
		PQueueCreateInfos:    queueCreateInfos,
	}

	if a.config.EnableValidationLayers {
		deviceCreateInfo.EnabledLayerCount = uint32(len(a.config.ValidationLayers))
		deviceCreateInfo.PpEnabledLayerNames = a.config.ValidationLayers
	}

	var device vk.Device
	if vk.CreateDevice(a.physicalDevice, &deviceCreateInfo, nil, &device) != vk.Success {
		return fmt.Errorf("could not create logical device")
	}

	a.logicalDevice = device

	var graphicsQueue vk.Queue
	vk.GetDeviceQueue(device, *indices.graphicsFamily, 0, &graphicsQueue)
	var presentQueue vk.Queue
	vk.GetDeviceQueue(device, *indices.presentFamily, 0, &presentQueue)

	return nil
}

func isDeviceSuitable(device vk.PhysicalDevice, surface vk.Surface) bool {

	var deviceProperties vk.PhysicalDeviceProperties
	var deviceFeatures vk.PhysicalDeviceFeatures
	vk.GetPhysicalDeviceProperties(device, &deviceProperties)
	vk.GetPhysicalDeviceFeatures(device, &deviceFeatures)

	indices := findQueueFamilies(device, surface)
	if !indices.isComplete() {
		return false
	}

	return true
}

type queueFamilyIndices struct {
	graphicsFamily *uint32
	presentFamily  *uint32
}

func (q *queueFamilyIndices) isComplete() bool {
	return q.graphicsFamily != nil && q.presentFamily != nil
}

func findQueueFamilies(device vk.PhysicalDevice, surface vk.Surface) queueFamilyIndices {
	var indices queueFamilyIndices

	var propCount uint32
	vk.GetPhysicalDeviceQueueFamilyProperties(device, &propCount, nil)

	properties := make([]vk.QueueFamilyProperties, propCount)
	vk.GetPhysicalDeviceQueueFamilyProperties(device, &propCount, properties)

	for i, property := range properties {
		property.Deref()
		queueFlags := property.QueueFlags
		property.Free()

		if (uint32(queueFlags) & uint32(vk.QueueGraphicsBit)) != 0 {
			tmp := uint32(i)
			indices.graphicsFamily = &tmp
		}

		var isSupported vk.Bool32
		vk.GetPhysicalDeviceSurfaceSupport(device, uint32(i), surface, &isSupported)
		if isSupported == vk.True {
			tmp := uint32(i)
			indices.presentFamily = &tmp
		}

		if indices.isComplete() {
			break
		}
	}

	return indices
}

func (a *app) pickPhysicalDevice() error {

	var deviceCount uint32
	vk.EnumeratePhysicalDevices(a.instance, &deviceCount, nil)
	if deviceCount == 0 {
		return fmt.Errorf("failed to find gpus with vulkan support")
	}

	physicalDevices := make([]vk.PhysicalDevice, deviceCount)
	vk.EnumeratePhysicalDevices(a.instance, &deviceCount, physicalDevices)

	for _, physicalDevice := range physicalDevices {
		if isDeviceSuitable(physicalDevice, a.windowSurface) {
			a.physicalDevice = physicalDevice
			break
		}
	}

	if unsafe.Pointer(a.physicalDevice) == vk.NullHandle {
		return fmt.Errorf("failed to find a suitable gpu")
	}

	return nil
}

func defaultDebugCreateInfo() vk.DebugReportCallbackCreateInfo {
	return vk.DebugReportCallbackCreateInfo{
		SType: vk.StructureTypeDebugReportCallbackCreateInfo,
		Flags: vk.DebugReportFlags(vk.DebugReportErrorBit | vk.DebugReportWarningBit),
		PfnCallback: func(flags vk.DebugReportFlags, objectType vk.DebugReportObjectType, object uint64, location uint, messageCode int32, pLayerPrefix string, pMessage string, pUserData unsafe.Pointer) vk.Bool32 {
			switch {
			case flags&vk.DebugReportFlags(vk.DebugReportErrorBit) != 0:
				log.Printf("[ERROR %d] %s on layer %s", messageCode, pMessage, pLayerPrefix)
			case flags&vk.DebugReportFlags(vk.DebugReportWarningBit) != 0:
				log.Printf("[WARN %d] %s on layer %s", messageCode, pMessage, pLayerPrefix)
			default:
				log.Printf("[WARN] unknown debug message %d (layer %s)", messageCode, pLayerPrefix)
			}
			return vk.Bool32(vk.False)
		},
	}
}

func (a *app) setupDebugMessenger() error {
	dbgCreateInfo := defaultDebugCreateInfo()
	var dbg vk.DebugReportCallback
	err := vk.Error(vk.CreateDebugReportCallback(a.instance, &dbgCreateInfo, nil, &dbg))
	if err != nil {
		err = fmt.Errorf("vk.CreateDebugReportCallback failed with %s", err)
		return err
	}
	a.debugMessenger = dbg
	return nil
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

	dbgCreateInfo := defaultDebugCreateInfo()
	instanceCreateInfo := vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo:        &applicationInfo,
		EnabledExtensionCount:   uint32(len(requiredExtensions)),
		PpEnabledExtensionNames: requiredExtensions,
		PNext:                   unsafe.Pointer(dbgCreateInfo.Ref()),
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
