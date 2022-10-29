package transform_test

import (
	"testing"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func TestTransform(t *testing.T) {
	w := donburi.NewWorld()

	parent := w.Entry(w.Create(transform.Transform))
	transform.Reset(parent)
	transform.SetWorldPosition(parent, dmath.Vec2{X: 1, Y: 2})
	transform.SetWorldScale(parent, dmath.Vec2{X: 2, Y: 3})

	child := w.Entry(w.Create(transform.Transform))
	donburi.SetValue(child, transform.Transform, transform.TransformData{
		LocalPosition: dmath.Vec2{X: 1, Y: 2},
		LocalRotation: 90,
		LocalScale:    dmath.Vec2{X: 2, Y: 3},
	})

	transform.AppendChild(parent, child, false)

	testWorldTransform(t, child, &testTransform{
		Position: dmath.Vec2{X: 2, Y: 4},
		Rotation: 90,
		Scale:    dmath.Vec2{X: 4, Y: 9},
	})

	transform.SetWorldRotation(parent, 90)
	transform.SetWorldPosition(parent, dmath.Vec2{X: 0, Y: 0})
	transform.SetWorldScale(parent, dmath.Vec2{X: 1, Y: 1})

	testWorldTransform(t, child, &testTransform{
		Position: dmath.Vec2{X: 1, Y: 2},
		Rotation: 180,
		Scale:    dmath.Vec2{X: 2, Y: 3},
	})
}

func TestTransformKeepWorldPosition(t *testing.T) {
	w := donburi.NewWorld()

	parent := w.Entry(w.Create(transform.Transform))
	transform.Reset(parent)

	transform.SetWorldPosition(parent, dmath.Vec2{X: 1, Y: 2})
	transform.SetWorldRotation(parent, 90)
	transform.SetWorldScale(parent, dmath.Vec2{X: 2, Y: 2})

	child := w.Entry(w.Create(transform.Transform))
	donburi.SetValue(child, transform.Transform, transform.TransformData{
		LocalPosition: dmath.Vec2{X: 1, Y: 2},
		LocalRotation: 90,
		LocalScale:    dmath.Vec2{X: 1.5, Y: 1.5},
	})

	transform.AppendChild(parent, child, true)

	testWorldTransform(t, child, &testTransform{
		Position: dmath.Vec2{X: 1, Y: 2},
		Rotation: 90,
		Scale:    dmath.Vec2{X: 1.5, Y: 1.5},
	})
}

type testTransform struct {
	Position dmath.Vec2
	Rotation float64
	Scale    dmath.Vec2
}

func testWorldTransform(t *testing.T, entry *donburi.Entry, test *testTransform) {
	t.Helper()

	pos := transform.WorldPosition(entry)
	rot := transform.WorldRotation(entry)
	scale := transform.WorldScale(entry)

	if pos.X != test.Position.X || pos.Y != test.Position.Y {
		t.Errorf("expected position %v, but got %v", test.Position, pos)
	}

	if rot != test.Rotation {
		t.Errorf("expected rotation %v, but got %v", test.Rotation, rot)
	}

	if scale.X != test.Scale.X || scale.Y != test.Scale.Y {
		t.Errorf("expected scale %v, but got %v", test.Scale, scale)
	}
}
