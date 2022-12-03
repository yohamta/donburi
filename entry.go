package donburi

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/internal/entity"
	"github.com/yohamta/donburi/internal/storage"
)

// Entry is a struct that contains an entity and a location in an archetype.
type Entry struct {
	World *world

	id     entity.EntityId
	entity Entity
	loc    *storage.Location
}

// Get returns the component from the entry
func Get[T any](e *Entry, ctype component.IComponentType) *T {
	return (*T)(e.Component(ctype))
}

// Add adds the component to the entry.
func Add[T any](e *Entry, ctype component.IComponentType, component *T) {
	e.AddComponent(ctype, unsafe.Pointer(component))
}

// Set sets the comopnent of the entry.
func Set[T any](e *Entry, ctype component.IComponentType, component *T) {
	e.SetComponent(ctype, unsafe.Pointer(component))
}

// SetValue sets the value of the component.
func SetValue[T any](e *Entry, ctype component.IComponentType, value T) {
	c := Get[T](e, ctype)
	*c = value
}

// Remove removes the component from the entry.
func Remove[T any](e *Entry, ctype component.IComponentType) {
	e.RemoveComponent(ctype)
}

// Valid returns true if the entry is valid.
func Valid(e *Entry) bool {
	if e == nil {
		return false
	}
	return e.Valid()
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
func (e *Entry) Component(ctype component.IComponentType) unsafe.Pointer {
	c := e.loc.Component
	a := e.loc.Archetype
	return e.World.components.Storage(ctype).Component(a, c)
}

// SetComponent sets the component.
func (e *Entry) SetComponent(ctype component.IComponentType, component unsafe.Pointer) {
	c := e.loc.Component
	a := e.loc.Archetype
	e.World.components.Storage(ctype).SetComponent(a, c, component)
}

// AddComponent adds the component to the entity.
func (e *Entry) AddComponent(ctype component.IComponentType, components ...unsafe.Pointer) {
	if len(components) > 1 {
		panic("AddComponent: component argument must be a single value")
	}
	if !e.HasComponent(ctype) {
		c := e.loc.Component
		a := e.loc.Archetype

		base_layout := e.World.archetypes[a].Layout().Components()
		target_arc := e.World.getArchetypeForComponents(append(base_layout, ctype))
		e.World.TransferArchetype(a, target_arc, c)

		e.loc = e.World.Entry(e.entity).loc
	}
	if len(components) == 1 {
		e.SetComponent(ctype, components[0])
	}
}

// RemoveComponent removes the component from the entity.
func (e *Entry) RemoveComponent(ctype component.IComponentType) {
	if !e.Archetype().Layout().HasComponent(ctype) {
		return
	}

	c := e.loc.Component
	a := e.loc.Archetype

	base_layout := e.World.archetypes[a].Layout().Components()
	target_layout := make([]component.IComponentType, 0, len(base_layout)-1)
	for _, c2 := range base_layout {
		if c2 == ctype {
			continue
		}
		target_layout = append(target_layout, c2)
	}

	target_arc := e.World.getArchetypeForComponents(target_layout)
	e.World.TransferArchetype(e.loc.Archetype, target_arc, c)

	e.loc = e.World.Entry(e.entity).loc
}

// Remove removes the entity from the world.
func (e *Entry) Remove() {
	e.World.Remove(e.entity)
}

// Valid returns true if the entry is valid.
func (e *Entry) Valid() bool {
	return e.World.Valid(e.entity)
}

// Archetype returns the archetype.
func (e *Entry) Archetype() *storage.Archetype {
	a := e.loc.Archetype
	return e.World.archetypes[a]
}

// HasComponent returns true if the entity has the given component type.
func (e *Entry) HasComponent(componentType component.IComponentType) bool {
	return e.Archetype().Layout().HasComponent(componentType)
}

func (e *Entry) String() string {
	var out bytes.Buffer
	out.WriteString("Entry: {")
	out.WriteString(e.Entity().String())
	out.WriteString(", ")
	out.WriteString(e.Archetype().Layout().String())
	out.WriteString(", Valid: ")
	out.WriteString(fmt.Sprintf("%v", e.Valid()))
	out.WriteString("}")
	return out.String()
}
