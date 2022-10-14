package ecs

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

// Script is a arbitrary script that can be executed in the world.
type Script interface {
	Update(entry *donburi.Entry)
	Draw(entry *donburi.Entry, screen *ebiten.Image)
}

// ScriptSystem is a built-in system that manages scripts with queries.
type ScriptSystem struct {
	scripts []scriptEntry
}

// ScriptOpts represents options for a script.
// The Priority is the priority of the script. higher priority is executed first.
type ScriptOpts struct {
	ImageToDraw *ebiten.Image
	Priority    int
}

// NewScriptSystem creates a new ScriptSystem.
func NewScriptSystem() *ScriptSystem {
	return &ScriptSystem{
		scripts: []scriptEntry{},
	}
}

type scriptEntry struct {
	Query   query.Query
	Script  Script
	Options *ScriptOpts
}

// AddScript adds a script to the system.
// Target entities are specified by the query.
func (s *ScriptSystem) AddScript(q query.Query, script Script, opts *ScriptOpts) {
	if opts == nil {
		opts = &ScriptOpts{}
	}
	entry := scriptEntry{q, script, opts}
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
				return
			}
			script.Script.Draw(entry, screen)
		})
	}
}

// sortScriptEntries sorts script entries by priority. higher priority is executed first.
func sortScriptEntries(entries []scriptEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Options.Priority > entries[j].Options.Priority
	})
}

var (
	_ = (Updater)((*ScriptSystem)(nil))
	_ = (Drawer)((*ScriptSystem)(nil))
)
