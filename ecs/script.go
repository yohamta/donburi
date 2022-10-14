package ecs

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

// Script is a function that updates an entry.
type Script interface {
	Update(entry *donburi.Entry)
	Draw(entry *donburi.Entry, screen *ebiten.Image)
}

// ScriptSystem is a built-in system that manages scripts with queries.
type ScriptSystem struct {
	scripts []*script
}

// ScriptOpts represents options for a script.
type ScriptOpts struct {
	// Priority is the priority of the script. higher priority is executed first.
	Priority int
	// Image is the image that the script draws to.
	Image *ebiten.Image
}

// NewScriptSystem creates a new ScriptSystem.
func NewScriptSystem() *ScriptSystem {
	return &ScriptSystem{
		scripts: []*script{},
	}
}

type script struct {
	Query   *query.Query
	Script  Script
	Options *ScriptOpts
}

// AddScript adds a script to the system.
// Target entities are specified by the query.
func (ss *ScriptSystem) AddScript(q *query.Query, scr Script, opts *ScriptOpts) {
	if opts == nil {
		opts = &ScriptOpts{}
	}
	ss.addUpdater(&script{
		Query:   q,
		Script:  scr,
		Options: opts,
	})
}

func (ss *ScriptSystem) Update(ecs *ECS) {
	for _, script := range ss.scripts {
		script.Query.EachEntity(ecs.World, script.Script.Update)
	}
}

func (ss *ScriptSystem) Draw(ecs *ECS, screen *ebiten.Image) {
	for _, script := range ss.scripts {
		script.Query.EachEntity(ecs.World, func(entry *donburi.Entry) {
			if script.Options.Image != nil {
				script.Script.Draw(entry, script.Options.Image)
				return
			}
			script.Script.Draw(entry, screen)
		})
	}
}

func (ss *ScriptSystem) addUpdater(entry *script) {
	ss.scripts = append(ss.scripts, entry)
	sortEntryUpdater(ss.scripts)
}

func sortEntryUpdater(items []*script) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Options.Priority > items[j].Options.Priority
	})
}

var (
	_ = (System)((*ScriptSystem)(nil))
)
