package hierarchy

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type hierarchySystem struct {
	query *query.Query
}

// HierarchySystem is a system that removes children of invalid parents.
var HierarchySystem = &hierarchySystem{
	query: query.NewQuery(filter.Contains(parentComponent)),
}

// RemoveChildren removes children of invalid parents.
// This function is useful when you want to remove children of invalid parents.
// You don't need to use this system function when you use
// hierarchy.RemoveRecursive() or hierarchy.RemoveChildrenRecursive() instead.
func (hs *hierarchySystem) RemoveChildren(ecs *ecs.ECS) {
	hs.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		if !entry.Valid() {
			return
		}
		if pd, ok := getParentData(entry); ok {
			if ecs.World.Valid(pd.Parent) {
				return
			}
			c, ok := GetChildren(entry)
			if ok {
				for _, e := range c {
					RemoveRecursive(ecs.World.Entry(e))
				}
			}
			entry.Remove()
		}
	})
}
