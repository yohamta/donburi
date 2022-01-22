package component

import (
	"github.com/yohamta/donburi"
)

type PositionData struct {
	X, Y float64
}

var Position = donburi.NewComponentType(PositionData{})

func GetPositionData(entry *donburi.Entry) *PositionData {
	return (*PositionData)(entry.Component(Position))
}
