package ecs

import (
	"fmt"
	"reflect"

	"github.com/yohamta/donburi"
)

type Layer struct {
	*layer
	renderers map[string][]any
}

func newLayer(l *layer) *Layer {
	return &Layer{l, make(map[string][]any)}
}

func keyForType(typ reflect.Type) string {
	return fmt.Sprintf("%s/%s", typ.PkgPath(), typ.Name())
}

func invoke(fn any, e *ECS, arg any) {
	v := reflect.ValueOf(fn)
	v.Call([]reflect.Value{reflect.ValueOf(e), reflect.ValueOf(arg)})
}

func (l *Layer) draw(e *ECS, arg any) {
	key := keyForType(reflect.TypeOf(arg))
	for _, fn := range l.renderers[key] {
		invoke(fn, e, arg)
	}
}

func (l *Layer) addRenderer(r any) {
	// check renderer type is func(*ECS, any)
	typ := reflect.TypeOf(r)
	if typ.Kind() != reflect.Func {
		panic("renderer must be a function")
	}
	if typ.NumIn() != 2 {
		panic("renderer must have 2 arguments")
	}
	if typ.In(0) != reflect.TypeOf(&ECS{}) {
		panic("first argument must be *ECS")
	}
	if typ.NumOut() != 0 {
		panic("renderer must not have return values")
	}
	// add renderer
	key := keyForType(typ.In(1))
	l.renderers[key] = append(l.renderers[key], r)
}

var (
	layers []*layer
)

type layer struct {
	id  LayerID
	tag donburi.IComponentType
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
