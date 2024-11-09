package transform

import (
	"github.com/yohamta/donburi"
)

type hierarchyChildrenData struct {
	Children []*donburi.Entry
}

var hierarchyChildrenComponent = donburi.NewComponentType[hierarchyChildrenData]()

// GetHierarchyChildren returns children of the entry.
func GetHierarchyChildren(entry *donburi.Entry) ([]*donburi.Entry, bool) {
	if cd, ok := getHierarchyChildrenData(entry); ok {
		return cd.Children, true
	}
	return nil, false
}

// MustGetHierarchyChildren returns children of the entry.
func MustGetHierarchyChildren(entry *donburi.Entry) []*donburi.Entry {
	c := donburi.Get[hierarchyChildrenData](entry, hierarchyChildrenComponent)
	return c.Children
}

// HasHierarchyChildren returns true if the entry has children.
func HasHierarchyChildren(entry *donburi.Entry) bool {
	return entry.HasComponent(hierarchyChildrenComponent)
}

func getHierarchyChildrenData(entry *donburi.Entry) (*hierarchyChildrenData, bool) {
	if HasHierarchyChildren(entry) {
		c := donburi.Get[hierarchyChildrenData](entry, hierarchyChildrenComponent)
		return c, true
	}
	return nil, false
}
