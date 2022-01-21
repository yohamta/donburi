package donburi

import (
	"fmt"
	"sort"

	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/internal/component"
	"github.com/yohamta/donburi/internal/entity"
	"github.com/yohamta/donburi/internal/storage"
)

// WorldId is a unique identifier for a world.
type WorldId int

// World is a collection of entities and components.
type World interface {
	// Id returns the unique identifier for the world.
	Id() WorldId
	// Create creates a new entity with the specified components.
	Create(components ...*component.ComponentType) Entity
	// CreateMany creates a new entity with the specified components.
	CreateMany(n int, components ...*component.ComponentType) []Entity
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
}

// StorageAccessor is an accessor for the world's storage.
type StorageAccessor struct {
	// Index is the search index for the world.
	Index *storage.SearchIndex
	// Components is the component storage for the world.
	Components *storage.Components
	// Archetypes is the archetype storage for the world.
	Archetypes []*storage.Archetype
}

type world struct {
	id           WorldId
	index        *storage.SearchIndex
	entities     *storage.LocationMap
	components   *storage.Components
	archetypes   []*storage.Archetype
	destroyed    []Entity
	entries      []*Entry
	nextEntityId entity.EntityId
}

var nextWorldId WorldId = 0

// NewWorld creates a new world.
func NewWorld() World {
	worldId := nextWorldId
	nextWorldId++
	w := &world{
		id:           worldId,
		index:        storage.NewSearchIndex(),
		entities:     storage.NewLocationMap(),
		components:   storage.NewComponents(),
		archetypes:   make([]*storage.Archetype, 0),
		destroyed:    make([]Entity, 0, 256),
		entries:      make([]*Entry, 1, 256),
		nextEntityId: 1,
	}
	return w
}

func (w *world) Id() WorldId {
	return w.id
}

func (w *world) CreateMany(num int, components ...*component.ComponentType) []Entity {
	archetypeIndex := w.getArchetypeForComponents(components)
	entities := make([]Entity, 0, num)
	for i := 0; i < num; i++ {
		entities = append(entities, w.createEntity(archetypeIndex))
	}
	return entities
}

func (w *world) Create(components ...*component.ComponentType) Entity {
	archetypeIndex := w.getArchetypeForComponents(components)
	return w.createEntity(archetypeIndex)
}

func (w *world) createEntity(archetypeIndex storage.ArchetypeIndex) Entity {
	entity := w.nextEntity()
	archetype := w.archetypes[archetypeIndex]
	componentIndex := w.components.PushComponents(archetype.Layout().Components(), archetypeIndex)
	w.entities.Insert(entity.Id(), archetypeIndex, componentIndex)
	archetype.PushEntity(entity)
	w.createEntry(entity)
	return entity
}

func (w *world) createEntry(e Entity) {
	id := e.Id()
	if int(id) >= len(w.entries) {
		w.entries = append(w.entries, nil)
	}
	w.entries[id] = &Entry{
		entity: e,
		loc:    w.entities.LocationMap[id],
		world:  w,
	}
}

func (w *world) Valid(e Entity) bool {
	if e == Null {
		return false
	}
	if !w.entities.Contains(e.Id()) {
		return false
	}
	loc := w.entities.LocationMap[e.Id()]
	a := loc.Archetype
	c := loc.Component
	// If the version of the entity is not the same as the version of the archetype,
	// the entity is invalid (it means the entity is already destroyed).
	return loc.Valid && e == w.archetypes[a].Entities()[c]
}

func (w *world) Entry(entity Entity) *Entry {
	entry := w.entries[entity.Id()]
	entry.entity = entity
	return entry
}

func (w *world) Len() int {
	return w.entities.Len
}

func (w *world) Remove(ent Entity) {
	if !w.Valid(ent) {
		return
	}
	loc := w.entities.LocationMap[ent.Id()]
	w.entities.Remove(ent.Id())
	w.removeAtLocation(ent, loc)
}

func (w *world) removeAtLocation(ent Entity, loc *storage.EntityLocation) {
	arch_index := loc.Archetype
	component_index := loc.Component
	archetype := w.archetypes[arch_index]
	archetype.SwapRemove(int(component_index))
	for _, component_type := range archetype.Layout().Components() {
		storage := w.components.Storage(component_type)
		storage.SwapRemove(arch_index, component_index)
	}
	if int(component_index) < len(archetype.Entities()) {
		swapped := archetype.Entities()[component_index]
		w.entities.Set(swapped.Id(), loc)
	}
	w.destroyed = append(w.destroyed, ent.IncVersion())
}

func (w *world) TransferArchetype(from, to storage.ArchetypeIndex, idx storage.ComponentIndex) storage.ComponentIndex {
	if from == to {
		return idx
	}
	from_arch := w.archetypes[from]
	to_arch := w.archetypes[to]

	// move entity id
	ent := from_arch.SwapRemove(int(idx))
	to_arch.PushEntity(ent)
	w.entities.Insert(ent.Id(), to, storage.ComponentIndex(len(to_arch.Entities())-1))

	if len(from_arch.Entities()) > int(idx) {
		moved := from_arch.Entities()[idx]
		w.entities.Insert(moved.Id(), from, storage.ComponentIndex(idx))
	}

	// creates component if not exists in new layout
	from_layout := from_arch.Layout()
	to_layout := to_arch.Layout()
	for _, component_type := range to_layout.Components() {
		if !from_layout.HasComponent(component_type) {
			storage := w.components.Storage(component_type)
			storage.PushComponent(component_type, to)
		}
	}

	// move components
	for _, component_type := range from_layout.Components() {
		storage := w.components.Storage(component_type)
		if to_layout.HasComponent(component_type) {
			storage.MoveComponent(from, idx, to)
		} else {
			storage.SwapRemove(from, idx)
		}
	}

	return storage.ComponentIndex(len(to_arch.Entities()) - 1)
}

func (w *world) StorageAccessor() StorageAccessor {
	return StorageAccessor{
		w.index,
		w.components,
		w.archetypes,
	}
}

func (w *world) nextEntity() Entity {
	if len(w.destroyed) == 0 {
		id := w.nextEntityId
		w.nextEntityId++
		return entity.NewEntity(id)
	}
	entity := w.destroyed[len(w.destroyed)-1]
	w.destroyed = w.destroyed[:len(w.destroyed)-1]
	return entity
}

func (w *world) insertArcheType(layout *storage.EntityLayout) storage.ArchetypeIndex {
	w.index.Push(layout)
	arch_index := storage.ArchetypeIndex(len(w.archetypes))
	w.archetypes = append(w.archetypes, storage.NewArchetype(arch_index, layout))

	return arch_index
}

func (w *world) getArchetypeForComponents(components []*component.ComponentType) storage.ArchetypeIndex {
	if len(components) == 0 {
		panic("entity must have at least one component")
	}
	sort.Slice(components, func(i, j int) bool {
		return components[i].Id() < components[j].Id()
	})
	if ii := w.index.Search(filter.Exact(components)); ii.HasNext() {
		return ii.Next()
	}
	if !w.noDuplicates(components) {
		panic(fmt.Sprintf("duplicate component types: %v", components))
	}
	return w.insertArcheType(storage.NewEntityLayout(components))
}

func (w *world) noDuplicates(components []*ComponentType) bool {
	// check if there're duplicate values inside components slice
	for i := 0; i < len(components); i++ {
		for j := i + 1; j < len(components); j++ {
			if components[i] == components[j] {
				return false
			}
		}
	}
	return true
}
