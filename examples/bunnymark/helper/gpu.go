package helper

import (
	"github.com/jaypipes/ghw"
)

func GpuInfo() (gpu string) {
	gpu = "Unknown GPU"
	info, err := ghw.GPU()
	if err != nil {
		return
	}

	for _, gc := range info.GraphicsCards {
		if gc.DeviceInfo != nil {
			return gc.DeviceInfo.Product.Name
		}
	}
	return
}
