package debug

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/internal/storage"
)

type EntityCounts struct {
	Archetype *storage.Archetype
	Count     int
}

func (e EntityCounts) String() string {
	return fmt.Sprintf("Archetype %v has %d entities", e.Archetype.Layout(), e.Count)
}

func PrintEntityCounts(w donburi.World) {
	var out bytes.Buffer
	out.WriteString("Entity Counts:\n")
	for _, c := range GetEntityCounts(w) {
		out.WriteString(c.String())
		out.WriteString("\n")
	}
	out.WriteString("\n")
	fmt.Println(out.String())
}

func GetEntityCounts(w donburi.World) []EntityCounts {
	archetypes := w.Archetypes()
	ret := []EntityCounts{}
	for _, a := range archetypes {
		if a.Count() != 0 {
			ret = append(ret, EntityCounts{Archetype: a, Count: a.Count()})
		}
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Count > ret[j].Count
	})
	return ret
}
