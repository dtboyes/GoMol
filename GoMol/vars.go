package main

const (
	imageWidth     = 800
	aspectRatio    = 16.0 / 9.0
	imageHeight    = imageWidth / aspectRatio
	viewportHeight = 2.0
	viewportWidth  = viewportHeight * float64(imageWidth) / float64(imageHeight)
)

var (
	camera            *Camera
	light             *Light
	atoms1            []*Atom
	atoms2            []*Atom
	atoms1_sequence   string
	atoms2_sequence   string
	alignedSeq1       string
	alignedSeq2       string
	matchLine         string
	percentSimilarity float64
)

var (
	rotationX, rotationY   float64
	leftMouseButtonPressed bool
	lastX, lastY           float64
)

var (
	colorByChain            = false
	colorByAtom             = false
	colorByDifferingRegions = true
	onlyChainA              = false
	renderProtein1          = false
	renderProtein2          = false
	renderBoth              = true
)

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
	amino    string
	chain    string
	seqIndex int
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

// AminoPair represents a pair of amino acids
type AminoPair struct {
	First  rune
	Second rune
}

// BLOSUM62 scoring matrix
var BLOSUM62 = make(map[AminoPair]int)
