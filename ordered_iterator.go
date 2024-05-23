package donburi

import "sort"

// OrderedEntryIterator is an iterator for entries from a list of `[]Entity`.
type OrderedEntryIterator[T IOrderable] struct {
	current int
	entries []*Entry
}

// OrderedEntryIterator is an iterator for entries based on a list of `[]Entity`.
func NewOrderedEntryIterator[T IOrderable](current int, w World, entities []Entity, orderedBy *ComponentType[T]) OrderedEntryIterator[T] {
	entLen := len(entities)
	entries := make([]*Entry, entLen)
	orders := make([]int, entLen)

	for i := 0; i < entLen; i++ {
		entry := w.Entry(entities[i])
		entries[i] = entry
		orders[i] = (*orderedBy.Get(entry)).Order()
	}

	sort.Slice(entries, func(i, j int) bool {
		return orders[i] < orders[j]
	})

	return OrderedEntryIterator[T]{
		entries: entries,
		current: current,
	}
}

// HasNext returns true if there are more entries to iterate over.
func (it *OrderedEntryIterator[T]) HasNext() bool {
	return it.current < len(it.entries)
}

// Next returns the next entry.
func (it *OrderedEntryIterator[T]) Next() *Entry {
	nextIndex := it.entries[it.current]
	it.current++
	return nextIndex
}
