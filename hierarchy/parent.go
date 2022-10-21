package hierarchy

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// ParentSystem is a system that removes children of invalid parents.
var ParentSystem = &parent{
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

type parentData struct {
	Parent donburi.Entity
}

var parentComponent = donburi.NewComponentType[parentData]()

// GetParent returns a parent of the entry.
func GetParent(entry *donburi.Entry) (donburi.Entity, bool) {
	if entry.HasComponent(parentComponent) {
		p := donburi.Get[parentData](entry, parentComponent).Parent
		return p, true
	}
	return donburi.Null, false
}

// SetParent sets a parent of the entry.
func SetParent(parent *donburi.Entry, child *donburi.Entry) {
	if !parent.Valid() {
		panic("parent is not valid")
	}
	if !child.Valid() {
		panic("child is not valid")
	}
	if child.HasComponent(parentComponent) {
		panic("child already has a parent")
	}
	if !parent.HasComponent(childrenComponent) {
		parent.AddComponent(childrenComponent)
	}
	child.AddComponent(parentComponent)
	donburi.SetValue(child, parentComponent, parentData{Parent: parent.Entity()})
	children := donburi.Get[childrenData](parent, childrenComponent)
	children.Children = append(children.Children, child.Entity())
}

type parent struct {
	query *query.Query
}
