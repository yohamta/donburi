package component

import (
	"github.com/yohamta/donburi"
)

type PositionData struct {
	X, Y float64
	ID   int
}

func (p PositionData) Order() int {
	return -p.ID
}

var Position = donburi.NewComponentType[PositionData]()
