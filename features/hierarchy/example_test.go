package hierarchy

import (
	"testing"

	"github.com/yohamta/donburi"
	ecslib "github.com/yohamta/donburi/ecs"
)

func TestHierarchy(t *testing.T) {
	w := donburi.NewWorld()
	ecs := ecslib.NewECS(w)

	ecs.AddSystem(HierarchySystem.RemoveChildren)

	parent := donburi.NewTag().SetName("parent")
	child := donburi.NewTag().SetName("child")
	grandChild := donburi.NewTag().SetName("grandChild")

	pe := w.Entry(w.Create(parent))
	ce := w.Entry(w.Create(child))
	ge := w.Entry(w.Create(grandChild))

	SetParent(ce, pe)
	SetParent(ge, ce)

	testChildren(t, []childrenTest{
		{
			Parent:   pe,
			Children: []*donburi.Entry{ce},
		},
		{
			Parent:   ce,
			Children: []*donburi.Entry{ge},
		},
	})

	if p, ok := GetParent(ce); p.Entity() != pe.Entity() || !ok {
		t.Errorf("expected parent entity %d, got %d", pe.Entity(), p.Entity())
	}

	if p, ok := GetParent(ge); p.Entity() != ce.Entity() || !ok {
		t.Errorf("expected parent entity %d, got %d", ce.Entity(), p.Entity())
	}

	if HasParent(pe) {
		t.Errorf("expected parent entity %d, got %d", donburi.Null, pe.Entity())
	}

	children, ok := GetChildren(pe)
	if !ok {
		t.Errorf("expected children, got nil")
	}
	if children[0].Entity() != ce.Entity() {
		t.Errorf("expected child entity %d, got %d", ce.Entity(), children[0].Entity())
	}

	children, ok = GetChildren(ce)
	if children[0].Entity() != ge.Entity() {
		t.Errorf("expected child entity %d, got %d", ge.Entity(), children[0].Entity())
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

	testChildren(t, []childrenTest{
		{
			Parent:   pe,
			Children: []*donburi.Entry{ce},
		},
		{
			Parent:   ce,
			Children: []*donburi.Entry{ge},
		},
	})

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

	testChildren(t, []childrenTest{
		{
			Parent:   pe,
			Children: []*donburi.Entry{ce},
		},
		{
			Parent:   ce,
			Children: []*donburi.Entry{ge},
		},
	})

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

func TestFindChildren(t *testing.T) {
	w := donburi.NewWorld()

	tagToFind := donburi.NewTag().SetName("tag")

	parent := donburi.NewTag().SetName("parent")
	child := donburi.NewTag().SetName("child")

	pe := w.Entry(w.Create(parent))
	ce := w.Entry(w.Create(child, tagToFind))

	SetParent(ce, pe)

	found, ok := FindChildWithComponent(pe, tagToFind)
	if !ok {
		t.Errorf("expected to find child with component")
	}
	if found.Entity() != ce.Entity() {
		t.Errorf("expected to find child entity %d, got %d", ce.Entity(), found.Entity())
	}
}

func TestChangeParent(t *testing.T) {
	w := donburi.NewWorld()

	parent1 := donburi.NewTag().SetName("parent1")
	parent2 := donburi.NewTag().SetName("parent2")
	child := donburi.NewTag().SetName("child")

	p1e := w.Entry(w.Create(parent1))
	p2e := w.Entry(w.Create(parent2))
	ce := w.Entry(w.Create(child))

	// no parent exists
	ChangeParent(ce, p1e)
	testChildren(t, []childrenTest{
		{
			Parent:   p1e,
			Children: []*donburi.Entry{ce},
		},
		{
			Parent:   p2e,
			Children: []*donburi.Entry{},
		},
	})

	// change to same parent
	ChangeParent(ce, p1e)
	testChildren(t, []childrenTest{
		{
			Parent:   p1e,
			Children: []*donburi.Entry{ce},
		},
		{
			Parent:   p2e,
			Children: []*donburi.Entry{},
		},
	})

	// change parent
	ChangeParent(ce, p2e)
	testChildren(t, []childrenTest{
		{
			Parent:   p1e,
			Children: []*donburi.Entry{},
		},
		{
			Parent:   p2e,
			Children: []*donburi.Entry{ce},
		},
	})
}

type childrenTest struct {
	Parent   *donburi.Entry
	Children []*donburi.Entry
}

func testChildren(t *testing.T, tests []childrenTest) {
	for _, test := range tests {
		children, ok := GetChildren(test.Parent)
		if !ok {
			t.Errorf("expected children, got nil")
		}
		if len(children) != len(test.Children) {
			t.Errorf("expected %d children, got %d", len(test.Children), len(children))
		}
		for i, c := range children {
			if c.Entity() != test.Children[i].Entity() {
				t.Errorf("expected child entity %d, got %d", test.Children[i].Entity(), c.Entity())
			}
		}
	}
}
