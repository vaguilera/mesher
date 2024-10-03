package vec3

import (
	"github.com/chewxy/math32"
)

type Vec3 struct {
	X, Y, Z float32
}

func CrossProduct(a, b Vec3) Vec3 {
	return Vec3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (v *Vec3) Normalize() {
	length := math32.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	v = &Vec3{v.X / length, v.Y / length, v.Z / length}
}

func Add(v, other Vec3) Vec3 {
	return Vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func Sub(v, other Vec3) Vec3 {
	return Vec3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}
