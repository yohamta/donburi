package math

import "math"

const Epsilon = 0.00001

// Vec2 represents a 2D vector.
type Vec2 struct {
	X float64
	Y float64
}

// NewVec2 creates a new vector.
func NewVec2(x, y float64) Vec2 {
	return Vec2{
		X: x,
		Y: y,
	}
}

func (v Vec2) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vec2) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vec2) Normalized() Vec2 {
	m := v.Magnitude()

	if m > Epsilon {
		return v.DivScalar(m)
	}

	return v
}
func (v Vec2) MulScalar(scalar float64) Vec2 {
	return Vec2{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

func (v Vec2) DivScalar(scalar float64) Vec2 {
	return Vec2{
		X: v.X / scalar,
		Y: v.Y / scalar,
	}
}

func (v Vec2) Dot(other *Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Add adds x and y to the vector.
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{
		v.X + other.X,
		v.Y + other.Y,
	}
}

// Mul multiplies the vector by another vector.
func (v Vec2) Mul(other Vec2) Vec2 {
	return Vec2{
		v.X * other.X,
		v.Y * other.Y,
	}
}

// Sub subtracts x and y from the vector.
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{
		v.X - other.X,
		v.Y - other.Y,
	}
}

// RotateAround rotates the vector by the given angle around the given point.
func (v Vec2) RotateAround(point *Vec2, angle float64) Vec2 {
	x := v.X - point.X
	y := v.Y - point.Y
	return Vec2{
		X: x*math.Cos(angle) - y*math.Sin(angle) + point.X,
		Y: x*math.Sin(angle) + y*math.Cos(angle) + point.Y,
	}
}

// Equal returns true if the vector is equal to another vector.
func (v Vec2) Equal(other Vec2) bool {
	return v.X == other.X && v.Y == other.Y
}

// Roatate rotates the vector by the given angle.
func (v Vec2) Rotate(rad float64) Vec2 {
	return Vec2{
		X: v.X*math.Cos(rad) - v.Y*math.Sin(rad),
		Y: v.X*math.Sin(rad) + v.Y*math.Cos(rad),
	}
}

// Angle returns angle in radian to another vec2
func (v Vec2) Angle(other Vec2) float64 {
	x, y := v.X, v.Y
	x2, y2 := other.X, other.Y
	return math.Atan2(y2-y, x2-x)
}

// Distance returns the distance between the vector and another vector.
func (v Vec2) Distance(other Vec2) float64 {
	return math.Sqrt(math.Pow(v.X-other.X, 2) + math.Pow(v.Y-other.Y, 2))
}
