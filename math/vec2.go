package engine

import "math"

const Epsilon = 0.00001

// Vec2 represents a 2D vector.
type Vec2 struct {
	X float64
	Y float64
}

// NewVec2 creates a new vector.
func NewVec2(x, y float64) *Vec2 {
	return &Vec2{
		X: x,
		Y: y,
	}
}

func (v *Vec2) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vec2) DivScalar(scalar float64) {
	v.X, v.Y = v.X/scalar, v.Y/scalar
}

func (v *Vec2) Normalize() {
	m := v.Magnitude()

	if m > Epsilon {
		v.DivScalar(m)
	}
}
func (v *Vec2) MulScalar(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vec2) Dot(other *Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Clone clones the vector.
func (v *Vec2) Clone() *Vec2 {
	return NewVec2(v.X, v.Y)
}

// Set sets the vector.
func (v *Vec2) Set(x, y float64) {
	v.X = x
	v.Y = y
}

// Add adds x and y to the vector.
func (v *Vec2) Add(x, y float64) {
	v.X += x
	v.Y += y
}

// Sub subtracts x and y from the vector.
func (v *Vec2) Sub(x, y float64) {
	v.X -= x
	v.Y -= y
}

// SetFrom sets the vector from another vector.
func (v *Vec2) SetFrom(v2 *Vec2) {
	v.X = v2.X
	v.Y = v2.Y
}

// RotateAround rotates the vector by the given angle around the given point.
func (v *Vec2) RotateAround(point *Vec2, angle float64) {
	x := v.X - point.X
	y := v.Y - point.Y
	v.X = x*math.Cos(angle) - y*math.Sin(angle) + point.X
	v.Y = x*math.Sin(angle) + y*math.Cos(angle) + point.Y
}

// AddFrom adds another vector to the vector.
func (v *Vec2) Equal(v2 *Vec2) bool {
	return v.X == v2.X && v.Y == v2.Y
}

// Roatate rotates the vector by the given angle.
func (v *Vec2) Rotate(rad float64) {
	x := v.X
	y := v.Y
	v.X = x*math.Cos(rad) - y*math.Sin(rad)
	v.Y = x*math.Sin(rad) + y*math.Cos(rad)
}

// Angle returns angle in radian to another vec2
func (v *Vec2) Angle(v2 *Vec2) float64 {
	x, y := v.X, v.Y
	x2, y2 := v2.X, v2.Y
	return math.Atan2(y2-y, x2-x)
}

// Distance returns the distance between the vector and another vector.
func (v *Vec2) Distance(v2 *Vec2) float64 {
	return math.Sqrt(math.Pow(v.X-v2.X, 2) + math.Pow(v.Y-v2.Y, 2))
}
