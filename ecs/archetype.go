package ecs

import (
	"github.com/yohamta/donburi"
)

type Archetype struct {
	components []donburi.IComponentType
}

func NewArchetype(cs ...donburi.IComponentType) *Archetype {
	return &Archetype{
		components: cs,
	}
}

func (a *Archetype) Spawn(ecs *ECS, cs ...donburi.IComponentType) *donburi.Entry {
	return a.spawn(ecs, LayerDefault, cs...)
}

func (a *Archetype) SpawnOnLayer(ecs *ECS, layer LayerID, cs ...donburi.IComponentType) *donburi.Entry {
	return a.spawn(ecs, layer, cs...)
}

func (a *Archetype) spawn(ecs *ECS, layer LayerID, cs ...donburi.IComponentType) *donburi.Entry {
	e := ecs.World.Entry(ecs.Create(
		layer,
		append(a.components, cs...)...,
	))
	return e
}
