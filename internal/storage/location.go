package storage

// Location is a location of an entity in the storage.
type Location struct {
	Archetype ArchetypeIndex
	Component ComponentIndex
	Valid     bool
}

// NewLocation creates a new EntityLocation.
func NewLocation(archetype ArchetypeIndex, component ComponentIndex) *Location {
	return &Location{
		Archetype: archetype,
		Component: component,
		Valid:     true,
	}
}
