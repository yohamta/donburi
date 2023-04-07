package storage

// LocationMap is a storage of entity locations.
type LocationMap struct {
	LocationMap []*Location
	Len         int
}

// NewLocationMap creates an empty storage.
func NewLocationMap() *LocationMap {
	return &LocationMap{
		LocationMap: make([]*Location, 1, 256),
		Len:         0,
	}
}

// Contains returns true if the storage contains the given entity id.
func (lm *LocationMap) Contains(id EntityId) bool {
	val := lm.LocationMap[id]
	return val != nil && val.Valid
}

// Remove removes the given entity id from the storage.
func (lm *LocationMap) Remove(id EntityId) {
	lm.LocationMap[id].Valid = false
	lm.Len--
}

// Insert inserts the given entity id and archetype index to the storage.
func (lm *LocationMap) Insert(id EntityId, archetype ArchetypeIndex, component ComponentIndex) {
	if int(id) == len(lm.LocationMap) {
		loc := NewLocation(archetype, component)
		lm.LocationMap = append(lm.LocationMap, loc)
		lm.Len++
	} else {
		loc := lm.LocationMap[id]
		loc.Archetype = archetype
		loc.Component = component
		if !loc.Valid {
			lm.Len++
			loc.Valid = true
		}
	}
}

// Set sets the given entity id and archetype index to the storage.
func (lm *LocationMap) Set(id EntityId, loc *Location) {
	lm.Insert(id, loc.Archetype, loc.Component)
}

// Location returns the location of the given entity id.
func (lm *LocationMap) Location(id EntityId) *Location {
	return lm.LocationMap[id]
}

// Archetype returns the archetype of the given entity id.
func (lm *LocationMap) Archetype(id EntityId) ArchetypeIndex {
	return lm.LocationMap[id].Archetype
}

// Component returns the component of the given entity id.
func (lm *LocationMap) Component(id EntityId) ComponentIndex {
	return lm.LocationMap[id].Component
}
