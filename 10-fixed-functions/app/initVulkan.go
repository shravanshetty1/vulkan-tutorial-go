package app

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

func (a *app) initVulkan() error {

	procAddr := glfw.GetVulkanGetInstanceProcAddress()
	if procAddr == nil {
		return fmt.Errorf("GetInstanceProcAddress is nil")
	}
	vk.SetGetInstanceProcAddr(procAddr)

	err := vk.Init()
	if err != nil {
		return err
	}

	err = a.createInstance()
	if err != nil {
		return err
	}

	if a.config.EnableValidationLayers {
		err = a.setupDebugMessenger()
		if err != nil {
			return err
		}
	}

	err = a.createWindowSurface()
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

	err = a.createSwapChain()
	if err != nil {
		return err
	}

	err = a.createImageViews()
	if err != nil {
		return err
	}

	err = a.createGraphicsPipeline()
	if err != nil {
		return err
	}

	return nil
}

func (a *app) createGraphicsPipeline() error {

	//_, fileName, _, _ := runtime.Caller(1)
	//fileName, err := filepath.Abs(fileName)
	//if err != nil {
	//	return err
	//}
	//
	//fragCode, err := ioutil.ReadFile(filepath.Join(filepath.Dir(fileName), "../shaders/frag.spv"))
	//if err != nil {
	//	return err
	//}
	//
	//vertCode, err := ioutil.ReadFile(filepath.Join(filepath.Dir(fileName), "../shaders/vert.spv"))
	//if err != nil {
	//	return err
	//}
	//
	//fragModule, err := a.createShaderModule(fragCode)
	//if err != nil {
	//	return err
	//}
	//vertModule, err := a.createShaderModule(vertCode)
	//if err != nil {
	//	return err
	//}
	//
	//vertStageCreateInfo := vk.PipelineShaderStageCreateInfo{
	//	SType:  vk.StructureTypePipelineShaderStageCreateInfo,
	//	Stage:  vk.ShaderStageVertexBit,
	//	Module: vertModule,
	//	PName:  "main",
	//}
	//fragStageCreateInfo := vk.PipelineShaderStageCreateInfo{
	//	SType:  vk.StructureTypePipelineShaderStageCreateInfo,
	//	Stage:  vk.ShaderStageFragmentBit,
	//	Module: fragModule,
	//	PName:  "main",
	//}
	//
	//shaderStages := []vk.PipelineShaderStageCreateInfo{vertStageCreateInfo, fragStageCreateInfo}
	//
	//vertexInputStateCreateInfo := vk.PipelineVertexInputStateCreateInfo{
	//	SType: vk.StructureTypePipelineVertexInputStateCreateInfo,
	//}
	//
	//inputAssemblyStateCreateInfo := vk.PipelineInputAssemblyStateCreateInfo{
	//	SType:                  vk.StructureTypePipelineInputAssemblyStateCreateInfo,
	//	Topology:               vk.PrimitiveTopologyTriangleList,
	//	PrimitiveRestartEnable: vk.False,
	//}
	//
	//viewports := []vk.Viewport{{
	//	X:        0,
	//	Y:        0,
	//	Width:    float32(a.swapChainExtent.Width),
	//	Height:   float32(a.swapChainExtent.Height),
	//	MinDepth: 0,
	//	MaxDepth: 1,
	//}}
	//
	//scissors := []vk.Rect2D{{
	//	Offset: vk.Offset2D{
	//		X: 0,
	//		Y: 0,
	//	},
	//	Extent: a.swapChainExtent,
	//}}
	//
	//viewportStateCreateInfo := vk.PipelineViewportStateCreateInfo{
	//	SType:         vk.StructureTypePipelineViewportStateCreateInfo,
	//	ViewportCount: uint32(len(viewports)),
	//	PViewports:    viewports,
	//	ScissorCount:  uint32(len(scissors)),
	//	PScissors:     scissors,
	//}
	//
	//rasterizer := vk.PipelineRasterizationStateCreateInfo{
	//	SType:                   vk.StructureTypePipelineRasterizationStateCreateInfo,
	//	PNext:                   nil,
	//	Flags:                   0,
	//	DepthClampEnable:        vk.False,
	//	RasterizerDiscardEnable: vk.False,
	//	PolygonMode:             vk.PolygonModeFill,
	//	CullMode:                vk.CullModeFlags(vk.CullModeBackBit),
	//	FrontFace:               vk.FrontFaceClockwise,
	//	DepthBiasEnable:         vk.False,
	//	DepthBiasConstantFactor: 0,
	//	DepthBiasClamp:          0,
	//	DepthBiasSlopeFactor:    0,
	//	LineWidth:               1,
	//}
	//
	//multisamplingCreateInfo := vk.PipelineMultisampleStateCreateInfo{
	//	SType:                 vk.StructureTypePipelineMultisampleStateCreateInfo,
	//	PNext:                 nil,
	//	Flags:                 0,
	//	RasterizationSamples:  vk.SampleCount1Bit,
	//	SampleShadingEnable:   vk.False,
	//	MinSampleShading:      1,
	//	PSampleMask:           nil,
	//	AlphaToCoverageEnable: vk.False,
	//	AlphaToOneEnable:      vk.False,
	//}
	//
	//colorBlendAttachmentStates := []vk.PipelineColorBlendAttachmentState{{
	//	BlendEnable:         vk.False,
	//	SrcColorBlendFactor: vk.BlendFactorOne,
	//	DstColorBlendFactor: vk.BlendFactorZero,
	//	ColorBlendOp:        vk.BlendOpAdd,
	//	SrcAlphaBlendFactor: vk.BlendFactorOne,
	//	DstAlphaBlendFactor: vk.BlendFactorZero,
	//	AlphaBlendOp:        vk.BlendOpAdd,
	//	ColorWriteMask:      vk.ColorComponentFlags(vk.ColorComponentRBit | vk.ColorComponentGBit | vk.ColorComponentBBit | vk.ColorComponentABit),
	//}}
	//
	//colorBlendingCreateInfo := vk.PipelineColorBlendStateCreateInfo{
	//	SType:           vk.StructureTypePipelineColorBlendStateCreateInfo,
	//	PNext:           nil,
	//	Flags:           0,
	//	LogicOpEnable:   vk.False,
	//	LogicOp:         vk.LogicOpCopy,
	//	AttachmentCount: uint32(len(colorBlendAttachmentStates)),
	//	PAttachments:    colorBlendAttachmentStates,
	//	BlendConstants:  [4]float32{0, 0, 0, 0},
	//}
	//
	//dynamicStates := []vk.DynamicState{
	//	vk.DynamicStateViewport,
	//	vk.DynamicStateLineWidth,
	//}
	//
	//dynamicStateCreateInfo := vk.PipelineDynamicStateCreateInfo{
	//	SType:             vk.StructureTypePipelineDynamicStateCreateInfo,
	//	PNext:             nil,
	//	Flags:             0,
	//	DynamicStateCount: uint32(len(dynamicStates)),
	//	PDynamicStates:    dynamicStates,
	//}

	pipelineLayoutCreateInfo := vk.PipelineLayoutCreateInfo{
		SType: vk.StructureTypePipelineLayoutCreateInfo,
	}

	var pipelineLayout vk.PipelineLayout
	err := vk.Error(vk.CreatePipelineLayout(a.logicalDevice, &pipelineLayoutCreateInfo, nil, &pipelineLayout))
	if err != nil {
		return err
	}

	a.pipelineLayout = pipelineLayout

	//vk.DestroyShaderModule(a.logicalDevice, fragModule, nil)
	//vk.DestroyShaderModule(a.logicalDevice, vertModule, nil)
	return nil
}

