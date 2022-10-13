package ecs

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

type Script interface {
	Update(entry *donburi.Entry)
	Draw(entry *donburi.Entry, screen *ebiten.Image)
}

type ScriptSystem struct {
	scripts []scriptEntry
}

type ScriptOptions struct {
	ImageToDraw *ebiten.Image
}

var (
	_ = (Updater)((*ScriptSystem)(nil))
	_ = (Drawer)((*ScriptSystem)(nil))
)

func NewScriptSystem() *ScriptSystem {
	return &ScriptSystem{
		scripts: []scriptEntry{},
	}
}

type scriptEntry struct {
	Query    query.Query
	Script   Script
	Priority int
	Options  *ScriptOptions
}

func (s *ScriptSystem) AddScript(q query.Query, script Script, priority int, opts *ScriptOptions) {
	entry := scriptEntry{q, script, priority, opts}
	if entry.Options == nil {
		entry.Options = &ScriptOptions{}
	}
	s.scripts = append(s.scripts, entry)

	// sort script entries by priority. higher priority is executed first.
	sortScriptEntries(s.scripts)
}

func (s *ScriptSystem) Update(ecs *ECS) {
	for _, script := range s.scripts {
		script.Query.EachEntity(ecs.World, script.Script.Update)
	}
}

func (s *ScriptSystem) Draw(ecs *ECS, screen *ebiten.Image) {
	for _, script := range s.scripts {
		script.Query.EachEntity(ecs.World, func(entry *donburi.Entry) {
			if script.Options.ImageToDraw != nil {
				script.Script.Draw(entry, script.Options.ImageToDraw)
			}
			script.Script.Draw(entry, screen)
		})
	}
}

// sortScriptEntries sorts script entries by priority. higher priority is executed first.
func sortScriptEntries(entries []scriptEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Priority > entries[j].Priority
	})
}
