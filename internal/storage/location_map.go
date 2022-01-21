package storage

import (
	"github.com/yohamta/donburi/internal/entity"
)

// LocationMap is a storage of entity locations.
type LocationMap struct {
	LocationMap map[entity.EntityId]*EntityLocation
	Len         int
}

// NewLocationMap creates an empty storage.
func NewLocationMap() *LocationMap {
	return &LocationMap{
		LocationMap: make(map[entity.EntityId]*EntityLocation),
		Len:         0,
	}
}

// Contains returns true if the storage contains the given entity id.
func (lm *LocationMap) Contains(id entity.EntityId) bool {
	val := lm.LocationMap[id]
	return val != nil && val.Valid
}

// Remove removes the given entity id from the storage.
func (lm *LocationMap) Remove(id entity.EntityId) {
	lm.LocationMap[id].Valid = false
	lm.Len--
}

// Insert inserts the given entity id and archetype index to the storage.
func (lm *LocationMap) Insert(id entity.EntityId, archetype ArchetypeIndex, component ComponentIndex) {
	if loc, ok := lm.LocationMap[id]; ok {
		loc.Archetype = archetype
		loc.Component = component
		if !loc.Valid {
			lm.Len++
			loc.Valid = true
		}
	} else {
		lm.LocationMap[id] = NewEntityLocation(archetype, component)
		lm.Len++
	}
}

// Set sets the given entity id and archetype index to the storage.
func (lm *LocationMap) Set(id entity.EntityId, loc *EntityLocation) {
	lm.Insert(id, loc.Archetype, loc.Component)
}

// Location returns the location of the given entity id.
func (lm *LocationMap) Location(id entity.EntityId) *EntityLocation {
	return lm.LocationMap[id]
}

// Archetype returns the archetype of the given entity id.
func (lm *LocationMap) Archetype(id entity.EntityId) ArchetypeIndex {
	return lm.LocationMap[id].Archetype
}

// Component returns the component of the given entity id.
func (lm *LocationMap) Component(id entity.EntityId) ComponentIndex {
	return lm.LocationMap[id].Component
}
