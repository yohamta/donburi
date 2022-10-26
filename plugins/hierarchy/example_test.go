package hierarchy

import (
	"testing"

	"github.com/yohamta/donburi"
	ecslib "github.com/yohamta/donburi/ecs"
)

func TestHierarchy(t *testing.T) {
	w := donburi.NewWorld()
	ecs := ecslib.NewECS(w)

	ecs.AddSystem(ecslib.System{
		Update: HierarchySystem.RemoveChildren,
	})

	parent := donburi.NewTag().SetName("parent")
	child := donburi.NewTag().SetName("child")
	grandChild := donburi.NewTag().SetName("grandChild")

	pe := w.Entry(w.Create(parent))
	ce := w.Entry(w.Create(child))
	ge := w.Entry(w.Create(grandChild))

	SetParent(ce, pe)
	SetParent(ge, ce)

	if p, ok := GetParent(ce); p != pe.Entity() || !ok {
		t.Errorf("expected parent entity %d, got %d", pe.Entity(), p)
	}

	if p, ok := GetParent(ge); p != ce.Entity() || !ok {
		t.Errorf("expected parent entity %d, got %d", ce.Entity(), p)
	}

	if _, ok := GetParent(pe); ok {
		t.Errorf("expected parent entity %d, got %d", donburi.Null, pe.Entity())
	}

	children, ok := GetChildren(pe)
	if !ok {
		t.Errorf("expected children, got nil")
	}
	if children[0] != ce.Entity() {
		t.Errorf("expected child entity %d, got %d", ce.Entity(), children[0])
	}

	children, ok = GetChildren(ce)
	if children[0] != ge.Entity() {
		t.Errorf("expected child entity %d, got %d", ge.Entity(), children[0])
	}

	pe.Remove()
	ecs.Update()

	if w.Valid(ce.Entity()) {
		t.Errorf("expected child entity %d to be removed", ce.Entity())
	}
	if w.Valid(ge.Entity()) {
		t.Errorf("expected grand child entity %d to be removed", ge.Entity())
	}
	if w.Len() != 0 {
		t.Errorf("expected world to be empty")
	}
}

func TestRemoveChildrenRecursive(t *testing.T) {
	w := donburi.NewWorld()

	parent := donburi.NewTag().SetName("parent")
	child := donburi.NewTag().SetName("child")
	grandChild := donburi.NewTag().SetName("grandChild")

	pe := w.Entry(w.Create(parent))
	ce := w.Entry(w.Create(child))
	ge := w.Entry(w.Create(grandChild))

	SetParent(ce, pe)
	SetParent(ge, ce)

	RemoveChildrenRecursive(pe)

	if w.Valid(ce.Entity()) {
		t.Errorf("expected child entity %d to be removed", ce.Entity())
	}
	if w.Valid(ge.Entity()) {
		t.Errorf("expected grand child entity %d to be removed", ge.Entity())
	}
	if !w.Valid(pe.Entity()) {
		t.Errorf("expected parent entity %d to be valid", pe.Entity())
	}
}

func TestRemoveRecursive(t *testing.T) {
	w := donburi.NewWorld()

	parent := donburi.NewTag().SetName("parent")
	child := donburi.NewTag().SetName("child")
	grandChild := donburi.NewTag().SetName("grandChild")

	pe := w.Entry(w.Create(parent))
	ce := w.Entry(w.Create(child))
	ge := w.Entry(w.Create(grandChild))

	AppendChild(pe, ce)
	AppendChild(ce, ge)

	RemoveRecursive(pe)

	if w.Valid(ce.Entity()) {
		t.Errorf("expected child entity %d to be removed", ce.Entity())
	}
	if w.Valid(ge.Entity()) {
		t.Errorf("expected grand child entity %d to be removed", ge.Entity())
	}
	if w.Len() != 0 {
		t.Errorf("expected world to be empty")
	}
}
