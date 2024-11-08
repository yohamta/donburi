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

func RemoveParent(entry *donburi.Entry) {
	if !HasParent(entry) {
		return
	}

	parent, ok := GetParent(entry)
	if !ok || !parent.Valid() {
		entry.RemoveComponent(parentComponent)
		return
	}

	if children, ok := GetChildren(parent); ok {
		newChildren := make([]*donburi.Entry, 0, len(children))
		for _, child := range children {
			if child != entry {
				newChildren = append(newChildren, child)
			}
		}
		donburi.SetValue(parent, childrenComponent, childrenData{
			Children: newChildren,
		})
	}

	entry.RemoveComponent(parentComponent)
}

// RemoveRecursive removes the entry recursively.
func RemoveRecursive(entry *donburi.Entry) {
	RemoveChildrenRecursive(entry)
	RemoveParent(entry)
	entry.Remove()
}

// AppendChild appends a child to the entry.
func AppendChild(parent *donburi.Entry, child *donburi.Entry) {
	SetParent(child, parent)
}

// FindChildrenWithComponent
func FindChildWithComponent(entry *donburi.Entry, componentType donburi.IComponentType) (*donburi.Entry, bool) {
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

// ChangeParent changes a parent of the entry.
func ChangeParent(child *donburi.Entry, newParent *donburi.Entry) {
	if !newParent.Valid() {
		panic("newParent is not valid")
	}
	if !child.Valid() {
		panic("child is not valid")
	}
	if !newParent.HasComponent(childrenComponent) {
		newParent.AddComponent(childrenComponent)
	}

	if oldParent, ok := GetParent(child); ok {
		if oldParent == newParent {
			return
		}

		if oldParent.Valid() {
			oldChildren := donburi.Get[childrenData](oldParent, childrenComponent)
			for i, c := range oldChildren.Children {
				if c == child {
					oldChildren.Children = append(oldChildren.Children[:i], oldChildren.Children[i+1:]...)
					break
				}
			}
		}

		child.RemoveComponent(parentComponent)
	}

	SetParent(child, newParent)
}

func getParentData(entry *donburi.Entry) (*parentData, bool) {
	if HasParent(entry) {
		p := donburi.Get[parentData](entry, parentComponent)
		return p, true
	}
	return nil, false
}
