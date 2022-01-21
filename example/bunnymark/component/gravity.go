package component

import (
	"github.com/yohamta/donburi"
)

type GravityData struct {
	Value float64
}

var Gravity = donburi.NewComponentType(GravityData{})
