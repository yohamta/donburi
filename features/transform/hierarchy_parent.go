package transform

import (
	"github.com/yohamta/donburi"
)

type hierarchyParentData struct {
	Parent *donburi.Entry
}

var hierarchyParentComponent = donburi.NewComponentType[hierarchyParentData]()

// getHierarchyParent returns a parent of the entry.
func getHierarchyParent(entry *donburi.Entry) (*donburi.Entry, bool) {
	if pd, ok := getHierarchyParentData(entry); ok {
		if pd.Parent.Valid() {
			return pd.Parent, true
		}
	}
	return nil, false
}

// mustGetHierarchyParent returns a parent of the entry.
func mustGetHierarchyParent(entry *donburi.Entry) *donburi.Entry {
	p := donburi.Get[hierarchyParentData](entry, hierarchyParentComponent)
	return p.Parent
}

// removeHierarchyChildrenRecursive removes children of the entry recursively.
func removeHierarchyChildrenRecursive(entry *donburi.Entry) {
	if hasHierarchyChildren(entry) {
		children, ok := getHierarchyChildren(entry)
		if ok {
			for _, c := range children {
				if c.Valid() {
					removeHierarchyChildrenRecursive(c)
					c.Remove()
				}
			}
		}
	}
}

func removeHierarchyParent(entry *donburi.Entry) {
	if !hasHierarchyParent(entry) {
		return
	}

	parent, ok := getHierarchyParent(entry)
	if !ok || !parent.Valid() {
		entry.RemoveComponent(hierarchyParentComponent)
		return
	}

	if children, ok := getHierarchyChildren(parent); ok {
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

// removeHierarchyRecursive removes the entry recursively.
func removeHierarchyRecursive(entry *donburi.Entry) {
	removeHierarchyChildrenRecursive(entry)
	removeHierarchyParent(entry)
	entry.Remove()
}

// appendHierarchyChild appends a child to the entry.
func appendHierarchyChild(parent *donburi.Entry, child *donburi.Entry) {
	setHierarchyParent(child, parent)
}

// FindChildrenWithComponent
func findHierarchyChildWithComponent(entry *donburi.Entry, componentType donburi.IComponentType) (*donburi.Entry, bool) {
	if children, ok := getHierarchyChildren(entry); ok {
		for _, c := range children {
			if c.Valid() && c.HasComponent(componentType) {
				return c, true
			}
		}
	}
	return nil, false
}

// setHierarchyParent sets a parent of the entry.
func setHierarchyParent(child *donburi.Entry, parent *donburi.Entry) {
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

// hasHierarchyParent returns true if the entry has a parent.
func hasHierarchyParent(entry *donburi.Entry) bool {
	return entry.HasComponent(hierarchyParentComponent)
}

// changeHierarchyParent changes a parent of the entry.
func changeHierarchyParent(child *donburi.Entry, newParent *donburi.Entry) {
	if !newParent.Valid() {
		panic("newParent is not valid")
	}
	if !child.Valid() {
		panic("child is not valid")
	}
	if !newParent.HasComponent(hierarchyChildrenComponent) {
		newParent.AddComponent(hierarchyChildrenComponent)
	}

	if oldParent, ok := getHierarchyParent(child); ok {
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

	setHierarchyParent(child, newParent)
}

func getHierarchyParentData(entry *donburi.Entry) (*hierarchyParentData, bool) {
	if hasHierarchyParent(entry) {
		p := donburi.Get[hierarchyParentData](entry, hierarchyParentComponent)
		return p, true
	}
	return nil, false
}
