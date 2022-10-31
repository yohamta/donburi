// This code is adapted from https://github.com/m110/airplanes (author: m110)
package transform

import (
	"math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/hierarchy"
	dmath "github.com/yohamta/donburi/features/math"
)

// TransformData is a data of transform component.
type TransformData struct {
	LocalPosition dmath.Vec2
	LocalRotation float64
	LocalScale    dmath.Vec2

	hasPrent bool
}

var defaultValue = TransformData{
	LocalPosition: dmath.Vec2{},
	LocalRotation: 0,
	LocalScale:    dmath.Vec2{X: 1, Y: 1},
}

// Reset resets the entry's transform component to default value.
func Reset(entry *donburi.Entry) {
	if !entry.HasComponent(Transform) {
		panic("entry does not have transform component")
	}
	donburi.SetValue(entry, Transform, defaultValue)
}

// AppendChild appends child to the entry.
func AppendChild(parent, child *donburi.Entry, keepWorldPosition bool) {
	if !parent.HasComponent(Transform) {
		panic("parent does not have transform component")
	}
	hierarchy.AppendChild(parent, child)
	SetParent(child, parent, keepWorldPosition)
}

// SetParent sets parent to the entry.
func SetParent(entry, parent *donburi.Entry, keepWorldPosition bool) {
	if !entry.HasComponent(Transform) {
		panic("entry does not have transform component")
	}
	d := GetTransform(entry)
	d.hasPrent = true
	if keepWorldPosition {
		parentPos := WorldPosition(parent)

		d.LocalPosition = d.LocalPosition.Sub(parentPos)
		d.LocalRotation -= WorldRotation(parent)

		ws := WorldScale(parent)
		d.LocalScale.X = d.LocalScale.X / ws.X
		d.LocalScale.Y = d.LocalScale.Y / ws.Y
	}
}

// GetParent returns parent of the entry.
func GetParent(entry *donburi.Entry) (*donburi.Entry, bool) {
	d := GetTransform(entry)
	if !d.hasPrent {
		return nil, false
	}
	return hierarchy.GetParent(entry)
}

// GetChildren returns children of the entry.
func GetChildren(entry *donburi.Entry) ([]*donburi.Entry, bool) {
	return hierarchy.GetChildren(entry)
}

// RemoveRecursive removes the entry and its children recursively.
func RemoveRecursive(entry *donburi.Entry) {
	hierarchy.RemoveRecursive(entry)
}

// RemoveChildrenRecursive removes children recursively.
func RemoveChildrenRecursive(entry *donburi.Entry) {
	hierarchy.RemoveChildrenRecursive(entry)
}

// FindChildWithComponent finds child with specified component.
func FindChildWithComponent(entry *donburi.Entry, componentType *donburi.ComponentType) (*donburi.Entry, bool) {
	return hierarchy.FindChildWithComponent(entry, componentType)
}

// SetWorldPosition sets world position to the entry.
func SetWorldPosition(entry *donburi.Entry, pos dmath.Vec2) {
	d := GetTransform(entry)
	if !d.hasPrent {
		d.LocalPosition = pos
		return
	}

	parent := hierarchy.MustGetParent(entry)
	parentPos := WorldPosition(parent)
	d.LocalPosition.X = pos.X - parentPos.X
	d.LocalPosition.Y = pos.Y - parentPos.Y
}

// WorldPosition returns world position of the entry.
func WorldPosition(entry *donburi.Entry) dmath.Vec2 {
	d := GetTransform(entry)
	if !d.hasPrent {
		return d.LocalPosition
	}

	parent := hierarchy.MustGetParent(entry)
	parentPos := WorldPosition(parent)
	return parentPos.Add(d.LocalPosition)
}

// SetWorldRotation sets world rotation to the entry.
func SetWorldRotation(entry *donburi.Entry, rotation float64) {
	d := GetTransform(entry)
	if !d.hasPrent {
		d.LocalRotation = rotation
		return
	}

	parent := hierarchy.MustGetParent(entry)
	d.LocalRotation = rotation - WorldRotation(parent)
}

// SetWorldScale sets world scale to the entry.
func SetWorldScale(entry *donburi.Entry, scale dmath.Vec2) {
	d := GetTransform(entry)
	if !d.hasPrent {
		d.LocalScale = scale
		return
	}

	parent := hierarchy.MustGetParent(entry)
	parentScale := WorldScale(parent)
	d.LocalScale.X = scale.X / parentScale.X
	d.LocalScale.Y = scale.Y / parentScale.Y
}

// WorldRotation returns world rotation of the entry.
func WorldRotation(entry *donburi.Entry) float64 {
	d := GetTransform(entry)
	if !d.hasPrent {
		return d.LocalRotation
	}

	parent := hierarchy.MustGetParent(entry)
	rot := WorldRotation(parent)
	return rot + d.LocalRotation
}

// WorldScale returns world scale of the entry.
func WorldScale(entry *donburi.Entry) dmath.Vec2 {
	d := GetTransform(entry)
	if !d.hasPrent {
		return d.LocalScale
	}

	parent := hierarchy.MustGetParent(entry)
	parentScale := WorldScale(parent)
	return parentScale.Mul(d.LocalScale)
}

// Right returns right vector of the entry.
func Right(entry *donburi.Entry) dmath.Vec2 {
	radians := dmath.ToRadians(WorldRotation(entry))
	return dmath.Vec2{
		X: math.Cos(radians),
		Y: math.Sin(radians),
	}
}

// Up returns up vector of the entry.
func Up(entry *donburi.Entry) dmath.Vec2 {
	radians := dmath.ToRadians(WorldRotation(entry) - 90)
	return dmath.Vec2{
		X: math.Cos(radians),
		Y: math.Sin(radians),
	}
}

// LookAt looks at the target.
func LookAt(entry *donburi.Entry, target dmath.Vec2) {
	x := target.X - WorldPosition(entry).X
	y := target.Y - WorldPosition(entry).Y
	radians := math.Atan2(y, x)
	SetWorldRotation(entry, radians)
}

var Transform = donburi.NewComponentType[TransformData]()

func GetTransform(entry *donburi.Entry) *TransformData {
	return donburi.Get[TransformData](entry, Transform)
}
