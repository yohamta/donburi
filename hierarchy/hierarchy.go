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

func (ps *parent) RemoveChildren(ecs *ecs.ECS) {
	ps.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		if p, ok := GetParent(entry); ok {
			if ecs.World.Valid(p) {
				return
			}
			c, ok := GetChildren(entry)
			if ok {
				for _, e := range c {
					ecs.World.Remove(e)
				}
			}
			entry.Remove()
		}
	})
}
