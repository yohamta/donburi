package component

import (
	"github.com/yohamta/donburi"
)

type VelocityData struct {
	X, Y float64
}

var Velocity = donburi.NewComponentType[VelocityData]()

func GetVelocity(entry *donburi.Entry) *VelocityData {
	return donburi.Get[VelocityData](entry, Velocity)
}
