package main

type vec3 struct {
	x, y, z float64
}

// define ray object, which has an origin and direction
type Ray struct {
	origin    vec3
	direction vec3
	at        vec3
	color     Color
}

type Camera struct {
	position       vec3
	focalLength    float64
	viewportHeight float64
	viewportWidth  float64
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
