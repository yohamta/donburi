package component

import (
	"github.com/yohamta/donburi"
)

type GravityData struct {
	Value float64
}

var Gravity = donburi.NewComponentType(GravityData{})

func GetGravity(entry *donburi.Entry) *GravityData {
	return (*GravityData)(entry.Component(Gravity))
}
