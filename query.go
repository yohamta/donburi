package donburi

import (
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/internal/storage"
)

type cache struct {
	archetypes []storage.ArchetypeIndex
	seen       int
}

// Query represents a query for entities.
// It is used to filter entities based on their components.
// It receives arbitrary filters that are used to filter entities.
// It contains a cache that is used to avoid re-evaluating the query.
// So it is not recommended to create a new query every time you want
// to filter entities with the same query.
type Query struct {
	layout_matches map[WorldId]*cache
	filter         filter.LayoutFilter
}

// NewQuery creates a new query.
// It receives arbitrary filters that are used to filter entities.
func NewQuery(filter filter.LayoutFilter) *Query {
	return &Query{
		layout_matches: make(map[WorldId]*cache),
		filter:         filter,
	}
}

// Each iterates over all entities that match the query.
func (q *Query) Each(w World, callback func(*Entry)) {
	accessor := w.StorageAccessor()
	result := q.evaluateQuery(w, &accessor)
	iter := storage.NewEntityIterator(0, accessor.Archetypes, result)
	f := func(entity storage.Entity) {
		entry := w.Entry(entity)
		callback(entry)
	}
	for iter.HasNext() {
		entities := iter.Next()
		for _, entity := range entities {
			f(entity)
		}
	}
}

// deprecated: use Each instead
func (q *Query) EachEntity(w World, callback func(*Entry)) {
	q.Each(w, callback)
}

// Count returns the number of entities that match the query.
func (q *Query) Count(w World) int {
	accessor := w.StorageAccessor()
	result := q.evaluateQuery(w, &accessor)
	iter := storage.NewEntityIterator(0, accessor.Archetypes, result)
	ret := 0
	for iter.HasNext() {
		entities := iter.Next()
		ret += len(entities)
	}
	return ret
}

// First returns the first entity that matches the query.
func (q *Query) First(w World) (entry *Entry, ok bool) {
	accessor := w.StorageAccessor()
	result := q.evaluateQuery(w, &accessor)
	iter := storage.NewEntityIterator(0, accessor.Archetypes, result)
	if !iter.HasNext() {
		return nil, false
	}
	for iter.HasNext() {
		entities := iter.Next()
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
	if _, ok := q.layout_matches[w]; !ok {
		q.layout_matches[w] = &cache{
			archetypes: make([]storage.ArchetypeIndex, 0),
			seen:       0,
		}
	}
	cache := q.layout_matches[w]
	for it := accessor.Index.SearchFrom(q.filter, cache.seen); it.HasNext(); {
		cache.archetypes = append(cache.archetypes, it.Next())
	}
	cache.seen = len(accessor.Archetypes)
	return cache.archetypes
}
