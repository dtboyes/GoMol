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
	atoms  []*Atom
)

var (
	pixels []uint8
)

var (
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
	atoms = ParsePDB("pdbfiles/" + os.Args[1] + ".pdb")
	fmt.Println(atoms[0])

	// initialize camera and light
	camera = InitializeCamera(atoms)
	light = ParseLight("input/light.txt")

	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetKeyCallback(keyCallback)
	window.SetScrollCallback(scrollCallback)
	// RenderScene(camera, light, atoms)
	// pixels = make([]uint8, 4*imageWidth*imageHeight)
	// RenderScene(camera, light, atoms, 0, imageHeight, pixels)
	pixels = make([]uint8, 4*imageWidth*imageHeight)
	RenderMultiProc(pixels, numProcs)
	// main loop to render scene
	for !window.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
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
		camera.yaw += camera.speed * (xpos - lastX)
		camera.pitch += camera.speed * (ypos - lastY)
		if camera.pitch > 89.0 {
			camera.pitch = 89.0
		}
		if camera.pitch < -89.0 {
			camera.pitch = -89.0
		}
		camera.position.x += camera.radius * math.Cos(degToRad(camera.yaw)) * math.Cos(degToRad(camera.pitch))
		camera.position.y = camera.radius * math.Sin(degToRad(camera.pitch))
		camera.position.z += camera.radius * math.Sin(degToRad(camera.yaw)) * math.Cos(degToRad(camera.pitch))
	}
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
		} else if key == glfw.KeyW {
			light.position.y += 0.1
		} else if key == glfw.KeyS {
			light.position.y -= 0.1
		} else if key == glfw.KeyA {
			light.position.x -= 0.1
		} else if key == glfw.KeyD {
			camera.position.x += 0.1
		}
	}
}

func scrollCallback(window *glfw.Window, xoff, yoff float64) {
	camera.position.z += yoff * 0.1
	camera.position.z -= xoff * 0.1
}

func degToRad(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
