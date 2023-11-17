package main

type vec3 struct {
	x, y, z float64
}

// define ray object, which has an origin and direction
type Ray struct {
	origin    vec3
	direction vec3
	color     Color
}

type Camera struct {
	position       vec3
	radius         float64
	yaw            float64
	pitch          float64
	speed          float64
	focalLength    float64
	viewportHeight float64
	viewportWidth  float64
	pixel00        vec3
	pixelDeltaU    vec3
	pixelDeltaV    vec3
}

type Light struct {
	position vec3
}

type Atom struct {
	number   int
	element  string
	residue  string
	chain    string
	sequence string
	x, y, z  float64
	radius   float64
}

type Color struct {
	r, g, b, a uint8
}

type Collision struct {
	point  vec3
	normal vec3
	color  vec3
}
