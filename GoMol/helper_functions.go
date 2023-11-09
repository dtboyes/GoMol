package main

import "math"

// implementing basic vector operations to use in ray generation
func (v1 vec3) vectorAddition(v2 vec3) vec3 {
	return vec3{v1.x + v2.x, v1.y + v2.y, v1.z + v2.z}
}

func (v1 vec3) vectorSubtraction(v2 vec3) vec3 {
	return vec3{v1.x - v2.x, v1.y - v2.y, v1.z - v2.z}
}

func normalize(v vec3) vec3 {
	return vec3{v.x / v.length(),
		v.y / v.length(),
		v.z / v.length()}
}

func crossProduct(v1, v2 vec3) vec3 {
	return vec3{v1.y*v2.z - v1.z*v2.y,
		v1.z*v2.x - v1.x*v2.z,
		v1.x*v2.y - v1.y*v2.x}
}

func (v1 vec3) dotProduct(v2 vec3) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func (v vec3) multiplyScalar(s float64) vec3 {
	return vec3{v.x * s, v.y * s, v.z * s}
}

func (v vec3) unitVector() vec3 {
	return v.multiplyScalar(1 / v.length())
}

func (v vec3) length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

// Ray getter functions
func (r Ray) getOrigin() vec3    { return r.origin }
func (r Ray) getDirection() vec3 { return r.direction }

// Ray helper functions
func (r Ray) getAt(t float64) vec3 {
	return vec3{r.origin.x + r.direction.x*t,
		r.origin.y + r.direction.y*t,
		r.origin.z + r.direction.z*t}
}

// camera function definitions
func (c Camera) getPosition() vec3 { return c.position }
