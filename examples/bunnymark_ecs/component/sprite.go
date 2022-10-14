package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type SpriteData struct {
	Image *ebiten.Image
}

var Sprite = donburi.NewComponentType[SpriteData]()

func GetSprite(entry *donburi.Entry) *SpriteData {
	return donburi.Get[SpriteData](entry, Sprite)
}
