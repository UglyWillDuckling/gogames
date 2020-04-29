package vec3

import "math"

// Vector3 struct
type Vector3 struct {
	X, Y, Z float32
}

// Add will
func Add(a, b Vector3) Vector3 {
	return Vector3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Mult will
func Mult(a Vector3, b float32) Vector3 {
	return Vector3{a.X * b, a.Y * b, a.Z * b}
}

// Length will
func (a Vector3) Length() float32 {
	return float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y + a.Z*a.Z)))
}

// Distance will
func Distance(a, b Vector3) float32 {
	xDiff := a.X - b.X
	yDiff := a.Y - b.Y
	zDiff := a.Z - b.Z

	return float32(math.Sqrt(float64(xDiff*xDiff + yDiff*yDiff + zDiff*zDiff)))
}

// DistanceSquared will
func DistanceSquared(a, b Vector3) float32 {
	xDiff := a.X - b.X
	yDiff := a.Y - b.Y
	zDiff := a.Z - b.Z

	return xDiff*xDiff + yDiff*yDiff + zDiff*zDiff
}

// Normalize will resize the vector to a length of 1
func Normalize(a Vector3) Vector3 {
	len := a.Length()

	return Vector3{a.X / len, a.Y / len, a.Z / len}
}
