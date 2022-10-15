package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

type script struct {
	query    *query.Query
	callback Script
}

func newScript(scr Script, q *query.Query) *script {
	return &script{
		query:    q,
		callback: scr,
	}
}

// Script is a function that runs on the entities matched by the query.
type Script interface {
	Update(entry *donburi.Entry)
	Draw(entry *donburi.Entry, screen *ebiten.Image)
}

type scriptSystem struct {
	scripts []*script
}

func newScriptSystem() *scriptSystem {
	return &scriptSystem{
		scripts: []*script{},
	}
}

func (ss *scriptSystem) AddScript(scr *script) {
	ss.scripts = append(ss.scripts, scr)
}

func (ss *scriptSystem) Update(ecs *ECS) {
	for _, script := range ss.scripts {
		script.query.EachEntity(ecs.World, script.callback.Update)
	}
}

func (ss *scriptSystem) Draw(ecs *ECS, screen *ebiten.Image) {
	for _, script := range ss.scripts {
		script.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
			script.callback.Draw(entry, screen)
		})
	}
}

var (
	_ = (System)((*scriptSystem)(nil))
)
