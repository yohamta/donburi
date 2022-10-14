package component

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/helper"
)

type SettingsData struct {
	Ticker   *time.Ticker
	Sprite   *ebiten.Image
	Colorful bool
	Amount   int
	Gpu      string
	Tps      *helper.Plot
	Fps      *helper.Plot
	Objects  *helper.Plot
}

var Settings = donburi.NewComponentType[SettingsData]()

func GetSettings(entry *donburi.Entry) *SettingsData {
	return donburi.Get[SettingsData](entry, Settings)
}
