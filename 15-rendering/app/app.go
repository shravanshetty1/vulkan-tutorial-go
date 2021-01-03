package app

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

const width = 800
const height = 600

type app struct {
	window                  *glfw.Window
	physicalDevice          vk.PhysicalDevice
	instance                vk.Instance
	config                  AppConfig
	debugMessenger          vk.DebugReportCallback
	logicalDevice           vk.Device
	windowSurface           vk.Surface
	graphicsQueue           vk.Queue
	presentQueue            vk.Queue
	swapChain               vk.Swapchain
	swapChainImages         []vk.Image
	swapChainExtent         vk.Extent2D
	swapChainImageFormat    vk.Format
	swapChainImageViews     []vk.ImageView
	renderPass              vk.RenderPass
	pipelineLayout          vk.PipelineLayout
	graphicsPipeline        vk.Pipeline
	swapChainFrameBuffers   []vk.Framebuffer
	commandPool             vk.CommandPool
	commandBuffers          []vk.CommandBuffer
	imageAvailableSemaphore vk.Semaphore
	renderFinishedSemaphore vk.Semaphore
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

	err = a.mainLoop()
	if err != nil {
		return err
	}
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

func (a *app) mainLoop() error {
	for !a.window.ShouldClose() {
		glfw.PollEvents()
		err := a.drawFrame()
		if err != nil {
			return err
		}
	}

	vk.DeviceWaitIdle(a.logicalDevice)

	return nil
}

func (a *app) drawFrame() error {
	var imageIndex uint32
	vk.AcquireNextImage(a.logicalDevice, a.swapChain, vk.MaxUint64, a.imageAvailableSemaphore, vk.NullFence, &imageIndex)

	waitsemaphores := []vk.Semaphore{a.imageAvailableSemaphore}
	signalsemaphores := []vk.Semaphore{a.renderFinishedSemaphore}
	waitStages := []vk.PipelineStageFlags{vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit)}

	submitInfo := []vk.SubmitInfo{{
		SType:                vk.StructureTypeSubmitInfo,
		PNext:                nil,
		WaitSemaphoreCount:   uint32(len(waitsemaphores)),
		PWaitSemaphores:      waitsemaphores,
		PWaitDstStageMask:    waitStages,
		CommandBufferCount:   1,
		PCommandBuffers:      []vk.CommandBuffer{a.commandBuffers[imageIndex]},
		SignalSemaphoreCount: uint32(len(signalsemaphores)),
		PSignalSemaphores:    signalsemaphores,
	}}

	err := vk.Error(vk.QueueSubmit(a.graphicsQueue, 1, submitInfo, vk.NullFence))
	if err != nil {
		return err
	}

	presentInfo := vk.PresentInfo{
		SType:              vk.StructureTypePresentInfo,
		PNext:              nil,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    signalsemaphores,
		SwapchainCount:     1,
		PSwapchains:        []vk.Swapchain{a.swapChain},
		PImageIndices:      []uint32{imageIndex},
		PResults:           nil,
	}

	vk.QueuePresent(a.presentQueue, &presentInfo)

	return nil
}

func (a *app) cleanup() {
	vk.DestroySemaphore(a.logicalDevice, a.renderFinishedSemaphore, nil)
	vk.DestroySemaphore(a.logicalDevice, a.imageAvailableSemaphore, nil)
	vk.DestroyCommandPool(a.logicalDevice, a.commandPool, nil)
	for _, v := range a.swapChainFrameBuffers {
		vk.DestroyFramebuffer(a.logicalDevice, v, nil)
	}

	vk.DestroyPipeline(a.logicalDevice, a.graphicsPipeline, nil)
	vk.DestroyPipelineLayout(a.logicalDevice, a.pipelineLayout, nil)
	vk.DestroyRenderPass(a.logicalDevice, a.renderPass, nil)
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
