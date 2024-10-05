// donburi is an Entity Component System library for Go/Ebitengine inspired by legion.
//
// It aims to be a feature rich and high-performance ECS Library.
package donburi

import (
	"fmt"
	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/internal/storage"
)

// WorldId is a unique identifier for a world.
type WorldId int

// World is a collection of entities and components.
type World interface {
	// Id returns the unique identifier for the world.
	Id() WorldId
	// Create creates a new entity with the specified components.
	Create(components ...component.IComponentType) Entity
	// CreateMany creates a new entity with the specified components.
	CreateMany(n int, components ...component.IComponentType) []Entity
	// Entry returns an entry for the specified entity.
	Entry(entity Entity) *Entry
	// Remove removes the specified entity.
	Remove(entity Entity)
	// Valid returns true if the specified entity is valid.
	Valid(e Entity) bool
	// Len returns the number of entities in the world.
	Len() int
	// StorageAccessor returns an accessor for the world's storage.
	// It is used to access components and archetypes by queries.
	StorageAccessor() StorageAccessor
	// ArcheTypes returns the archetypes in the world.
	Archetypes() []*storage.Archetype

	// OnCreate registers a callback function that gets triggered when an entity is created.
	OnCreate(callback func(world World, entity Entity))

	// OnRemove registers a callback function that gets triggered when an entity is removed.
	// Note that it is called before the entity is removed from the ECS.
	OnRemove(callback func(world World, entity Entity))
}

// StorageAccessor is an accessor for the world's storage.
type StorageAccessor struct {
	// Index is the search index for the world.
	Index *storage.Index
	// Components is the component storage for the world.
	Components *storage.Components
	// Archetypes is the archetype storage for the world.
	Archetypes []*storage.Archetype
}

type initializer func(w World)

type world struct {
	id           WorldId
	index        *storage.Index
	entities     *storage.LocationMap
	components   *storage.Components
	archetypes   []*storage.Archetype
	destroyed    []Entity
	entries      []*Entry
	nextEntityId storage.EntityId

	removeCallbacks []func(world World, entity Entity)
	createCallbacks []func(world World, entity Entity)
}

var nextWorldId WorldId = 0

var registeredInitializers []initializer

// RegisterInitializer registers an initializer for a world.
func RegisterInitializer(initializer initializer) {
	registeredInitializers = append(registeredInitializers, initializer)
}

// NewWorld creates a new world.
func NewWorld() World {
	worldId := nextWorldId
	nextWorldId++
	w := &world{
		id:           worldId,
		index:        storage.NewIndex(),
		entities:     storage.NewLocationMap(),
		components:   storage.NewComponents(),
		archetypes:   make([]*storage.Archetype, 0),
		destroyed:    make([]Entity, 0, 256),
		entries:      make([]*Entry, 1, 256),
		nextEntityId: 1,
	}
	for _, initializer := range registeredInitializers {
		initializer(w)
	}
	return w
}

func (w *world) Id() WorldId {
	return w.id
}

func (w *world) CreateMany(num int, components ...component.IComponentType) []Entity {
	archetypeIndex := w.getArchetypeForComponents(components)
	entities := make([]Entity, 0, num)
	for i := 0; i < num; i++ {
		entities = append(entities, w.createEntity(archetypeIndex))
	}
	return entities
}

func (w *world) Create(components ...component.IComponentType) Entity {
	return w.createEntity(w.getArchetypeForComponents(components))
}

func (w *world) createEntity(archetypeIndex storage.ArchetypeIndex) Entity {
	entity := w.nextEntity()
	archetype := w.archetypes[archetypeIndex]
	componentIndex := w.components.PushComponents(archetype.Layout().Components(), archetypeIndex)
	w.entities.Insert(entity.Id(), archetypeIndex, componentIndex)
	archetype.PushEntity(entity)

	archetype.Lock()
	defer archetype.Unlock()

	w.createEntry(entity)

	er := entity.Ready()
	w.archetypes[archetypeIndex].Entities()[componentIndex] = er
	w.entries[er.Id()].entity = er.Ready()
	for _, callback := range w.createCallbacks {
		callback(w, er)
	}

	return er
}

func (w *world) createEntry(e Entity) *Entry {
	id := e.Id()
	if int(id) >= len(w.entries) {
		w.entries = append(w.entries, &Entry{id: id, entity: e, loc: w.entities.Location(id), World: w})
		return w.entries[id]
	}
	w.entries[id].loc = w.entities.Location(id)
	return w.entries[id]
}

