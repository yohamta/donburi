package hierarchy

import (
	"github.com/yohamta/donburi"
)

type ChildrenData struct {
	Children []donburi.Entity
}

var Children = donburi.NewComponentType[ChildrenData]()

func GetChildren(entry *donburi.Entry) ([]donburi.Entity, bool) {
	if entry.HasComponent(Children) {
		c := donburi.Get[ChildrenData](entry, Children).Children
		return c, true
	}
	return nil, false
}
