package component

import (
	"github.com/yohamta/donburi"
)

type VelocityData struct {
	X, Y float64
}

var Velocity = donburi.NewComponentType(VelocityData{})