func (w *world) Valid(e Entity) bool {
	if e == Null {
		return false
	}
	if !w.entities.Contains(e.Id()) {
		return false
	}
	loc := w.entities.LocationMap[e.Id()]
	// If the version of the entity is not the same as the version of the archetype,
	// the entity is invalid (it means the entity is already destroyed).
	return e.IsReady() && loc.Valid && e == w.archetypes[loc.Archetype].Entities()[loc.Component]
}

func (w *world) Entry(entity Entity) *Entry {
	id := entity.Id()
	entry := w.entries[id]
	entry.entity = entity
	entry.loc = w.entities.LocationMap[id]
	return entry
}

func (w *world) Len() int {
	return w.entities.Len
}

func (w *world) Remove(ent Entity) {
	if w.Valid(ent) {
		// Called before any operations so that user code can access all the data it might need
		for _, callback := range w.removeCallbacks {
			callback(w, ent)
		}
		w.entities.Remove(ent.Id())
		w.removeAtLocation(ent, w.entities.LocationMap[ent.Id()])
	}
}

func (w *world) removeAtLocation(ent Entity, loc *storage.Location) {
	componentIndex := loc.Component
	archetype := w.archetypes[loc.Archetype]
	archetype.SwapRemove(int(componentIndex))
	w.components.Remove(archetype, componentIndex)
	if int(componentIndex) < len(archetype.Entities()) {
		swapped := archetype.Entities()[componentIndex]
		w.entities.Set(swapped.Id(), loc)
	}
	w.destroyed = append(w.destroyed, ent.IncVersion())
}

func (w *world) TransferArchetype(from, to storage.ArchetypeIndex, idx storage.ComponentIndex) storage.ComponentIndex {
	if from == to {
		return idx
	}
	fromArchetype := w.archetypes[from]
	toArchetype := w.archetypes[to]

	// move entity id
	ent := fromArchetype.SwapRemove(int(idx))
	toArchetype.PushEntity(ent)
	w.entities.Insert(ent.Id(), to, storage.ComponentIndex(len(toArchetype.Entities())-1))

	if len(fromArchetype.Entities()) > int(idx) {
		moved := fromArchetype.Entities()[idx]
		w.entities.Insert(moved.Id(), from, idx)
	}

	// creates component if not exists in new layout
	fromLayout := fromArchetype.Layout()
	toLayout := toArchetype.Layout()
	for _, c := range toLayout.Components() {
		if !fromLayout.HasComponent(c) {
			st := w.components.Storage(c)
			st.PushComponent(c, to)
		}
	}

	// move components
	for _, c := range fromLayout.Components() {
		st := w.components.Storage(c)
		if toLayout.HasComponent(c) {
			st.MoveComponent(from, idx, to)
		} else {
			st.SwapRemove(from, idx)
		}
	}
	w.components.Move(from, to)

	return storage.ComponentIndex(len(toArchetype.Entities()) - 1)
}

func (w *world) StorageAccessor() StorageAccessor {
	return StorageAccessor{
		w.index,
		w.components,
		w.archetypes,
	}
}

func (w *world) OnCreate(callback func(world World, entity Entity)) {
	w.createCallbacks = append(w.createCallbacks, callback)
}

func (w *world) OnRemove(callback func(world World, entity Entity)) {
	w.removeCallbacks = append(w.removeCallbacks, callback)
}

func (w *world) Archetypes() []*storage.Archetype {
	return w.archetypes
}

func (w *world) nextEntity() Entity {
	if len(w.destroyed) == 0 {
		id := w.nextEntityId
		w.nextEntityId++
		return storage.NewEntity(id)
	}
	entity := w.destroyed[len(w.destroyed)-1]
	w.destroyed = w.destroyed[:len(w.destroyed)-1]
	return entity
}

func (w *world) insertArchetype(layout *storage.Layout) storage.ArchetypeIndex {
	w.index.Push(layout)
	archetypeIndex := storage.ArchetypeIndex(len(w.archetypes))
	w.archetypes = append(w.archetypes, storage.NewArchetype(archetypeIndex, layout))

	return archetypeIndex
}

func (w *world) getArchetypeForComponents(components []component.IComponentType) storage.ArchetypeIndex {
	if len(components) == 0 {
		panic("entity must have at least one component")
	}
	if i := w.index.Search(filter.Exact(components)); i.HasNext() {
		return i.Next()
	}
	if !w.noDuplicates(components) {
		panic(fmt.Sprintf("duplicate component types: %v", components))
	}
	return w.insertArchetype(storage.NewLayout(components))
}

func (w *world) noDuplicates(components []component.IComponentType) bool {
	// check if there are duplicate values inside components slice
	for i := 0; i < len(components); i++ {
		for j := i + 1; j < len(components); j++ {
			if components[i] == components[j] {
				return false
			}
		}
	}
	return true
}
