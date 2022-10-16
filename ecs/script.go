package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

// UpdateScript is a script that updates entities.
type UpdateScript func(entry *donburi.Entry)

// DrawScript is a script that draws entities.
type DrawScript func(entry *donburi.Entry, screen *ebiten.Image)

// Script represents a script.
type Script struct {
	Update    UpdateScript
	DrawLayer DrawLayer
	Draw      DrawScript
	Query     *query.Query
}

type updateScript struct {
	query    *query.Query
	callback UpdateScript
}

type drawScript struct {
	query    *query.Query
	callback DrawScript
}

type scriptSystem struct {
	updateScripts []*updateScript
	drawScripts   []*drawScript
}

func newScriptSystem() *scriptSystem {
	return &scriptSystem{
		updateScripts: []*updateScript{},
		drawScripts:   []*drawScript{},
	}
}

func (ss *scriptSystem) AddUpdateScript(s UpdateScript, q *query.Query) {
	ss.updateScripts = append(ss.updateScripts, &updateScript{
		query:    q,
		callback: s,
	})
}

func (ss *scriptSystem) AddDrawScript(s DrawScript, q *query.Query) {
	ss.drawScripts = append(ss.drawScripts, &drawScript{
		query:    q,
		callback: s,
	})
}

func (ss *scriptSystem) Update(ecs *ECS) {
	for _, script := range ss.updateScripts {
		script.query.EachEntity(ecs.World, script.callback)
	}
}

func (ss *scriptSystem) Draw(ecs *ECS, screen *ebiten.Image) {
	for _, script := range ss.drawScripts {
		script.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
			script.callback(entry, screen)
		})
	}
}
