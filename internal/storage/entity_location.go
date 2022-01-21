package storage

// EntityLocation is a location of an entity in the storage.
type EntityLocation struct {
	Archetype ArchetypeIndex
	Component ComponentIndex
	Valid     bool
}

// NewEntityLocation creates a new EntityLocation.
func NewEntityLocation(archetype ArchetypeIndex, component ComponentIndex) *EntityLocation {
	return &EntityLocation{
		Archetype: archetype,
		Component: component,
		Valid:     true,
	}
}
