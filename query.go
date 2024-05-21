package donburi

import (
	"cmp"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/internal/storage"
	"slices"
)

type cache struct {
	archetypes []storage.ArchetypeIndex
	seen       int
}

type IOrderable interface {
	Order() int
}

// Query represents a query for entities.
// It is used to filter entities based on their components.
// It receives arbitrary filters that are used to filter entities.
// It contains a cache that is used to avoid re-evaluating the query.
// So it is not recommended to create a new query every time you want
// to filter entities with the same query.
type Query struct {
	layoutMatches map[WorldId]*cache
	filter        filter.LayoutFilter
}

// NewQuery creates a new query.
// It receives arbitrary filters that are used to filter entities.
func NewQuery(filter filter.LayoutFilter) *Query {
	return &Query{
		layoutMatches: make(map[WorldId]*cache),
		filter:        filter,
	}
}

// Each iterates over all entities that match the query.
func (q *Query) Each(w World, callback func(*Entry)) {
	accessor := w.StorageAccessor()
	iter := storage.NewEntityIterator(0, accessor.Archetypes, q.evaluateQuery(w, &accessor))
	for iter.HasNext() {
		archetype := iter.Next()
		archetype.Lock()
		for _, entity := range archetype.Entities() {
			entry := w.Entry(entity)
			if entry.entity.IsReady() {
				callback(entry)
			}
		}
		archetype.Unlock()
	}
}

func (q *Query) EachOrdered(w World, callback func(*Entry), orderComponent IComponentType) {
	accessor := w.StorageAccessor()
	iter := storage.NewEntityIterator(0, accessor.Archetypes, q.evaluateQuery(w, &accessor))
	for iter.HasNext() {
		archetype := iter.Next()
		archetype.Lock()

		ents := archetype.Entities()
		slices.SortFunc(ents, func(firstEnt, secondEnt Entity) int {
			first := w.Entry(firstEnt)
			second := w.Entry(secondEnt)

			var firstOrder, secondOrder int

			// Convert first entry component
			firstPtr := first.Component(orderComponent)
			if firstOrderable, canOrder := convertToOrderable(firstPtr, orderComponent.Typ()); canOrder {
				firstOrder = firstOrderable.Order()
			}

			// Convert second entry component
			secondPtr := second.Component(orderComponent)
			if secondOrderable, canOrder := convertToOrderable(secondPtr, orderComponent.Typ()); canOrder {
				secondOrder = secondOrderable.Order()
			}

			return cmp.Compare(firstOrder, secondOrder)
		})

		for _, entity := range ents {
			entry := w.Entry(entity)
			if entry.entity.IsReady() {
				callback(entry)
			}
		}

		archetype.Unlock()
	}
}

// deprecated: use Each instead
func (q *Query) EachEntity(w World, callback func(*Entry)) {
	q.Each(w, callback)
}

// Count returns the number of entities that match the query.
func (q *Query) Count(w World) int {
	accessor := w.StorageAccessor()
	iter := storage.NewEntityIterator(0, accessor.Archetypes, q.evaluateQuery(w, &accessor))
	ret := 0
	for iter.HasNext() {
		archetype := iter.Next()
		ret += len(archetype.Entities())
	}
	return ret
}

// First returns the first entity that matches the query.
func (q *Query) First(w World) (entry *Entry, ok bool) {
	accessor := w.StorageAccessor()
	iter := storage.NewEntityIterator(0, accessor.Archetypes, q.evaluateQuery(w, &accessor))
	if !iter.HasNext() {
		return nil, false
	}
	for iter.HasNext() {
		archetype := iter.Next()
		entities := archetype.Entities()
		if len(entities) > 0 {
			return w.Entry(entities[0]), true
		}
	}
	return nil, false
}

// deprecated: use First instead
func (q *Query) FirstEntity(w World) (entry *Entry, ok bool) {
	return q.First(w)
}

func (q *Query) evaluateQuery(world World, accessor *StorageAccessor) []storage.ArchetypeIndex {
	w := world.Id()
	if _, ok := q.layoutMatches[w]; !ok {
		q.layoutMatches[w] = &cache{
			archetypes: make([]storage.ArchetypeIndex, 0),
			seen:       0,
		}
	}
	cache := q.layoutMatches[w]
	for it := accessor.Index.SearchFrom(q.filter, cache.seen); it.HasNext(); {
		cache.archetypes = append(cache.archetypes, it.Next())
	}
	cache.seen = len(accessor.Archetypes)
	return cache.archetypes
}
