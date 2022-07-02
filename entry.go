package donburi

import (
	"unsafe"

	"github.com/yohamta/donburi/internal/component"
	"github.com/yohamta/donburi/internal/entity"
	"github.com/yohamta/donburi/internal/storage"
)

// Entry is a struct that contains an entity and a location in an archetype.
type Entry struct {
	id     entity.EntityId
	entity Entity
	loc    *storage.EntityLocation
	world  *world
}

// Get returns the component from the entry
func Get[T any](e *Entry, ctype *component.ComponentType) *T {
	return (*T)(e.Component(ctype))
}

// Set returns the component from the entry
func Set[T any](e *Entry, ctype *component.ComponentType, component *T) {
	e.SetComponent(ctype, unsafe.Pointer(component))
}

// Add adds the component to the entry.
func Add[T any](e *Entry, ctype *component.ComponentType, component *T) {
	e.AddComponent(ctype, unsafe.Pointer(component))
}

// Remove removes the component from the entry.
func Remove(e *Entry, ctype *component.ComponentType) {
	e.RemoveComponent(ctype)
}

// Id returns the entity id.
func (e *Entry) Id() entity.EntityId {
	return e.id
}

// Entity returns the entity.
func (e *Entry) Entity() Entity {
	return e.entity
}

// Component returns the component.
func (e *Entry) Component(ctype *component.ComponentType) unsafe.Pointer {
	c := e.loc.Component
	a := e.loc.Archetype
	return e.world.components.Storage(ctype).Component(a, c)
}

// SetComponent sets the component.
func (e *Entry) SetComponent(ctype *component.ComponentType, component unsafe.Pointer) {
	c := e.loc.Component
	a := e.loc.Archetype
	e.world.components.Storage(ctype).SetComponent(a, c, component)
}

// AddComponent adds the component to the entity.
func (e *Entry) AddComponent(ctype *component.ComponentType, components ...unsafe.Pointer) {
	if len(components) > 1 {
		panic("AddComponent: component argument must be a single value")
	}

	c := e.loc.Component
	a := e.loc.Archetype

	base_layout := e.world.archetypes[a].Layout().Components()
	target_arc := e.world.getArchetypeForComponents(append(base_layout, ctype))
	e.world.TransferArchetype(a, target_arc, c)

	e.loc = e.world.Entry(e.entity).loc
	if len(components) == 1 {
		e.SetComponent(ctype, components[0])
	}
}

// RemoveComponent removes the component from the entity.
func (e *Entry) RemoveComponent(ctype *component.ComponentType) {
	if !e.Archetype().Layout().HasComponent(ctype) {
		return
	}

	c := e.loc.Component
	a := e.loc.Archetype

	base_layout := e.world.archetypes[a].Layout().Components()
	target_layout := make([]*component.ComponentType, 0, len(base_layout)-1)
	for _, c2 := range base_layout {
		if c2 == ctype {
			continue
		}
		target_layout = append(target_layout, c2)
	}

	target_arc := e.world.getArchetypeForComponents(target_layout)
	e.world.TransferArchetype(e.loc.Archetype, target_arc, c)

	e.loc = e.world.Entry(e.entity).loc
}

// Archetype returns the archetype.
func (e *Entry) Archetype() *storage.Archetype {
	a := e.loc.Archetype
	return e.world.archetypes[a]
}

// HasComponent returns true if the entity has the given component type.
func (e *Entry) HasComponent(componentType *component.ComponentType) bool {
	return e.Archetype().Layout().HasComponent(componentType)
}
