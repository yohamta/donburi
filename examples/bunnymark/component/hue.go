package component

import (
	"github.com/yohamta/donburi"
)

type HueData struct {
	Colorful *bool
	Value    float64
}

var Hue = donburi.NewComponentType[HueData]()