func (a *app) createShaderModule(code []byte) (vk.ShaderModule, error) {
	createInfo := vk.ShaderModuleCreateInfo{
		SType:    vk.StructureTypeShaderModuleCreateInfo,
		PNext:    nil,
		Flags:    0,
		CodeSize: uint(len(code)),
		PCode:    repackUint32(code),
	}

	var shaderModule vk.ShaderModule
	err := vk.Error(vk.CreateShaderModule(a.logicalDevice, &createInfo, nil, &shaderModule))
	if err != nil {
		return nil, fmt.Errorf("could not create shader module - " + err.Error())
	}

	return shaderModule, nil
}

func repackUint32(data []byte) []uint32 {
	buf := make([]uint32, len(data)/4)
	vk.Memcopy(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buf)).Data), data)
	return buf
}

type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

func (a *app) createImageViews() error {
	a.swapChainImageViews = make([]vk.ImageView, len(a.swapChainImages))

	for i, image := range a.swapChainImages {
		createInfo := vk.ImageViewCreateInfo{
			SType:    vk.StructureTypeImageViewCreateInfo,
			PNext:    nil,
			Flags:    0,
			Image:    image,
			ViewType: vk.ImageViewType2d,
			Format:   a.swapChainImageFormat,
			Components: vk.ComponentMapping{
				R: vk.ComponentSwizzleIdentity,
				G: vk.ComponentSwizzleIdentity,
				B: vk.ComponentSwizzleIdentity,
				A: vk.ComponentSwizzleIdentity,
			},
			SubresourceRange: vk.ImageSubresourceRange{
				AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
				BaseMipLevel:   0,
				LevelCount:     1,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
		}
		var imageView vk.ImageView
		err := vk.Error(vk.CreateImageView(a.logicalDevice, &createInfo, nil, &imageView))
		if err != nil {
			return err
		}

		a.swapChainImageViews[i] = imageView
	}

	return nil
}

func (a *app) createSwapChain() error {
	swapChainSupport := querySwapChainSupport(a.physicalDevice, a.windowSurface)
	swapChainSupport.capabilities.Deref()
	swapChainSupport.capabilities.Free()

	surfaceFormat := chooseSwapSurfaceFormat(swapChainSupport.surfaceFormats...)
	presentationMode := chooseSwapPresentMode(swapChainSupport.presentationModes...)
	swapExtent := chooseSwapExtent(swapChainSupport.capabilities, a.window)
	imageCount := swapChainSupport.capabilities.MinImageCount + 1

	if swapChainSupport.capabilities.MaxImageCount > 0 && imageCount > swapChainSupport.capabilities.MaxImageCount {
		imageCount = swapChainSupport.capabilities.MaxImageCount
	}

	createInfo := vk.SwapchainCreateInfo{
		SType:            vk.StructureTypeSwapchainCreateInfo,
		Surface:          a.windowSurface,
		MinImageCount:    imageCount,
		ImageFormat:      surfaceFormat.Format,
		ImageColorSpace:  surfaceFormat.ColorSpace,
		ImageExtent:      swapExtent,
		ImageArrayLayers: 1,
		ImageUsage:       vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit),
		PreTransform:     swapChainSupport.capabilities.CurrentTransform,
		CompositeAlpha:   vk.CompositeAlphaOpaqueBit,
		PresentMode:      presentationMode,
		Clipped:          vk.True,
		OldSwapchain:     vk.NullSwapchain,
	}

	indices := findQueueFamilies(a.physicalDevice, a.windowSurface)
	queueFamilies := []uint32{*indices.presentFamily, *indices.graphicsFamily}

	if *indices.graphicsFamily != *indices.presentFamily {
		createInfo.ImageSharingMode = vk.SharingModeConcurrent
		createInfo.QueueFamilyIndexCount = uint32(len(queueFamilies))
		createInfo.PQueueFamilyIndices = queueFamilies
	} else {
		createInfo.ImageSharingMode = vk.SharingModeExclusive
	}

	var swapChain vk.Swapchain
	err := vk.Error(vk.CreateSwapchain(a.logicalDevice, &createInfo, nil, &swapChain))
	if err != nil {
		return err
	}

	a.swapChain = swapChain

	var imagesCount uint32
	vk.GetSwapchainImages(a.logicalDevice, a.swapChain, &imagesCount, nil)
	a.swapChainImages = make([]vk.Image, imageCount)
	vk.GetSwapchainImages(a.logicalDevice, a.swapChain, &imagesCount, a.swapChainImages)

	a.swapChainExtent = swapExtent
	a.swapChainImageFormat = surfaceFormat.Format

	return nil
}

