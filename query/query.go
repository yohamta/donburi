package query

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// Query represents a query for entities.
// It is used to filter entities based on their components.
// It receives arbitrary filters that are used to filter entities.
// It contains a cache that is used to avoid re-evaluating the query.
// So it is not recommended to create a new query every time you want
// to filter entities with the same query.
// deprecated: use donburi.Query instead
type Query = donburi.Query

// NewQuery creates a new query.
// It receives arbitrary filters that are used to filter entities.
// deprecated: use donburi.NewQuery instead
func NewQuery(filter filter.LayoutFilter) *Query {
	return donburi.NewQuery(filter)
}
