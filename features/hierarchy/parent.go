package hierarchy

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

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

// RemoveChildrenRecursive removes children of the entry recursively.
func RemoveChildrenRecursive(entry *donburi.Entry) {
	if entry.HasComponent(childrenComponent) {
		c, ok := GetChildren(entry)
		if ok {
			for _, e := range c {
				RemoveChildrenRecursive(entry.World.Entry(e))
				entry.World.Remove(e)
			}
		}
	}
}

// RemoveRecursive removes the entry recursively.
func RemoveRecursive(entry *donburi.Entry) {
	RemoveChildrenRecursive(entry)
	entry.Remove()
}

// AppendChild appends a child to the entry.
func AppendChild(parent *donburi.Entry, child *donburi.Entry) {
	SetParent(child, parent)
}

// SetParent sets a parent of the entry.
func SetParent(child *donburi.Entry, parent *donburi.Entry) {
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
