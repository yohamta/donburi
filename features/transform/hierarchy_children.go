package transform

import (
	"github.com/yohamta/donburi"
)

type hierarchyChildrenData struct {
	Children []*donburi.Entry
}

var hierarchyChildrenComponent = donburi.NewComponentType[hierarchyChildrenData]()

// getHierarchyChildren returns children of the entry.
func getHierarchyChildren(entry *donburi.Entry) ([]*donburi.Entry, bool) {
	if cd, ok := getHierarchyChildrenData(entry); ok {
		return cd.Children, true
	}
	return nil, false
}

// hasHierarchyChildren returns true if the entry has children.
func hasHierarchyChildren(entry *donburi.Entry) bool {
	return entry.HasComponent(hierarchyChildrenComponent)
}

func getHierarchyChildrenData(entry *donburi.Entry) (*hierarchyChildrenData, bool) {
	if hasHierarchyChildren(entry) {
		c := donburi.Get[hierarchyChildrenData](entry, hierarchyChildrenComponent)
		return c, true
	}
	return nil, false
}
