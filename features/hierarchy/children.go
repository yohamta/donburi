package hierarchy

import (
	"github.com/yohamta/donburi"
)

type childrenData struct {
	Children []donburi.Entity
}

var childrenComponent = donburi.NewComponentType[childrenData]()

// GetChildren returns children of the entry.
func GetChildren(entry *donburi.Entry) ([]donburi.Entity, bool) {
	if cd, ok := getChildrenData(entry); ok {
		return cd.Children, true
	}
	return nil, false
}

// HasChildren returns true if the entry has children.
func HasChildren(entry *donburi.Entry) bool {
	if entry.Valid() {
		return entry.HasComponent(childrenComponent)
	}
	return false
}

func getChildrenData(entry *donburi.Entry) (*childrenData, bool) {
	if entry.Valid() {
		c := donburi.Get[childrenData](entry, childrenComponent)
		return c, true
	}
	return nil, false
}
