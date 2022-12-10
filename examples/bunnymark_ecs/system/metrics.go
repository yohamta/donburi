package system

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type Metrics struct {
	bounds   *image.Rectangle
	settings *component.SettingsData
}

func NewMetrics(bounds *image.Rectangle) *Metrics {
	return &Metrics{
		bounds: bounds,
	}
}

func (m *Metrics) Update(ecs *ecs.ECS) {
	if m.settings == nil {
		if entry, ok := component.Settings.First(ecs.World); ok {
			m.settings = component.Settings.Get(entry)
		}
	}
	select {
	case <-m.settings.Ticker.C:
		m.settings.Objects.Update(float64(ecs.World.Len() - 1))
		m.settings.Tps.Update(ebiten.CurrentTPS())
		m.settings.Fps.Update(ebiten.CurrentFPS())
	default:
	}
}

func (m *Metrics) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	str := fmt.Sprintf(
		"GPU: %s\nTPS: %.2f, FPS: %.2f, Objects: %.f\nBatching: %t, Amount: %d\nResolution: %dx%d",
		m.settings.Gpu, m.settings.Tps.Last(), m.settings.Fps.Last(), m.settings.Objects.Last(),
		!m.settings.Colorful, m.settings.Amount,
		m.bounds.Dx(), m.bounds.Dy(),
	)

	rect := text.BoundString(basicfont.Face7x13, str)
	width, height := float64(rect.Dx()), float64(rect.Dy())

	padding := 20.0
	rectW, rectH := width+padding, height+padding
	plotW, plotH := 100.0, 40.0

	ebitenutil.DrawRect(screen, 0, 0, rectW, rectH, color.RGBA{A: 128})
	text.Draw(screen, str, basicfont.Face7x13, int(padding)/2, 10+int(padding)/2, colornames.White)

	m.settings.Tps.Draw(screen, 0, padding+rectH, plotW, plotH)
	m.settings.Fps.Draw(screen, 0, padding+rectH*2, plotW, plotH)
	m.settings.Objects.Draw(screen, 0, padding+rectH*3, plotW, plotH)
}
