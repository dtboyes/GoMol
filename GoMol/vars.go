package main

type vec3 struct {
	x, y, z float64
}

type windowSize struct {
	width, height int
}

// define ray object, which has an origin and direction
type Ray struct {
	origin    vec3
	direction vec3
}

type CameraFunctions interface {
	getPosition()
	getViewDirection()
	getUpVector()
	getRightVector()
	getR()
	getL()
	getT()
	getB()
	getD()
}

type Sphere struct {
	radius float64
	x, y   float64
}

type Camera struct {
	position, viewDirection vec3
	up, right               vec3
	r, l, t, b, d           float64
}
