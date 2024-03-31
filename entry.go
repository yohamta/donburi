package donburi

import (
	"bytes"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/internal/storage"
)

// Entry is a struct that contains an entity and a location in an archetype.
type Entry struct {
	World *world

	id     storage.EntityId
	entity Entity
	loc    *storage.Location
}

// Get returns the component from the entry
func Get[T any](e *Entry, c component.IComponentType) *T {
	return (*T)(e.Component(c))
}

// GetComponents uses reflection to convert the unsafe.Pointers of an entry into its component data instances.
// Note that this is likely to be slow and should not be used in hot paths or if not necessary.
func GetComponents(e *Entry) []any {
	var instances []any
	for _, c := range e.World.StorageAccessor().Archetypes[e.loc.Archetype].ComponentTypes() {
		instances = append(instances, reflect.Indirect(
			reflect.NewAt(c.Typ(), e.Component(c))).Interface(),
		)
	}
	return instances
}

// Add adds the component to the entry.
func Add[T any](e *Entry, c component.IComponentType, component *T) {
	e.AddComponent(c, unsafe.Pointer(component))
}

// Set sets the comopnent of the entry.
func Set[T any](e *Entry, c component.IComponentType, component *T) {
	e.SetComponent(c, unsafe.Pointer(component))
}

// SetValue sets the value of the component.
func SetValue[T any](e *Entry, c component.IComponentType, value T) {
	*Get[T](e, c) = value
}

// GetValue gets the value of the component.
func GetValue[T any](e *Entry, c component.IComponentType) T {
	return *Get[T](e, c)
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
func (e *Entry) Id() storage.EntityId {
	return e.id
}

// Entity returns the entity.
func (e *Entry) Entity() Entity {
	return e.entity
}

// Component returns the component.
func (e *Entry) Component(c component.IComponentType) unsafe.Pointer {
	return e.World.components.Storage(c).Component(e.loc.Archetype, e.loc.Component)
}

// SetComponent sets the component.
func (e *Entry) SetComponent(c component.IComponentType, component unsafe.Pointer) {
	e.World.components.Storage(c).SetComponent(e.loc.Archetype, e.loc.Component, component)
}

// AddComponent adds the component to the entity.
func (e *Entry) AddComponent(c component.IComponentType, components ...unsafe.Pointer) {
	if len(components) > 1 {
		panic("AddComponent: component argument must be a single value")
	}
	if !e.HasComponent(c) {
		archetypeIndex := e.loc.Archetype
		targetArchetype := e.World.getArchetypeForComponents(
			append(e.World.archetypes[archetypeIndex].Layout().Components(), c),
		)
		e.World.TransferArchetype(archetypeIndex, targetArchetype, e.loc.Component)
		e.loc = e.World.Entry(e.entity).loc
	}
	if len(components) == 1 {
		e.SetComponent(c, components[0])
	}
}

// RemoveComponent removes the component from the entity.
func (e *Entry) RemoveComponent(c component.IComponentType) {
	if !e.Archetype().Layout().HasComponent(c) {
		return
	}

	baseLayout := e.World.archetypes[e.loc.Archetype].Layout().Components()
	targetLayout := make([]component.IComponentType, 0, len(baseLayout)-1)
	for _, c2 := range baseLayout {
		if c2 == c {
			continue
		}
		targetLayout = append(targetLayout, c2)
	}

	targetArchetype := e.World.getArchetypeForComponents(targetLayout)
	e.World.TransferArchetype(e.loc.Archetype, targetArchetype, e.loc.Component)
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
func (e *Entry) HasComponent(c component.IComponentType) bool {
	return e.Archetype().Layout().HasComponent(c)
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
