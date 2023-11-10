package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

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

var (
	rotationX, rotationY float64
	lastX, lastY         float64
	isDragging           bool
)

func init() {
	runtime.LockOSThread()
}

func main() {

	// Initializing Camera

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

	window, err := glfw.CreateWindow(imageWidth, imageHeight, "GoMol: "+strings.ToUpper(os.Args[1]), nil, nil)
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

	// Main loop
	camera := InitializeCamera(atoms)
	light := ParseLight("input/light.txt")

	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetKeyCallback(keyCallback)
	RenderScene(camera, light, atoms)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.DrawPixels(imageWidth, imageHeight, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
		gl.LoadIdentity()

		gl.Rotatef(float32(rotationX), float32(camera.position.x), float32(camera.position.y), float32(camera.position.z))
		gl.Rotatef(float32(rotationY), float32(camera.position.x), float32(camera.position.y), float32(camera.position.z))
		window.SwapBuffers()
		glfw.PollEvents()

	}
}

func mouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft {
		if action == glfw.Press {
			isDragging = true
			lastX, lastY = window.GetCursorPos()
		} else if action == glfw.Release {
			isDragging = false
		}
	}
}

func cursorPosCallback(window *glfw.Window, xpos, ypos float64) {
	if isDragging {
		dx := xpos - lastX
		dy := ypos - lastY

		rotationX += dy * 0.1
		rotationY += dx * 0.1

		lastX, lastY = xpos, ypos
	}
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press && key == glfw.KeyEscape {
		window.SetShouldClose(true)
	}
}
