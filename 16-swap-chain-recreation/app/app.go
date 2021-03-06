package app

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

const width = 800
const height = 600
const maxFramesInFlight = 2

type app struct {
	window                   *glfw.Window
	physicalDevice           vk.PhysicalDevice
	instance                 vk.Instance
	config                   AppConfig
	debugMessenger           vk.DebugReportCallback
	logicalDevice            vk.Device
	windowSurface            vk.Surface
	graphicsQueue            vk.Queue
	presentQueue             vk.Queue
	swapChain                vk.Swapchain
	swapChainImages          []vk.Image
	swapChainExtent          vk.Extent2D
	swapChainImageFormat     vk.Format
	swapChainImageViews      []vk.ImageView
	renderPass               vk.RenderPass
	pipelineLayout           vk.PipelineLayout
	graphicsPipeline         vk.Pipeline
	swapChainFrameBuffers    []vk.Framebuffer
	commandPool              vk.CommandPool
	commandBuffers           []vk.CommandBuffer
	imageAvailableSemaphores []vk.Semaphore
	renderFinishedSemaphores []vk.Semaphore
	inFlightFences           []vk.Fence
	imagesInFlight           []vk.Fence
	currentFrame             int
	frameBufferResized       bool
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
	glfw.WindowHint(glfw.Resizable, glfw.True)
	win, err := glfw.CreateWindow(width, height, "Vulkan", nil, nil)
	if err != nil {
		return err
	}

	a.window = win

	win.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		a.frameBufferResized = true
	})

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
	vk.WaitForFences(a.logicalDevice, 1, []vk.Fence{a.inFlightFences[a.currentFrame]}, vk.True, vk.MaxUint64)

	res := vk.AcquireNextImage(a.logicalDevice, a.swapChain, vk.MaxUint64, a.imageAvailableSemaphores[a.currentFrame], vk.NullFence, &imageIndex)
	if res == vk.ErrorOutOfDate {
		return a.recreateSwapChain()
	} else if res != vk.Success && res != vk.Suboptimal {
		return fmt.Errorf("failed to acquire swapchain image")
	}

	if a.imagesInFlight[imageIndex] != vk.NullFence {
		vk.WaitForFences(a.logicalDevice, 1, []vk.Fence{a.imagesInFlight[imageIndex]}, vk.True, vk.MaxUint64)
	}

	a.imagesInFlight[imageIndex] = a.inFlightFences[a.currentFrame]

	waitsemaphores := []vk.Semaphore{a.imageAvailableSemaphores[a.currentFrame]}
	signalsemaphores := []vk.Semaphore{a.renderFinishedSemaphores[a.currentFrame]}
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

	vk.ResetFences(a.logicalDevice, 1, []vk.Fence{a.inFlightFences[a.currentFrame]})
	err := vk.Error(vk.QueueSubmit(a.graphicsQueue, 1, submitInfo, a.inFlightFences[a.currentFrame]))
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

	res = vk.QueuePresent(a.presentQueue, &presentInfo)
	if res == vk.ErrorOutOfDate || res == vk.Suboptimal || a.frameBufferResized {
		a.frameBufferResized = false
		return a.recreateSwapChain()
	} else if res != vk.Success {
		return fmt.Errorf("failed to present swapchain image")
	}

	//vk.QueueWaitIdle(a.presentQueue)

	a.currentFrame = (a.currentFrame + 1) % maxFramesInFlight

	return nil
}

func (a *app) cleanup() {
	a.cleanupSwapChain()

	for i := 0; i < maxFramesInFlight; i++ {
		vk.DestroySemaphore(a.logicalDevice, a.renderFinishedSemaphores[i], nil)
		vk.DestroySemaphore(a.logicalDevice, a.imageAvailableSemaphores[i], nil)
		vk.DestroyFence(a.logicalDevice, a.inFlightFences[i], nil)
	}
	vk.DestroyCommandPool(a.logicalDevice, a.commandPool, nil)

	vk.DestroyDevice(a.logicalDevice, nil)
	if a.config.EnableValidationLayers {
		vk.DestroyDebugReportCallback(a.instance, a.debugMessenger, nil)
	}
	vk.DestroySurface(a.instance, a.windowSurface, nil)
	vk.DestroyInstance(a.instance, nil)
	a.window.Destroy()
	glfw.Terminate()
}

func (a *app) cleanupSwapChain() {
	for _, v := range a.swapChainFrameBuffers {
		vk.DestroyFramebuffer(a.logicalDevice, v, nil)
	}

	vk.FreeCommandBuffers(a.logicalDevice, a.commandPool, uint32(len(a.commandBuffers)), a.commandBuffers)

	vk.DestroyPipeline(a.logicalDevice, a.graphicsPipeline, nil)
	vk.DestroyPipelineLayout(a.logicalDevice, a.pipelineLayout, nil)
	vk.DestroyRenderPass(a.logicalDevice, a.renderPass, nil)
	for i := range a.swapChainImageViews {
		vk.DestroyImageView(a.logicalDevice, a.swapChainImageViews[i], nil)
	}
	vk.DestroySwapchain(a.logicalDevice, a.swapChain, nil)
}

func (a *app) recreateSwapChain() error {
	w, h := a.window.GetFramebufferSize()
	for w == 0 || h == 0 {
		w, h = a.window.GetFramebufferSize()
		glfw.WaitEvents()
	}

	vk.DeviceWaitIdle(a.logicalDevice)
	a.cleanupSwapChain()

	err := a.createSwapChain()
	if err != nil {
		return err
	}

	err = a.createImageViews()
	if err != nil {
		return err
	}

	err = a.createRenderPass()
	if err != nil {
		return err
	}

	err = a.createGraphicsPipeline()
	if err != nil {
		return err
	}

	err = a.createFrameBuffers()
	if err != nil {
		return err
	}

	err = a.createCommandBuffers()
	if err != nil {
		return err
	}

	return nil
}
