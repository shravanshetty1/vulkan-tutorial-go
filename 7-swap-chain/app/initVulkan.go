package app

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

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
		SType:                   vk.StructureTypeDeviceCreateInfo,
		QueueCreateInfoCount:    uint32(len(queueCreateInfos)),
		PQueueCreateInfos:       queueCreateInfos,
		EnabledExtensionCount:   uint32(len(a.config.RequiredDeviceExtensions)),
		PpEnabledExtensionNames: a.config.RequiredDeviceExtensions,
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

func (a *app) isDeviceSuitable(device vk.PhysicalDevice) bool {

	if !checkDeviceExtensionsSupport(device, a.config.RequiredDeviceExtensions) {
		return false
	}

	indices := findQueueFamilies(device, a.windowSurface)
	if !indices.isComplete() {
		return false
	}

	return true
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

func (a *app) pickPhysicalDevice() error {

	var deviceCount uint32
	vk.EnumeratePhysicalDevices(a.instance, &deviceCount, nil)
	if deviceCount == 0 {
		return fmt.Errorf("failed to find gpus with vulkan support")
	}

	physicalDevices := make([]vk.PhysicalDevice, deviceCount)
	vk.EnumeratePhysicalDevices(a.instance, &deviceCount, physicalDevices)

	for _, physicalDevice := range physicalDevices {
		if a.isDeviceSuitable(physicalDevice) {
			a.physicalDevice = physicalDevice
			break
		}
	}

	if unsafe.Pointer(a.physicalDevice) == vk.NullHandle {
		return fmt.Errorf("failed to find a suitable gpu")
	}

	return nil
}
