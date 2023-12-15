package main

import "math"

// implementing basic vector operations to use in ray generation

// Add returns the sum of two vectors
func (v1 vec3) Add(v2 vec3) vec3 {
	return vec3{v1.x + v2.x, v1.y + v2.y, v1.z + v2.z}
}

// Subtract returns the difference between two vectors
func (v1 vec3) Subtract(v2 vec3) vec3 {
	return vec3{v1.x - v2.x, v1.y - v2.y, v1.z - v2.z}
}

// Cross returns the cross product of two vectors
func (v1 vec3) Cross(v2 vec3) vec3 {
	return vec3{v1.y*v2.z - v1.z*v2.y,
		v1.z*v2.x - v1.x*v2.z,
		v1.x*v2.y - v1.y*v2.x}
}

// Dot returns the dot product of two vectors
func (v1 vec3) Dot(v2 vec3) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

// Scale returns a vector multiplied by a scalar
func (v vec3) Scale(s float64) vec3 {
	return vec3{v.x * s, v.y * s, v.z * s}
}

// Normalize returns a normalized vector by scaling by 1 / length of vector
func (v vec3) Normalize() vec3 {
	return v.Scale(1.0 / v.Length())
}

// Length returns the length of a vector
func (v vec3) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

// EqualsZero returns true if all components of a vector are zero
func (v vec3) EqualsZero() bool {
	return v.x == 0 && v.y == 0 && v.z == 0
}

// Ray getter functions

// getOrigin returns the current origin location of a ray
func (r Ray) getOrigin() vec3 { return r.origin }

// getDirection returns the current direction in vector form of a ray
func (r Ray) getDirection() vec3 { return r.direction }

// camera function definitions
// getPosition returns the current x,y,z position of the camera
func (c Camera) getPosition() vec3 { return c.position }

// light function definitions
// getPosition returns the current x,y,z position of a point light
func (l Light) getPosition() vec3 { return l.position }

// collision function definitions
// getPoint returns a point of collision between a ray and a sphere
func (c Collision) getPoint() vec3 { return c.point }

// getNormal returns the normal vector of a sphere at the point of collision
func (c Collision) getNormal() vec3 { return c.normal }

// CenterOfMass calculates the center of mass of a slice of atoms
func CenterOfMass(atoms []*Atom) vec3 {
	sum := vec3{0, 0, 0}
	for i := 0; i < len(atoms); i++ {
		sum = sum.Add(vec3{atoms[i].x, atoms[i].y, atoms[i].z})
	}
	return sum.Scale(1.0 / float64(len(atoms)))
}
