package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	imageWidth     = 800
	aspectRatio    = 16.0 / 9.0
	imageHeight    = imageWidth / aspectRatio
	viewportHeight = 2.0
	viewportWidth  = viewportHeight * float64(imageWidth) / float64(imageHeight)
)

var (
	pixels []uint8
)

func init() {
	runtime.LockOSThread()
}

func main() {

	// Initializing Camera
	camera := InitializeCamera()
	viewport_u := vec3{camera.viewportWidth, 0, 0}
	viewport_v := vec3{0, -camera.viewportHeight, 0}

	// Initializing viewport, pixel delta, and top left pixel location

	// pixel_delta_u = viewport_u / imageWidth
	pixel_delta_u := viewport_u.multiplyScalar(1.0 / float64(imageWidth))
	// pixel_delta_u = viewport_v / imageHeight
	pixel_delta_v := viewport_v.multiplyScalar(1.0 / float64(imageHeight))

	// uppper left of viewport is the camera position minus half of the viewport width and height minus the focal length
	viewport_upper_left := camera.position.vectorSubtraction(viewport_u.multiplyScalar(0.5)).vectorSubtraction(viewport_v.multiplyScalar(0.5)).vectorSubtraction(vec3{0, 0, camera.focalLength})

	// top left pixel location is the upper left viewport location plus half of the pixel width and height
	pixel00Location := viewport_upper_left.vectorAddition(pixel_delta_u.multiplyScalar(0.5).vectorAddition(pixel_delta_v.multiplyScalar(0.5)))

	// DOWNLOAD PDB FILE
	pdbURL := "https://files.rcsb.org/download/" + os.Args[1] + ".pdb"
	localPath := "pdbfiles/" + os.Args[1] + ".pdb"
	err := downloadPDB(pdbURL, localPath)
	if err != nil {
		fmt.Println("Error downloading PDB file:", err)
	}

	// Initialize GLFW and create a window
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(imageWidth, imageHeight, "Ray Generation", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	window.MakeContextCurrent()
	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	gl.Viewport(0, 0, imageWidth, imageHeight)

	// parse pdb file to get list of atom objects
	atoms := ParsePDB("pdbfiles/" + os.Args[1] + ".pdb")
	fmt.Println(atoms[0])

	// Generate rays and render them
	RenderScene(camera, atoms, pixel00Location, pixel_delta_u, pixel_delta_v, viewport_upper_left)

	// Main loop
	for !window.ShouldClose() {
		gl.DrawPixels(imageWidth, imageHeight, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
		window.SwapBuffers()
		glfw.PollEvents()

	}
}
