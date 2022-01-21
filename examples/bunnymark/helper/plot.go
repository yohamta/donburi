package helper

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Plot struct {
	Values []float64
	Max    float64
	Bars   int
}

func NewPlot(bars int, max float64) *Plot {
	p := &Plot{
		Values: make([]float64, bars),
		Bars:   bars,
		Max:    max,
	}
	return p
}

func (p *Plot) Update(value float64) {
	p.Values = append(p.Values, value)
	if len(p.Values) > p.Bars {
		p.Values = p.Values[1:]
	}
}

func (p *Plot) Draw(screen *ebiten.Image, x, y, w, h float64) {
	ebitenutil.DrawRect(screen, x, y, w, h, color.RGBA{A: 128})
	if len(p.Values) < 2 {
		return
	}
	barW := w / float64(p.Bars)
	for i := 0; i < p.Bars; i++ {
		c := color.RGBA{R: 118, G: 222, B: 211, A: 255}
		if i%2 == 0 {
			c = color.RGBA{R: 106, G: 196, B: 186, A: 255}
		}
		height := 0.0
		if i < len(p.Values) {
			relH := p.Values[i] / p.Max
			if relH > 1 {
				relH = 1
			}
			height = relH * h
		}
		ebitenutil.DrawRect(screen, x+float64(i)*barW, y+h-height+height/2, barW, height/2, c)
	}
}

func (p *Plot) Last() float64 {
	if len(p.Values) < 1 {
		return 0
	}
	return p.Values[len(p.Values)-1]
}