func (a *app) createInstance() error {

	requiredExtensions := a.window.GetRequiredInstanceExtensions()
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

func (a *app) createWindowSurface() error {

	surfaceAddr, err := a.window.CreateWindowSurface(a.instance, nil)
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

	return nil
}

func (a *app) isDeviceSuitable(device vk.PhysicalDevice) bool {

	if !checkDeviceExtensionsSupport(device, a.config.RequiredDeviceExtensions) {
		return false
	}

	swapChainSupport := querySwapChainSupport(device, a.windowSurface)
	if len(swapChainSupport.surfaceFormats) == 0 || len(swapChainSupport.presentationModes) == 0 {
		return false
	}

	indices := findQueueFamilies(device, a.windowSurface)
	if !indices.isComplete() {
		return false
	}

	return true
}

func chooseSwapSurfaceFormat(surfaceFormats ...vk.SurfaceFormat) vk.SurfaceFormat {
	if len(surfaceFormats) < 1 {
		return vk.SurfaceFormat{}
	}

	for _, surfaceFormat := range surfaceFormats {
		surfaceFormat.Deref()
		surfaceFormat.Free()

		if surfaceFormat.Format == vk.FormatB8g8r8a8Srgb && surfaceFormat.ColorSpace == vk.ColorspaceSrgbNonlinear {
			return surfaceFormat
		}
	}

	return surfaceFormats[0]
}

func chooseSwapPresentMode(presentModes ...vk.PresentMode) vk.PresentMode {
	if len(presentModes) < 1 {
		return 0
	}

	for _, presentMode := range presentModes {
		if presentMode == vk.PresentModeMailbox {
			return presentMode
		}
	}

	return vk.PresentModeFifo
}

func chooseSwapExtent(surfaceCapabilities vk.SurfaceCapabilities, win *glfw.Window) vk.Extent2D {
	surfaceCapabilities.Deref()
	surfaceCapabilities.Free()
	surfaceCapabilities.CurrentExtent.Deref()
	surfaceCapabilities.CurrentExtent.Free()
	surfaceCapabilities.MaxImageExtent.Deref()
	surfaceCapabilities.MaxImageExtent.Free()
	surfaceCapabilities.MinImageExtent.Deref()
	surfaceCapabilities.MinImageExtent.Free()

	//if surfaceCapabilities.CurrentExtent.Width != vk.MaxUint32 {
	//	return surfaceCapabilities.CurrentExtent
	//}

	w, h := win.GetFramebufferSize()

	actualExtent := vk.Extent2D{
		Width:  uint32(w),
		Height: uint32(h),
	}

	if actualExtent.Width > surfaceCapabilities.MaxImageExtent.Width {
		actualExtent.Width = surfaceCapabilities.MaxImageExtent.Width
	}
	if actualExtent.Width < surfaceCapabilities.MinImageExtent.Width {
		actualExtent.Width = surfaceCapabilities.MinImageExtent.Width
	}

	if actualExtent.Height > surfaceCapabilities.MaxImageExtent.Height {
		actualExtent.Height = surfaceCapabilities.MaxImageExtent.Height
	}
	if actualExtent.Height < surfaceCapabilities.MinImageExtent.Height {
		actualExtent.Height = surfaceCapabilities.MinImageExtent.Height
	}

	return actualExtent
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
