package hierarchy

import (
	"github.com/yohamta/donburi"
)

type parentData struct {
	Parent *donburi.Entry
}

var parentComponent = donburi.NewComponentType[parentData]()

// GetParent returns a parent of the entry.
func GetParent(entry *donburi.Entry) (*donburi.Entry, bool) {
	if pd, ok := getParentData(entry); ok {
		if pd.Parent.Valid() {
			return pd.Parent, true
		}
	}
	return nil, false
}

// MustGetParent returns a parent of the entry.
func MustGetParent(entry *donburi.Entry) *donburi.Entry {
	p := donburi.Get[parentData](entry, parentComponent)
	return p.Parent
}

// RemoveChildrenRecursive removes children of the entry recursively.
func RemoveChildrenRecursive(entry *donburi.Entry) {
	if HasChildren(entry) {
		children, ok := GetChildren(entry)
		if ok {
			for _, c := range children {
				if c.Valid() {
					RemoveChildrenRecursive(c)
					c.Remove()
				}
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

// FindChildrenWithComponent
func FindChildWithComponent(entry *donburi.Entry, componentType donburi.ComponentType) (*donburi.Entry, bool) {
	if children, ok := GetChildren(entry); ok {
		for _, c := range children {
			if c.Valid() && c.HasComponent(componentType) {
				return c, true
			}
		}
	}
	return nil, false
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
	donburi.SetValue(child, parentComponent, parentData{Parent: parent})
	children := donburi.Get[childrenData](parent, childrenComponent)
	children.Children = append(children.Children, child)
}

// HasParent returns true if the entry has a parent.
func HasParent(entry *donburi.Entry) bool {
	return entry.HasComponent(parentComponent)
}

func getParentData(entry *donburi.Entry) (*parentData, bool) {
	if HasParent(entry) {
		p := donburi.Get[parentData](entry, parentComponent)
		return p, true
	}
	return nil, false
}
