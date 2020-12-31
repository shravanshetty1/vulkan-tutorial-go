package app

import (
	"fmt"
	"log"
	"strings"
	"unicode"
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

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

func checkDeviceExtensionsSupport(device vk.PhysicalDevice, requiredDeviceExtensions []string) bool {
	var count uint32
	vk.EnumerateDeviceExtensionProperties(device, "", &count, nil)
	extensionProperties := make([]vk.ExtensionProperties, count)
	vk.EnumerateDeviceExtensionProperties(device, "", &count, extensionProperties)

	supportedExtensions := make(map[string]bool, len(extensionProperties))
	for _, ep := range extensionProperties {
		ep.Deref()
		supportedExtensions[vk.ToString(ep.ExtensionName[:])] = true
		ep.Free()
	}

	for _, requiredExtension := range requiredDeviceExtensions {
		if !supportedExtensions[requiredExtension] {
			return false
		}
	}

	return true
}

type swapChainSupportDetails struct {
	capabilities      vk.SurfaceCapabilities
	surfaceFormats    []vk.SurfaceFormat
	presentationModes []vk.PresentMode
}

func querySwapChainSupport(device vk.PhysicalDevice, surface vk.Surface) swapChainSupportDetails {
	var details swapChainSupportDetails

	vk.GetPhysicalDeviceSurfaceCapabilities(device, surface, &details.capabilities)
	details.capabilities.Deref()
	//details.capabilities.Free()

	var formatCount uint32
	vk.GetPhysicalDeviceSurfaceFormats(device, surface, &formatCount, nil)
	if formatCount != 0 {
		surfaceFormats := make([]vk.SurfaceFormat, formatCount)
		vk.GetPhysicalDeviceSurfaceFormats(device, surface, &formatCount, surfaceFormats)
		details.surfaceFormats = surfaceFormats
	}

	var modeCount uint32
	vk.GetPhysicalDeviceSurfacePresentModes(device, surface, &modeCount, nil)
	if modeCount != 0 {
		presentationModes := make([]vk.PresentMode, modeCount)
		vk.GetPhysicalDeviceSurfacePresentModes(device, surface, &modeCount, presentationModes)
		details.presentationModes = presentationModes
	}

	return details
}
