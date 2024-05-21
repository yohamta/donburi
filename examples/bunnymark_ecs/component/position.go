package component

import (
	"github.com/yohamta/donburi"
)

type PositionData struct {
	X, Y float64
}

func (p PositionData) Order() int {
	return int(p.Y * 600)
}

var Position = donburi.NewComponentType[PositionData]()
