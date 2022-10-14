package ecs

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

type Script struct {
	query    *query.Query
	callback ScriptCallback
	options  *ScriptOpts
}

func NewScript(q *query.Query, scr ScriptCallback, opts *ScriptOpts) *Script {
	if opts == nil {
		opts = &ScriptOpts{}
	}
	return &Script{
		query:    q,
		callback: scr,
		options:  opts,
	}
}

// ScriptCallback is a function that updates an entry.
type ScriptCallback interface {
	Update(entry *donburi.Entry)
	Draw(entry *donburi.Entry, screen *ebiten.Image)
}

// ScriptSystem is a built-in system that manages scripts with queries.
type ScriptSystem struct {
	scripts []*Script
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
		scripts: []*Script{},
	}
}

// AddScript adds a script to the system.
// Target entities are specified by the query.
func (ss *ScriptSystem) AddScript(scr *Script) {
	ss.scripts = append(ss.scripts, scr)
	sortScripts(ss.scripts)
}

func (ss *ScriptSystem) Update(ecs *ECS) {
	for _, script := range ss.scripts {
		script.query.EachEntity(ecs.World, script.callback.Update)
	}
}

func (ss *ScriptSystem) Draw(ecs *ECS, screen *ebiten.Image) {
	for _, script := range ss.scripts {
		script.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
			if script.options.Image != nil {
				script.callback.Draw(entry, script.options.Image)
				return
			}
			script.callback.Draw(entry, screen)
		})
	}
}

func sortScripts(scripts []*Script) {
	sort.Slice(scripts, func(i, j int) bool {
		return scripts[i].options.Priority > scripts[j].options.Priority
	})
}

var (
	_ = (System)((*ScriptSystem)(nil))
)
