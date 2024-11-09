package transform

import (
	"github.com/yohamta/donburi"
)

type hierarchyParentData struct {
	Parent *donburi.Entry
}

var hierarchyParentComponent = donburi.NewComponentType[hierarchyParentData]()

// GetHierarchyParent returns a parent of the entry.
func GetHierarchyParent(entry *donburi.Entry) (*donburi.Entry, bool) {
	if pd, ok := getHierarchyParentData(entry); ok {
		if pd.Parent.Valid() {
			return pd.Parent, true
		}
	}
	return nil, false
}

// MustGetHierarchyParent returns a parent of the entry.
func MustGetHierarchyParent(entry *donburi.Entry) *donburi.Entry {
	p := donburi.Get[hierarchyParentData](entry, hierarchyParentComponent)
	return p.Parent
}

// RemoveHierarchyChildrenRecursive removes children of the entry recursively.
func RemoveHierarchyChildrenRecursive(entry *donburi.Entry) {
	if HasHierarchyChildren(entry) {
		children, ok := GetHierarchyChildren(entry)
		if ok {
			for _, c := range children {
				if c.Valid() {
					RemoveHierarchyChildrenRecursive(c)
					c.Remove()
				}
			}
		}
	}
}

func RemoveHierarchyParent(entry *donburi.Entry) {
	if !HasHierarchyParent(entry) {
		return
	}

	parent, ok := GetHierarchyParent(entry)
	if !ok || !parent.Valid() {
		entry.RemoveComponent(hierarchyParentComponent)
		return
	}

	if children, ok := GetHierarchyChildren(parent); ok {
		newChildren := make([]*donburi.Entry, 0, len(children))
		for _, child := range children {
			if child != entry {
				newChildren = append(newChildren, child)
			}
		}
		donburi.SetValue(parent, hierarchyChildrenComponent, hierarchyChildrenData{
			Children: newChildren,
		})
	}

	entry.RemoveComponent(hierarchyParentComponent)
}

// RemoveHierarchyRecursive removes the entry recursively.
func RemoveHierarchyRecursive(entry *donburi.Entry) {
	RemoveHierarchyChildrenRecursive(entry)
	RemoveHierarchyParent(entry)
	entry.Remove()
}

// AppendHierarchyChild appends a child to the entry.
func AppendHierarchyChild(parent *donburi.Entry, child *donburi.Entry) {
	SetHierarchyParent(child, parent)
}

// FindChildrenWithComponent
func FindHierarchyChildWithComponent(entry *donburi.Entry, componentType donburi.IComponentType) (*donburi.Entry, bool) {
	if children, ok := GetHierarchyChildren(entry); ok {
		for _, c := range children {
			if c.Valid() && c.HasComponent(componentType) {
				return c, true
			}
		}
	}
	return nil, false
}

// SetHierarchyParent sets a parent of the entry.
func SetHierarchyParent(child *donburi.Entry, parent *donburi.Entry) {
	if !parent.Valid() {
		panic("parent is not valid")
	}
	if !child.Valid() {
		panic("child is not valid")
	}
	if child.HasComponent(hierarchyParentComponent) {
		panic("child already has a parent")
	}
	if !parent.HasComponent(hierarchyChildrenComponent) {
		parent.AddComponent(hierarchyChildrenComponent)
	}
	child.AddComponent(hierarchyParentComponent)
	donburi.SetValue(child, hierarchyParentComponent, hierarchyParentData{Parent: parent})
	children := donburi.Get[hierarchyChildrenData](parent, hierarchyChildrenComponent)
	children.Children = append(children.Children, child)
}

// HasHierarchyParent returns true if the entry has a parent.
func HasHierarchyParent(entry *donburi.Entry) bool {
	return entry.HasComponent(hierarchyParentComponent)
}

// ChangeHierarchyParent changes a parent of the entry.
func ChangeHierarchyParent(child *donburi.Entry, newParent *donburi.Entry) {
	if !newParent.Valid() {
		panic("newParent is not valid")
	}
	if !child.Valid() {
		panic("child is not valid")
	}
	if !newParent.HasComponent(hierarchyChildrenComponent) {
		newParent.AddComponent(hierarchyChildrenComponent)
	}

	if oldParent, ok := GetHierarchyParent(child); ok {
		if oldParent == newParent {
			return
		}

		if oldParent.Valid() {
			oldChildren := donburi.Get[hierarchyChildrenData](oldParent, hierarchyChildrenComponent)
			for i, c := range oldChildren.Children {
				if c == child {
					oldChildren.Children = append(oldChildren.Children[:i], oldChildren.Children[i+1:]...)
					break
				}
			}
		}

		child.RemoveComponent(hierarchyParentComponent)
	}

	SetHierarchyParent(child, newParent)
}

func getHierarchyParentData(entry *donburi.Entry) (*hierarchyParentData, bool) {
	if HasHierarchyParent(entry) {
		p := donburi.Get[hierarchyParentData](entry, hierarchyParentComponent)
		return p, true
	}
	return nil, false
}
