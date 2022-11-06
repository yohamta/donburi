package ecs

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Layer struct {
	*layer
	renderers []Renderer
}

func newLayer(l *layer) *Layer {
	return &Layer{l, []Renderer{}}
}

func (l *Layer) draw(e *ECS, i *ebiten.Image) {
	screen := i
	for _, s := range l.renderers {
		s(e, screen)
	}
}

func (l *Layer) addRenderer(r Renderer) {
	l.renderers = append(l.renderers, r)
}

var (
	layers []*layer
)

type layer struct {
	id  LayerID
	tag *donburi.ComponentType
}

func getLayer(layerID LayerID) *layer {
	if int(layerID) >= len(layers) {
		layers = append(layers, make([]*layer, int(layerID)-len(layers)+1)...)
	}
	if layers[layerID] == nil {
		layers[layerID] = &layer{
			id:  layerID,
			tag: donburi.NewTag().SetName(fmt.Sprintf("Layer%d", layerID)),
		}
	}
	return layers[layerID]
}
