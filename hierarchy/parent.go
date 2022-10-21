package hierarchy

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type ParentData struct {
	Parent donburi.Entity
}

var Parent = donburi.NewComponentType[ParentData]()

func GetParent(entry *donburi.Entry) (donburi.Entity, bool) {
	if entry.HasComponent(Parent) {
		p := donburi.Get[ParentData](entry, Parent).Parent
		return p, true
	}
	return donburi.Null, false
}

func SetParent(parent *donburi.Entry, child *donburi.Entry) {
	if !parent.Valid() {
		panic("parent is not valid")
	}
	if !child.Valid() {
		panic("child is not valid")
	}
	if child.HasComponent(Parent) {
		panic("child already has a parent")
	}
	if !parent.HasComponent(Children) {
		parent.AddComponent(Children)
	}
	child.AddComponent(Parent)
	donburi.SetValue(child, Parent, ParentData{Parent: parent.Entity()})
	children := donburi.Get[ChildrenData](parent, Children)
	children.Children = append(children.Children, child.Entity())
}

type parent struct {
	query *query.Query
}

var ParentSystem = &parent{
	query: query.NewQuery(filter.Contains(Parent)),
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
