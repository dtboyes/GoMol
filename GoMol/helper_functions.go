package main

import "math"

// implementing basic vector operations to use in ray generation
func (v1 vec3) Add(v2 vec3) vec3 {
	return vec3{v1.x + v2.x, v1.y + v2.y, v1.z + v2.z}
}

func (v1 vec3) Subtract(v2 vec3) vec3 {
	return vec3{v1.x - v2.x, v1.y - v2.y, v1.z - v2.z}
}

func (v1 vec3) Cross(v2 vec3) vec3 {
	return vec3{v1.y*v2.z - v1.z*v2.y,
		v1.z*v2.x - v1.x*v2.z,
		v1.x*v2.y - v1.y*v2.x}
}

func (v1 vec3) Dot(v2 vec3) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func (v vec3) Scale(s float64) vec3 {
	return vec3{v.x * s, v.y * s, v.z * s}
}

func (v vec3) Normalize() vec3 {
	return v.Scale(1.0 / v.Length())
}

func (v vec3) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v vec3) EqualsZero() bool {
	return v.x == 0 && v.y == 0 && v.z == 0
}

// Ray getter functions
func (r Ray) getOrigin() vec3    { return r.origin }
func (r Ray) getDirection() vec3 { return r.direction }

// camera function definitions
func (c Camera) getPosition() vec3 { return c.position }

// light function definitions
func (l Light) getPosition() vec3 { return l.position }

// collision function definitions
func (c Collision) getPoint() vec3  { return c.point }
func (c Collision) getNormal() vec3 { return c.normal }

// atom functions
func CenterOfMass(atoms []*Atom) vec3 {
	sum := vec3{0, 0, 0}
	for i := 0; i < len(atoms); i++ {
		sum = sum.Add(vec3{atoms[i].x, atoms[i].y, atoms[i].z})
	}
	return sum.Scale(1.0 / float64(len(atoms)))
}
