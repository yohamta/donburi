package donburi

import (
	"iter"
	"sort"

	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/internal/storage"
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

// OrderedQuery is a special extension of Query which has a type parameter used
// when running ordered queries using `EachOrdered`.
type OrderedQuery[T IOrderable] struct {
	Query
}

// NewOrderedQuery creates a new ordered query.
// It takes a filter parameter that is used when evaluating the query.
// Use `OrderedQuery.EachOrdered` to run a Each query in ordered mode.
func NewOrderedQuery[T IOrderable](filter filter.LayoutFilter) *OrderedQuery[T] {
	return &OrderedQuery[T]{
		//orderedBy: orderedBy,
		Query: Query{
			layoutMatches: make(map[WorldId]*cache),
			filter:        filter,
		},
	}
}

// IterOrdered returns an iterator over all entities within the query filter,
// ordered by the specified component.
func (q *OrderedQuery[T]) IterOrdered(w World, orderBy *ComponentType[T]) iter.Seq[*Entry] {
	return func(yield func(*Entry) bool) {
		accessor := w.StorageAccessor()
		iter := storage.NewEntityIterator(0, accessor.Archetypes, q.evaluateQuery(w, &accessor))

		var allEntries []*Entry

		for iter.HasNext() {
			archetype := iter.Next()
			archetype.Lock()

			ents := archetype.Entities()
			for _, entity := range ents {
				entry := w.Entry(entity)
				if entry.entity.IsReady() {
					allEntries = append(allEntries, entry)
				}
			}

			archetype.Unlock()
		}

		// Sort all entries at once
		sort.Slice(allEntries, func(i, j int) bool {
			return orderBy.GetValue(allEntries[i]).Order() < orderBy.GetValue(allEntries[j]).Order()
		})

		// Yield sorted entries
		for _, entry := range allEntries {
			if !yield(entry) {
				return
			}
		}
	}
}

// EachOrdered iterates over all entities within the query filter, and uses the `orderBy` parameter to
// figure out which property to order using.
// `T` must implement `IOrderable`
func (q *OrderedQuery[T]) EachOrdered(w World, orderBy *ComponentType[T], callback func(*Entry)) {
	for entry := range q.IterOrdered(w, orderBy) {
		callback(entry)
	}
}

func (q *Query) Iter(w World) iter.Seq[*Entry] {
	return func(yield func(*Entry) bool) {
		accessor := w.StorageAccessor()
		iter := storage.NewEntityIterator(0, accessor.Archetypes, q.evaluateQuery(w, &accessor))

		for iter.HasNext() {
			archetype := iter.Next()
			archetype.Lock()
			for _, entity := range archetype.Entities() {
				entry := w.Entry(entity)
				if entry.entity.IsReady() {
					if !yield(entry) {
						archetype.Unlock()
						return
					}
				}
			}
			archetype.Unlock()
		}
	}
}

// Each iterates over all entities that match the query.
func (q *Query) Each(w World, callback func(*Entry)) {
	for entry := range q.Iter(w) {
		callback(entry)
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
