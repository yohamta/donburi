package ecs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func TestArchetype(t *testing.T) {
	ca := donburi.NewTag().SetName("A")
	cb := donburi.NewTag().SetName("B")
	cc := donburi.NewTag().SetName("C")
	cd := donburi.NewTag().SetName("D")

	archetypeA := ecs.NewArchetype(ca)
	archetypeB := ecs.NewArchetype(cb, cc)

	ecsInstance := ecs.NewECS(donburi.NewWorld())

	e1 := archetypeA.Spawn(ecsInstance)
	assert.True(t, e1.HasComponent(ca))
	assert.False(t, e1.HasComponent(cb))
	assert.False(t, e1.HasComponent(cc))

	e2 := archetypeB.Spawn(ecsInstance)
	assert.False(t, e2.HasComponent(ca))
	assert.True(t, e2.HasComponent(cb))
	assert.True(t, e2.HasComponent(cc))

	e3 := archetypeA.Spawn(ecsInstance, cd)
	assert.True(t, e3.HasComponent(ca))
	assert.False(t, e3.HasComponent(cb))
	assert.False(t, e3.HasComponent(cc))
	assert.True(t, e3.HasComponent(cd))
}
