// Package math implements 3D vector math for Minecraft.
package math

import (
	"math"
)

// Vec3 represents a 3D vector.
type Vec3 struct {
	X float64
	Y float64
	Z float64
}

// NewVec3 creates a new 3D vector.
func NewVec3(x, y, z float64) *Vec3 {
	return &Vec3{X: x, Y: y, Z: z}
}

// Add adds another vector to this vector.
func (v *Vec3) Add(o *Vec3) *Vec3 {
	return &Vec3{
		X: v.X + o.X,
		Y: v.Y + o.Y,
		Z: v.Z + o.Z,
	}
}

// Sub subtracts another vector from this vector.
func (v *Vec3) Sub(o *Vec3) *Vec3 {
	return &Vec3{
		X: v.X - o.X,
		Y: v.Y - o.Y,
		Z: v.Z - o.Z,
	}
}

// Scale scales the vector by a scalar.
func (v *Vec3) Scale(s float64) *Vec3 {
	return &Vec3{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}

// Length returns the length (magnitude) of the vector.
func (v *Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// LengthSquared returns the squared length.
// This is faster than Length() and useful for comparisons.
func (v *Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Normalize returns a normalized (unit length) vector.
func (v *Vec3) Normalize() *Vec3 {
	len := v.Length()
	if len == 0 {
		return &Vec3{X: 0, Y: 0, Z: 0}
	}
	return &Vec3{
		X: v.X / len,
		Y: v.Y / len,
		Z: v.Z / len,
	}
}

// DistanceTo returns the distance to another vector.
func (v *Vec3) DistanceTo(o *Vec3) float64 {
	dx := v.X - o.X
	dy := v.Y - o.Y
	dz := v.Z - o.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// DistanceSquaredTo returns the squared distance to another vector.
func (v *Vec3) DistanceSquaredTo(o *Vec3) float64 {
	dx := v.X - o.X
	dy := v.Y - o.Y
	dz := v.Z - o.Z
	return dx*dx + dy*dy + dz*dz
}

// Dot returns the dot product with another vector.
func (v *Vec3) Dot(o *Vec3) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

// Cross returns the cross product with another vector.
func (v *Vec3) Cross(o *Vec3) *Vec3 {
	return &Vec3{
		X: v.Y*o.Z - v.Z*o.Y,
		Y: v.Z*o.X - v.X*o.Z,
		Z: v.X*o.Y - v.Y*o.X,
	}
}

// Floor returns a vector with each component floored to an integer.
func (v *Vec3) Floor() *Vec3 {
	return &Vec3{
		X: math.Floor(v.X),
		Y: math.Floor(v.Y),
		Z: math.Floor(v.Z),
	}
}

// Ceil returns a vector with each component ceilinged to an integer.
func (v *Vec3) Ceil() *Vec3 {
	return &Vec3{
		X: math.Ceil(v.X),
		Y: math.Ceil(v.Y),
		Z: math.Ceil(v.Z),
	}
}

// Abs returns a vector with absolute values of each component.
func (v *Vec3) Abs() *Vec3 {
	return &Vec3{
		X: math.Abs(v.X),
		Y: math.Abs(v.Y),
		Z: math.Abs(v.Z),
	}
}

// Min returns the minimum value among all components.
func (v *Vec3) Min() float64 {
	return math.Min(v.X, math.Min(v.Y, v.Z))
}

// Max returns the maximum value among all components.
func (v *Vec3) Max() float64 {
	return math.Max(v.X, math.Max(v.Y, v.Z))
}

// Clone returns a copy of the vector.
func (v *Vec3) Clone() *Vec3 {
	return &Vec3{X: v.X, Y: v.Y, Z: v.Z}
}

// Equals returns true if the vectors are approximately equal.
func (v *Vec3) Equals(o *Vec3) bool {
	const epsilon = 1e-9
	dx := v.X - o.X
	dy := v.Y - o.Y
	dz := v.Z - o.Z
	return dx*dx < epsilon && dy*dy < epsilon && dz*dz < epsilon
}

// Set sets the vector components.
func (v *Vec3) Set(x, y, z float64) {
	v.X = x
	v.Y = y
	v.Z = z
}

// Offset offsets the vector by the given amounts.
func (v *Vec3) Offset(x, y, z float64) {
	v.X += x
	v.Y += y
	v.Z += z
}

// ToBlockPos returns the block position (integer coordinates).
func (v *Vec3) ToBlockPos() *BlockPos {
	return &BlockPos{
		X: int(math.Floor(v.X)),
		Y: int(math.Floor(v.Y)),
		Z: int(math.Floor(v.Z)),
	}
}

// BlockPos represents a block position (integer coordinates).
type BlockPos struct {
	X int
	Y int
	Z int
}

// NewBlockPos creates a new block position.
func NewBlockPos(x, y, z int) *BlockPos {
	return &BlockPos{X: x, Y: y, Z: z}
}

// ToVec3 converts to a Vec3 (block center).
func (b *BlockPos) ToVec3() *Vec3 {
	return &Vec3{
		X: float64(b.X) + 0.5,
		Y: float64(b.Y) + 0.5,
		Z: float64(b.Z) + 0.5,
	}
}

// DistanceTo returns the distance to another block position.
func (b *BlockPos) DistanceTo(o *BlockPos) float64 {
	dx := float64(b.X - o.X)
	dy := float64(b.Y - o.Y)
	dz := float64(b.Z - o.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// Equals returns true if the positions are equal.
func (b *BlockPos) Equals(o *BlockPos) bool {
	return b.X == o.X && b.Y == o.Y && b.Z == o.Z
}
