package hierarchy

import (
	"github.com/yohamta/donburi"
)

type childrenData struct {
	Children []*donburi.Entry
}

var childrenComponent = donburi.NewComponentType[childrenData]()

// GetChildren returns children of the entry.
func GetChildren(entry *donburi.Entry) ([]*donburi.Entry, bool) {
	if cd, ok := getChildrenData(entry); ok {
		return cd.Children, true
	}
	return nil, false
}

// MustGetChildren returns children of the entry.
func MustGetChildren(entry *donburi.Entry) []*donburi.Entry {
	c := donburi.Get[childrenData](entry, childrenComponent)
	return c.Children
}

// HasChildren returns true if the entry has children.
func HasChildren(entry *donburi.Entry) bool {
	return entry.HasComponent(childrenComponent)
}

func getChildrenData(entry *donburi.Entry) (*childrenData, bool) {
	if HasChildren(entry) {
		c := donburi.Get[childrenData](entry, childrenComponent)
		return c, true
	}
	return nil, false
}
