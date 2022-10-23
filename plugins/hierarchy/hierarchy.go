package hierarchy

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// HierarchySystem is a system that removes children of invalid parents.
var HierarchySystem = &parent{
	query: query.NewQuery(filter.Contains(parentComponent)),
}

// RemoveChildren removes children of invalid parents.
// This function is useful when you want to remove children of invalid parents.
// You don't need to use this system function when you use hierarchy.RemoveChildrenRecursive() function.
func (ps *parent) RemoveChildren(ecs *ecs.ECS) {
	ps.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		if !entry.Valid() {
			return
		}
		if p, ok := GetParent(entry); ok {
			if ecs.World.Valid(p) {
				return
			}
			c, ok := GetChildren(entry)
			if ok {
				for _, e := range c {
					defer ecs.World.Remove(e)
				}
			}
			entry.Remove()
		}
	})
}
