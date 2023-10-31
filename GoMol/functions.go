package main

import (
	"math"
)

// func (r Ray) getOrigin() vec3    { return r.origin }
// func (r Ray) getDirection() vec3 { return r.direction }

func vectorAddition(v1, v2 vec3) vec3 {
	return vec3{v1.x + v2.x, v1.y + v2.y, v1.z + v2.z}
}

func vectorSubtraction(v1, v2 vec3) vec3 {
	return vec3{v1.x - v2.x, v1.y - v2.y, v1.z - v2.z}
}

func normalize(v vec3) vec3 {
	return vec3{v.x / math.Sqrt(v.x*v.x+v.y*v.y+v.z*v.z),
		v.y / math.Sqrt(v.x*v.x+v.y*v.y+v.z*v.z),
		v.z / math.Sqrt(v.x*v.x+v.y*v.y+v.z*v.z)}
}

func crossProduct(v1, v2 vec3) vec3 {
	return vec3{v1.y*v2.z - v1.z*v2.y,
		v1.z*v2.x - v1.x*v2.z,
		v1.x*v2.y - v1.y*v2.x}
}

func dotProduct(v1, v2 vec3) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func multiplyScalar(v vec3, s float64) vec3 {
	return vec3{v.x * s, v.y * s, v.z * s}
}

// camera function definitions
func (c Camera) getPosition() vec3      { return c.position }
func (c Camera) getViewDirection() vec3 { return c.viewDirection }
func (c Camera) getUpVector() vec3      { return c.up }
func (c Camera) getRightVector() vec3   { return c.right }
func (c Camera) getR() float64          { return c.r }
func (c Camera) getL() float64          { return c.l }
func (c Camera) getT() float64          { return c.t }
func (c Camera) getB() float64          { return c.b }
func (c Camera) getD() float64          { return c.d }

func rayGen(cam *Camera, max_height, max_width int, i, j int, view int) Ray {
	d := cam.getD()
	r := cam.getR()
	l := cam.getL()
	t := cam.getT()
	b := cam.getB()
	width := imageWidth
	height := imageHeight
	theta := l + (r-l)*(float64(i)+0.5)/float64(width)
	sigma := b + (t-b)*(float64(j)+0.5)/float64(height)
	w_ray := vec3{0, 0, 1}
	u_ray := vec3{0, 1, 0}
	v_ray := vec3{1, 0, 0}

	var new_ray Ray

	if view == 0 {
		w_ray = multiplyScalar(w_ray, d*-1)
		u_ray = multiplyScalar(u_ray, theta)
		v_ray = multiplyScalar(v_ray, sigma)
		var dir_ray = vec3{w_ray.x + u_ray.x + v_ray.x, w_ray.y + u_ray.y + v_ray.y, w_ray.z + u_ray.z + v_ray.z}
		new_ray = Ray{cam.getPosition(), dir_ray}
	} else {
		var dir_ray = multiplyScalar(u_ray, -1)
		var ray_origin = vec3{u_ray.x + v_ray.x, u_ray.y + v_ray.y, u_ray.z + v_ray.z}
		new_ray = Ray{ray_origin, dir_ray}
	}
	return new_ray
}

// func ray_sphere_collision(ray *Ray, sphere *Sphere) vec3 {

// }
