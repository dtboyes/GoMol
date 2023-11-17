package main

import (
	"fmt"
	"log"
	"math"
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
	camera *Camera
	light  *Light
	atoms1 []*Atom
	atoms2 []*Atom
)

var (
	rotationX, rotationY   float64
	leftMouseButtonPressed bool
	lastX, lastY           float64
)

func init() {
	runtime.LockOSThread()
}

func main() {
	// initialize number of processors
	numProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(numProcs)

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
	atoms1 = ParsePDB("pdbfiles/" + os.Args[1] + ".pdb")
	fmt.Println(atoms1[0])

	atoms2 = ParsePDB("pdbfiles/" + os.Args[2] + ".pdb")

	atoms1_sequence := getQuerySequence(atoms1)
	atoms2_sequence := getQuerySequence(atoms2)

	fmt.Println(NeedlemanWunsch(atoms1_sequence, atoms2_sequence, 2, -1, -2))

	// initialize camera and light
	camera = InitializeCamera(atoms1)
	light = ParseLight("input/light.txt")
	light = InitializeLight(atoms1)

	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetKeyCallback(keyCallback)
	window.SetScrollCallback(scrollCallback)

	// finished := make(chan bool, numProcs)
	// pixels := make([]uint8, 4*imageWidth*imageHeight)
	// RenderScene(camera, light, atoms, 0, imageHeight, pixels, finished)
	// main loop to render scene
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		RotateAtoms(atoms1, rotationX, rotationY)
		pixels := make([]uint8, 4*imageWidth*imageHeight)
		RenderMultiProc(pixels, numProcs)
		gl.DrawPixels(imageWidth, imageHeight, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
		gl.LoadIdentity()
		window.SwapBuffers()
		glfw.PollEvents()

	}
}

func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft {
		if action == glfw.Press {
			leftMouseButtonPressed = true
			lastX, lastY = w.GetCursorPos()
		} else if action == glfw.Release {
			leftMouseButtonPressed = false
		}
	}
}

func cursorPosCallback(window *glfw.Window, xpos, ypos float64) {
	if leftMouseButtonPressed {
		dx := xpos - lastX
		dy := ypos - lastY

		// Adjust rotation angles based on mouse movement
		rotationY += dx * 0.005
		rotationX += dy * 0.005

		// Limit rotation around X-axis to prevent flipping
		rotationX = math.Max(-math.Pi/2, math.Min(math.Pi/2, rotationX))
	} else {
		rotationX = 0.0
		rotationY = 0.0
	}

	lastX, lastY = xpos, ypos
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
		} else if key == glfw.KeyW {
			camera.position.y += 0.1
		} else if key == glfw.KeyS {
			camera.position.y -= 0.1
		} else if key == glfw.KeyA {
			camera.position.x += 0.1
		} else if key == glfw.KeyD {
			camera.position.x -= 0.1
		}
	}
}

func scrollCallback(window *glfw.Window, xoff, yoff float64) {
	camera.position.z += yoff * 0.1
	camera.position.z -= xoff * 0.1
}
